package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/token"
	"strings"
	"text/template"
)

func structToHtml(a *ast.Field, structName string) (err error) {
	defer func() {
		if err != nil {
			ast.Print(token.NewFileSet(), a)
		}
	}()

	var f field
	f.Parse(a, structName)

	// header
	str += fmt.Sprintf("\nfunc (value %s) ToHtml() (out string) {\n", structName)
	// footer
	defer func() {
		str += "\treturn\n"
		str += "}\n\n"
	}()

	// convert types
	switch v := a.Type.(type) {
	case *ast.StructType:
		// imports
		AddImport("fmt")

		str += fmt.Sprintf(
			"\tout += fmt.Sprintf(\"\\n<br><strong>%s</strong><br>\\n\")\n", f.Docs)

		// parse nested struct
		for _, fss := range v.Fields.List {
			var s string
			s, err = structToHtml(fss, f.Name)
			if err != nil {
				return
			}
			str += s
		}

	case *ast.Ident:
		// Example of html form:
		//      out += fmt.Printf(
		//      "\n%s :<br>\n<input type=\"text\" name=\"%s\" value=\"%s\"><br>\n",
		//      "A is some value","P.A", fmt.Sprintf("%v",p.A))

		index := strings.Index(f.Name, ".")
		if index < 0 {
			err = fmt.Errorf("cannot find point of struct : %v", f.Name)
			return
		}
		f.ValueName = "value" + f.Name[index:]

		// imports
		AddImport("fmt")

		// template
		tmpl := `out += fmt.Sprintf(
	"\n{{ .Docs }} :<br>\n<input type=\"text\" name=\"{{ .Name }}\" value=\"%s\"><br>\n",
	fmt.Sprintf("%v", {{ .ValueName }}))
	`

		t := template.New("Ident template")
		t, err = t.Parse(tmpl)
		if err != nil {
			return
		}

		var buf bytes.Buffer
		err = t.Execute(&buf, f)
		if err != nil {
			return
		}

		str += "\n"
		str += fmt.Sprintf("	/"+"/ %v\n", f.Name)
		str += buf.String()
		str += "\n"

	case *ast.StarExpr:
		// TODO

	case *ast.ArrayType:
		// TODO

	default:
		err = fmt.Errorf("Type is not supported: %T", v)
		return
	}

	return
}
