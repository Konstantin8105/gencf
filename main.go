package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
)

// result of generation:
// P is very big struct
// type P struct {
//		A int // A is some value
// }
//
// func (p P) ToHtml() (out string){
//      out += fmt.Printf("P is very big struct\n")
//      out += fmt.Printf(
//      "\n%s :<br>\n<input type=\"text\" name=\"%s\" value=\"%s\"><br>\n",
//      "A is some value","P.A", fmt.Sprintf("%v",p.A))
//      return
// }

var Parameter = struct {
	InputFilename  []string
	OutputFilename string
	Structs        []string
}{}

type arrayStrings []string

func (a *arrayStrings) String() string {
	return fmt.Sprintf("%v", []string(*a))
}

func (a *arrayStrings) Set(value string) error {
	v := []string(*a)
	v = append(v, value)
	*a = arrayStrings(v)
	return nil
}

func main() {
	// CLI design
	// gensf -struct=foo -struct=buz -o=out_file.go file1.go file2.go
	{
		pif := arrayStrings(Parameter.InputFilename)
		pst := arrayStrings(Parameter.Structs)

		flag.Var(&pif, "i", "input filename for example : 'main.go'")
		flag.Var(&pst, "struct", "name of struct")
		flag.StringVar(&Parameter.OutputFilename, "o", "out_gen.go", "name of output filename")
		flag.Parse()

		Parameter.InputFilename = []string(pif)
		Parameter.Structs = []string(pst)
	}

	// run parsing
	err := run()
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(-1)
	}
}

var out string
var par string

func run() error {

	// print input data
	fmt.Println(Parameter)

	if len(Parameter.InputFilename) == 0 {
		return fmt.Errorf("Please enter input file/files")
	}

	if len(Parameter.Structs) == 0 {
		return fmt.Errorf("Please enter input struct name/names")
	}

	if Parameter.OutputFilename == "" {
		return fmt.Errorf("Please enter not empty output filename")
	}

	// TODO: input filename is exist

	// example of os.Args:
	// [/tmp/go-build649008261/b001/exe/gen -- datasheet]
	if len(os.Args) != 3 || os.Args[1] != "--" {
		panic(fmt.Errorf("Not correct arguments : ", os.Args))
	}

	// start
	name := os.Args[2]
	fmt.Printf("Generate struct `%s`\n", name)

	// get present folder
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	goFiles, err := filepath.Glob("*.go")
	if err != nil {
		panic(err)
	}

	for _, filename := range goFiles {
		out = ""
		par = ""

		// parsing file
		f, err := parser.ParseFile(token.NewFileSet(), pwd+"/"+filename, nil, parser.ParseComments)
		if err != nil {
			return err
		}

		// find information
		for i := range f.Decls {
			if decl, ok := f.Decls[i].(*ast.GenDecl); ok {
				info(decl, name)
			}
		}

		if out == "" || par == "" {
			continue
		}

		// header
		out = "package main\n\n" +
			"import \"fmt\"\n" +
			"import \"strconv\"\n" +
			"import errors \"github.com/Konstantin8105/errors\"\n" +
			"import \"net/http\"\n\n" +
			fmt.Sprintf("func (value %s) ToHtml() (out string) {\n", name) +
			out

		par = fmt.Sprintf("func (value *%s) FromHtml(r *http.Request) (err error) {\n", name) +
			"	et := errors.New(\"Errors of convert\")\n" +
			par

		// footer
		out += "	return\n"
		out += "}\n"

		par += "	if (et.IsError()){\n"
		par += "		return et\n"
		par += "	}\n"
		par += "	return\n"
		par += "}\n"

		// adding
		out += "\n" + par

		filename = pwd + "/" + strings.ToLower(name) + "_gen.go"

		err = os.Remove(filename)
		if err != nil {
			panic(err)
		}

		err = ioutil.WriteFile(filename, []byte(out), 0644)
		if err != nil {
			panic(err)
		}
	}

	return nil
}

func info(decl *ast.GenDecl, name string) {
	if decl.Tok != token.TYPE {
		return
	}
	if len(decl.Specs) != 1 {
		return
	}
	if _, ok := decl.Specs[0].(*ast.TypeSpec); !ok {
		return
	}
	tc := decl.Specs[0].(*ast.TypeSpec)
	if tc.Name.Name != name {
		return
	}

	fl := tc.Type.(*ast.StructType).Fields

	for _, fs := range fl.List {
		parse(name, fs)
	}
}

func parse(basename string, fs *ast.Field) {

	y := struct {
		Name      string
		Docs      string
		ValueName string
	}{}

	if len(fs.Names) != 1 {
		// debug
		ast.Print(token.NewFileSet(), fs)
		panic("Too many names")
	}
	y.Name = basename + "." + fs.Names[0].Name

	if fs.Doc != nil {
		for i := 0; i < len(fs.Doc.List); i++ {
			y.Docs += fs.Doc.List[i].Text[2:]
		}
		y.Docs = strings.TrimSpace(y.Docs)
	}
	if len(y.Docs) == 0 {
		// if docs is empty then error
		fmt.Println("Error empty documentation!!!")
	}

	y.Name = strconv.Quote(y.Name)
	y.Docs = strconv.Quote(y.Docs)
	if len(y.Name) > 2 {
		y.Name = y.Name[1 : len(y.Name)-1]
	}
	if len(y.Docs) > 2 {
		y.Docs = y.Docs[1 : len(y.Docs)-1]
	}

	switch v := fs.Type.(type) {
	case *ast.StructType:
		out += fmt.Sprintf("	out += fmt.Sprintf(\"\\n<br><strong>%s</strong><br>\\n\")\n", y.Docs)
		for _, fss := range v.Fields.List {
			parse(y.Name, fss)
		}

	case *ast.ArrayType:
		// TODO specific
		fmt.Println("Type is array of :", v.Elt.(*ast.Ident).Name)

	case *ast.Ident:
		// Example of html form:
		//      out += fmt.Printf(
		//      "\n%s :<br>\n<input type=\"text\" name=\"%s\" value=\"%s\"><br>\n",
		//      "A is some value","P.A", fmt.Sprintf("%v",p.A))

		index := strings.Index(y.Name, ".")
		if index < 0 {
			panic(y.Name)
		}
		y.ValueName = "value" + y.Name[index:]

		tmpl := `	out += fmt.Sprintf(
		"\n{{ .Docs }} :<br>\n<input type=\"text\" name=\"{{ .Name }}\" value=\"%s\"><br>\n",
		fmt.Sprintf("%v", {{ .ValueName }}))`

		t := template.New("Person template")

		t, err := t.Parse(tmpl)
		if err != nil {
			panic(err)
		}

		var buf bytes.Buffer
		err = t.Execute(&buf, y)
		if err != nil {
			panic(err)
		}

		out += "\n"
		out += fmt.Sprintf("	/"+"/ %v\n", y.Name)
		out += buf.String()
		out += "\n"

		// -----------------
		buf.Reset()
		switch v.Name {

		// float64
		case "float64":
			tmpl = `
	{
		// check if exist
		if str, ok := r.Form["{{ .Name }}"]; ok{
			if len(str) == 1 {
				v, err := strconv.ParseFloat(str[0],64)
				if err != nil {
					et.Add(err)
				} else {
					{{ .ValueName }} = v
				}
			}
		}
	}
`
			t := template.New("Person template")

			t, err := t.Parse(tmpl)
			if err != nil {
				panic(err)
			}

			err = t.Execute(&buf, y)
			if err != nil {
				panic(err)
			}

		// string
		case "string":
			tmpl = `
	{
		if str, ok := r.Form["{{ .Name }}"]; ok{
			if len(str) == 1 {
				{{ .ValueName }} = str[0]
			}
		}
	}
`
			t := template.New("Person template")

			t, err := t.Parse(tmpl)
			if err != nil {
				panic(err)
			}

			err = t.Execute(&buf, y)
			if err != nil {
				panic(err)
			}

		// int
		case "int":
			tmpl = `
	{
		// check if exist
		if str, ok := r.Form["{{ .Name }}"]; ok{
			if len(str) == 1 {
				v, err := strconv.ParseInt(str[0],10,64)
				if err != nil {
					et.Add(err)
				} else {
					{{ .ValueName }} = int(v)
				}
			}
		}
	}
`
			t := template.New("Person template")

			t, err := t.Parse(tmpl)
			if err != nil {
				panic(err)
			}

			err = t.Execute(&buf, y)
			if err != nil {
				panic(err)
			}

		default:
			fmt.Println("> Type >", v.Name)
		}

		par += "\n"
		par += fmt.Sprintf("	/"+"/ Name : %v\n", y.Name)
		par += fmt.Sprintf("	/"+"/ Type : %v\n", v.Name)
		par += buf.String()
		par += "\n"

	default:
		// debug
		ast.Print(token.NewFileSet(), fs.Type)
	}
}
