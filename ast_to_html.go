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

	f2 := struct {
		field
		FieldNameWithFirstPoint string
	}{}
	{
		var f field
		err = f.Parse(a, structName)
		if err != nil {
			return
		}
		f2.FieldName = f.FieldName
		f2.Docs = f.Docs
	}

	var buf bytes.Buffer
	Parameter.Source.WriteString("\n")
	Parameter.Source.WriteString(fmt.Sprintf("	/"+"/ %v\n", structName+f2.FieldName)) // comment
	// add docs
	if f2.Docs != "" {
		Parameter.Source.WriteString(fmt.Sprintf(
			"\n\n\tout += fmt.Sprintf(\"\\n<br><strong>%s</strong><br>\\n\")\n", f2.Docs))
	}

	// imports
	AddImport("fmt")

	// convert types
	switch v := a.Type.(type) {
	case *ast.StructType:
		// parse nested struct
		for _, fss := range v.Fields.List {
			err = structToHtml(fss, structName+f2.FieldName+".")
			if err != nil {
				return
			}
		}

	case *ast.Ident:
		index := strings.Index(structName, ".")
		if index < 0 {
			return fmt.Errorf("cannot find point : %s", structName)
		}
		{
			fn := f2.FieldName
			f2.FieldName = structName[index:] + fn
			f2.FieldNameWithFirstPoint = structName[index+1:] + fn
		}

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
	"\n<input type=\"text\" name=\"%s{{ .FieldNameWithFirstPoint }}\" value=\"%s\"><br>\n",
	prefix, fmt.Sprintf("%v", value{{ .FieldName }}))`

			t := template.New("Ident template")
			if t, err = t.Parse(tmpl); err != nil {
				return
			}

			if err = t.Execute(&buf, f2); err != nil {
				return
			}

		default: // user struct
			buf.WriteString(
				"out += value.toHtml(fmt.Sprintf(\"%s" + f2.FieldNameWithFirstPoint + ".\",prefix))")
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
