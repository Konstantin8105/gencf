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
	Parameter.Source.WriteString("\n")
	Parameter.Source.WriteString(fmt.Sprintf("	/"+"/ %v\n", f.FieldName)) // comment
	// add docs
	if f.Docs != "" {
		Parameter.Source.WriteString(fmt.Sprintf(
			"\n\n\tout += fmt.Sprintf(\"\\n<br><strong>%s</strong><br>\\n\")\n", f.Docs))
	}

	// imports
	AddImport("fmt")

	// convert types
	switch v := a.Type.(type) {
	case *ast.StructType:
		// parse nested struct
		Parameter.Source.WriteString(
			fmt.Sprintf("// found nested struct : %s in %s", f.FieldName, structName))
		for _, fss := range v.Fields.List {
			err = structToHtml(fss, structName+"."+f.FieldName+".")
			if err != nil {
				return
			}
		}

	case *ast.Ident:
		// Example of html form:
		//      out += fmt.Printf(
		//      "\n%s :<br>\n<input type=\"text\" name=\"%s\" value=\"%s\"><br>\n",
		//      "A is some value","P.A", fmt.Sprintf("%v",p.A))

		// imports
		AddImport("fmt")

		index := strings.Index(structName, ".")
		if index < 0 {
			return fmt.Errorf("cannot find point : %s", structName)
		}
		f.FieldName = structName[index+1:] + f.FieldName

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

			// template
			tmpl := `out += fmt.Sprintf(
	"\n<input type=\"text\" name=\"%s{{ .FieldName }}\" value=\"%s\"><br>\n",
	prefix, fmt.Sprintf("%v", value.{{ .FieldName }}))`

			t := template.New("Ident template")
			if t, err = t.Parse(tmpl); err != nil {
				return
			}

			if err = t.Execute(&buf, f); err != nil {
				return
			}

		default: // user struct
			buf.WriteString(
				"out += value.toHtml(fmt.Sprintf(\"%s" + f.FieldName + ".\",prefix))")
		}

	// case *ast.ArrayType:
	// TODO

	default:
		err = fmt.Errorf("Type is not supported: %T", v)
		return
	}

	Parameter.Source.WriteString(buf.String())
	Parameter.Source.WriteString("\n\n\n")

	return
}
