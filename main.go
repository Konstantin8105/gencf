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

	"github.com/Konstantin8105/errors"
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
	PackageName    string
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
	// gensf -struct=foo -struct=buz -o=out_file.go -i=file1.go -i=file2.go

	pif := arrayStrings(Parameter.InputFilename)
	pst := arrayStrings(Parameter.Structs)

	flag.Var(&pif, "i", "input filename for example : 'main.go'")
	flag.Var(&pst, "struct", "name of struct")
	flag.StringVar(&Parameter.OutputFilename, "o", "out_gen.go", "name of output filename")
	flag.StringVar(&Parameter.PackageName, "p", "main", "package in generate file")
	flag.Parse()

	Parameter.InputFilename = []string(pif)
	Parameter.Structs = []string(pst)

	// run parsing
	err := run()
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(-1)
	}
}

var out string
var par string

var osStdout = os.Stdout

func run() error {
	// check input data
	et := errors.New("Check input data")
	if len(Parameter.InputFilename) == 0 {
		et.Add(fmt.Errorf("input file/files is not added"))
	}
	if len(Parameter.Structs) == 0 {
		et.Add(fmt.Errorf("name of struct is not added"))
	}
	if Parameter.PackageName == "" {
		et.Add(fmt.Errorf("name of package is empty"))
	}
	if Parameter.OutputFilename == "" {
		et.Add(fmt.Errorf("name of output file is empty"))
	}
	for i := range Parameter.InputFilename {
		_, err := os.Stat(Parameter.InputFilename[i])
		if err != nil {
			et.Add(fmt.Errorf("input file `%s` is not exist", Parameter.InputFilename[i]))
		}
	}
	if et.IsError() {
		flag.PrintDefaults()
		return et
	}

	// print input data
	fmt.Fprintf(osStdout, "Generate HTML form from Go struct:\n")
	fmt.Fprintf(osStdout, "Package name: %s\n", Parameter.PackageName)
	fmt.Fprintf(osStdout, "Input go files:\n")
	for i := range Parameter.InputFilename {
		fmt.Fprintf(osStdout, "\t* %s\n", Parameter.InputFilename[i])
	}
	fmt.Fprintf(osStdout, "Parsing next Go structs:\n")
	for i := range Parameter.Structs {
		fmt.Fprintf(osStdout, "\t* %s\n", Parameter.Structs[i])
	}
	fmt.Fprintf(osStdout, "Output go file: %s\n", Parameter.OutputFilename)

	// get present folder
	pwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("Cannot get name of present folder : %v", err)
	}

	// parsing
	var files []*ast.File
	et.Name = "parsing Go files to AST"
	for _, filename := range Parameter.InputFilename {
		f, err := parser.ParseFile(
			token.NewFileSet(),
			filepath.FromSlash(pwd+"/"+filename),
			nil,
			parser.ParseComments)
		if err != nil {
			et.Add(fmt.Errorf("Cannot parse file : %s", filename)).
				Add(err)
		} else {
			files = append(files, f)
		}
	}
	if et.IsError() {
		return et
	}

	// parsing HTML, Go
	et.Name = "Parsing go to html, html to go"
	var h []H2go
	var s []S2html
	for i := range files {
		for k := range files[i].Decls {
			if decl, ok := files[i].Decls[k].(*ast.GenDecl); ok {
				for j := range Parameter.Structs {
					h2s, s2h, err := parsing(decl, Parameter.Structs[j])
					if err != nil {
						et.Add(err)
						continue
					}
					h = append(h, h2s)
					s = append(s, s2h)
				}
			}
		}
	}
	if et.IsError() {
		return et
	}

	// save structs into output file
	if _, err := os.Stat(Parameter.OutputFilename); err == nil {
		err = os.Remove(Parameter.OutputFilename)
		if err != nil {
			return fmt.Errorf("cannot remove file : %v", err)
		}
	}

	return ioutil.WriteFile(Parameter.OutputFilename, body(h, s), 0644)
}

func parsing(a *ast.GenDecl, structName string) (h2s H2go, s2h S2html, err error) {
	return
}

func body(h2s []H2go, s2h []S2html) (b []byte) {
	return
}

// for _, filename := range goFiles {
// 	out = ""
// 	par = ""
//
// 	// parsing file
//
// 	// find information
//
// 	if out == "" || par == "" {
// 		continue
// 	}
//
// 	// header
// 	out = "package main\n\n" +
// 		"import \"fmt\"\n" +
// 		"import \"strconv\"\n" +
// 		"import errors \"github.com/Konstantin8105/errors\"\n" +
// 		"import \"net/http\"\n\n" +
// 		fmt.Sprintf("func (value %s) ToHtml() (out string) {\n", name) +
// 		out
//
// 	par = fmt.Sprintf("func (value *%s) FromHtml(r *http.Request) (err error) {\n", name) +
// 		"	et := errors.New(\"Errors of convert\")\n" +
// 		par
//
// 	// footer
// 	out += "	return\n"
// 	out += "}\n"
//
// 	par += "	if (et.IsError()){\n"
// 	par += "		return et\n"
// 	par += "	}\n"
// 	par += "	return\n"
// 	par += "}\n"
//
// 	// adding
// 	out += "\n" + par
//
// 	filename = pwd + "/" + strings.ToLower(name) + "_gen.go"

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
