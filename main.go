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

func ResetParameter() {
	Parameter.InputFilename = []string{}
	Parameter.OutputFilename = ""
	Parameter.Structs = []string{}
	Parameter.PackageName = "main"
}

// pipe for outpur information
var osStdout = os.Stdout

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

	// flags
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

	// parsing file to ast
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

	// parsing to HTML, Go
	et.Name = "Parsing go to html, html to go"
	var h []H2go
	var s []S2html
	for i := range files {
		for k := range files[i].Decls {
			decl, ok := files[i].Decls[k].(*ast.GenDecl)
			if !ok {
				continue
			}
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

	err = ioutil.WriteFile(Parameter.OutputFilename, body(h, s), 0644)
	if err != nil {
		return err
	}

	// TODO check output file by gofmt

	return nil
}

type Parser interface {
	Parse(a *ast.Field, structName string)
}

func parsing(decl *ast.GenDecl, structName string) (h2s H2go, s2h S2html, err error) {
	// check : is this ast have struct name
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
	if tc.Name.Name != structName {
		return
	}

	// is this struct
	fl, ok := tc.Type.(*ast.StructType)
	if !ok {
		err = fmt.Errorf("Not correct type : %T", tc.Type)
		return
	}

	// parsing by parts
	ps := []Parser{&h2s, &s2h}
	for _, fs := range fl.Fields.List {
		for i := range ps {
			ps[i].Parse(fs, structName)
			fmt.Printf("> %#v\n", ps[i])
		}
	}
	return
}

func body(h2s []H2go, s2h []S2html) (b []byte) {
	var buf bytes.Buffer

	// general comment
	buf.WriteString("// Code generated by gensf. DO NOT EDIT.\n\n")

	// header
	buf.WriteString("package main\n\n" +
		"import \"fmt\"\n" +
		"import \"strconv\"\n" +
		"import errors \"github.com/Konstantin8105/errors\"\n" +
		"import \"net/http\"\n\n")

	for i := range h2s {
		buf.WriteString(h2s[i].String())
	}
	for i := range s2h {
		buf.WriteString(s2h[i].String())
	}

	return buf.Bytes()
}

type field struct {
	Name      string
	Docs      string
	ValueName string
}

func (f *field) Parse(a *ast.Field, structName string) {
	if len(a.Names) != 1 {
		// Panic with debug information for understood
		ast.Print(token.NewFileSet(), a)
		fmt.Fprintf(os.Stderr, "Too many names\n")
		return
	}

	f.Name = structName + "." + a.Names[0].Name

	if a.Doc != nil {
		for i := 0; i < len(a.Doc.List); i++ {
			f.Docs += a.Doc.List[i].Text[2:] // [2:] for remove words:"//","/*"
		}
		f.Docs = strings.TrimSpace(f.Docs)
	}
	// not allowable empty documentation
	if len(f.Docs) == 0 {
		// if docs is empty
		fmt.Fprintf(os.Stderr, "Struct `%s` haven`t documentation\n", structName)
	}

	f.Name = strconv.Quote(f.Name)
	f.Docs = strconv.Quote(f.Docs)
	if len(f.Name) > 2 {
		f.Name = f.Name[1 : len(f.Name)-1]
	}
	if len(f.Docs) > 2 {
		f.Docs = f.Docs[1 : len(f.Docs)-1]
	}
}

// func parse(basename string, fs *ast.Field) {
//
// 	// var out, par string
// 	//
// 	// y := struct {
// 	// 	Name      string
// 	// 	Docs      string
// 	// 	ValueName string
// 	// }{}
//
// 	switch v := fs.Type.(type) {
//
// 	case *ast.ArrayType:
// 		// TODO specific
// 		fmt.Println("Type is array of :", v.Elt.(*ast.Ident).Name)
//
// 	case *ast.Ident:
// 		// Example of html form:
// 		//      out += fmt.Printf(
// 		//      "\n%s :<br>\n<input type=\"text\" name=\"%s\" value=\"%s\"><br>\n",
// 		//      "A is some value","P.A", fmt.Sprintf("%v",p.A))
//
// 		index := strings.Index(y.Name, ".")
// 		if index < 0 {
// 			panic(y.Name)
// 		}
// 		y.ValueName = "value" + y.Name[index:]
//
// 		tmpl := `	out += fmt.Sprintf(
// 		"\n{{ .Docs }} :<br>\n<input type=\"text\" name=\"{{ .Name }}\" value=\"%s\"><br>\n",
// 		fmt.Sprintf("%v", {{ .ValueName }}))`
//
// 		t := template.New("Person template")
//
// 		t, err := t.Parse(tmpl)
// 		if err != nil {
// 			panic(err)
// 		}
//
// 		var buf bytes.Buffer
// 		err = t.Execute(&buf, y)
// 		if err != nil {
// 			panic(err)
// 		}
//
// 		out += "\n"
// 		out += fmt.Sprintf("	/"+"/ %v\n", y.Name)
// 		out += buf.String()
// 		out += "\n"
//
// 		// -----------------
// 		buf.Reset()
// 		switch v.Name {
//
// 		// float64
// 		case "float64":
// 			tmpl = `
// 	{
// 		// check if exist
// 		if str, ok := r.Form["{{ .Name }}"]; ok{
// 			if len(str) == 1 {
// 				v, err := strconv.ParseFloat(str[0],64)
// 				if err != nil {
// 					et.Add(err)
// 				} else {
// 					{{ .ValueName }} = v
// 				}
// 			}
// 		}
// 	}
// `
// 			t := template.New("Person template")
//
// 			t, err := t.Parse(tmpl)
// 			if err != nil {
// 				panic(err)
// 			}
//
// 			err = t.Execute(&buf, y)
// 			if err != nil {
// 				panic(err)
// 			}
//
// 		// string
// 		case "string":
// 			tmpl = `
// 	{
// 		if str, ok := r.Form["{{ .Name }}"]; ok{
// 			if len(str) == 1 {
// 				{{ .ValueName }} = str[0]
// 			}
// 		}
// 	}
// `
// 			t := template.New("Person template")
//
// 			t, err := t.Parse(tmpl)
// 			if err != nil {
// 				panic(err)
// 			}
//
// 			err = t.Execute(&buf, y)
// 			if err != nil {
// 				panic(err)
// 			}
//
// 		// int
// 		case "int":
// 			tmpl = `
// 	{
// 		// check if exist
// 		if str, ok := r.Form["{{ .Name }}"]; ok{
// 			if len(str) == 1 {
// 				v, err := strconv.ParseInt(str[0],10,64)
// 				if err != nil {
// 					et.Add(err)
// 				} else {
// 					{{ .ValueName }} = int(v)
// 				}
// 			}
// 		}
// 	}
// `
// 			t := template.New("Person template")
//
// 			t, err := t.Parse(tmpl)
// 			if err != nil {
// 				panic(err)
// 			}
//
// 			err = t.Execute(&buf, y)
// 			if err != nil {
// 				panic(err)
// 			}
//
// 		default:
// 			fmt.Println("> Type >", v.Name)
// 		}
//
// 		par += "\n"
// 		par += fmt.Sprintf("	/"+"/ Name : %v\n", y.Name)
// 		par += fmt.Sprintf("	/"+"/ Type : %v\n", v.Name)
// 		par += buf.String()
// 		par += "\n"
//
// 	default:
// 		// debug
// 		ast.Print(token.NewFileSet(), fs.Type)
// 	}
// }
