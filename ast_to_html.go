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
	err = f.Parse(a, structName)
	if err != nil {
		return
	}

	var buf bytes.Buffer

	// imports
	AddImport("fmt")

	// add docs
	if f.Docs != "" {
		buf.WriteString(fmt.Sprintf(
			"\n\n\tout += fmt.Sprintf(\"\\n<br><strong>%s</strong><br>\\n\")\n", f.Docs))
	}

	// convert types
	switch v := a.Type.(type) {
	case *ast.StructType:
		// parse nested struct
		for _, fss := range v.Fields.List {
			err = structToHtml(fss, f.Name)
			if err != nil {
				return
			}
		}

	case *ast.Ident:
		// Example of html form:
		//      out += fmt.Printf(
		//      "\n%s :<br>\n<input type=\"text\" name=\"%s\" value=\"%s\"><br>\n",
		//      "A is some value","P.A", fmt.Sprintf("%v",p.A))

		// change name
		index := strings.Index(f.Name, ".")
		if index < 0 {
			err = fmt.Errorf("cannot find point of struct : %v", f.Name)
			return
		}
		f.ValueName = "value" + f.Name[index:]

		switch v.Name {
		// Go`s basic types
		case "bool",
			"string",
			"int", "int8", "int16", " int32", "int64",
			"uint", "uint8", "uint16", "uint32", "uint64", "uintptr",
			"byte", // alias for uint8
			"rune", // alias for int32 represents a Unicode code point
			"float32", "float64",
			"complex64", "complex128":

			// imports
			AddImport("fmt")

			// template
			tmpl := `out += fmt.Sprintf(
	"\n<input type=\"text\" name=\"{{ .Name }}\" value=\"%s\"><br>\n",
	fmt.Sprintf("%v", {{ .ValueName }}))`

			t := template.New("Ident template")
			if t, err = t.Parse(tmpl); err != nil {
				return
			}

			if err = t.Execute(&buf, f); err != nil {
				return
			}

		default:

			buf.WriteString(fmt.Sprintf("out += %s.ToHtml()\n", f.ValueName))
			// ast.Print(token.NewFileSet(), a)

		}

	// case *ast.StarExpr:
	// TODO

	// case *ast.ArrayType:
	// TODO

	default:
		err = fmt.Errorf("Type is not supported: %T", v)
		return
	}

	Parameter.Source.WriteString("\n")
	Parameter.Source.WriteString(fmt.Sprintf("	/"+"/ %v\n", f.Name)) // comment
	Parameter.Source.WriteString(buf.String())
	Parameter.Source.WriteString("\n\n\n")

	return
}
