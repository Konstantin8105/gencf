package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/token"
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
		return fmt.Errorf("Field: %v\n%v", f, err)
	}

	var buf bytes.Buffer
	Parameter.Source.WriteString("\n")
	Parameter.Source.WriteString(fmt.Sprintf("	/"+"/ Field : %v\n", f.FieldName)) // comment
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
		for _, fss := range v.Fields.List {
			err = structToHtml(fss, f.StructName+f.FieldName+".")
			if err != nil {
				return
			}
		}

	case *ast.Ident:
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
	"\n<input type=\"text\" name=\"%s{{ .FieldName }}\" value=\"%s\"><br>\n",
	prefix, fmt.Sprintf("%v", value{{ .FieldNameWithFirstPoint }}))`

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

	case *ast.ArrayType:
		// Example
		//
		// *ast.ArrayType {
		// .  Lbrack: -
		// .  Elt: *ast.Ident {
		// .  .  NamePos: -
		// .  .  Name: "string"
		// .  }
		// }
		switch v.Elt.(*ast.Ident).Name {
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
			tmpl := `
	// 
	// Exist elements of field: {{ .FieldName }}
	//
	for i := range value{{ .FieldNameWithFirstPoint }}{
		out += fmt.Sprintf("Data %d<br>\n",i)
		out += fmt.Sprintf(
			"\n<input type=\"text\" name=\"%s{{ .FieldName }}[%d]\" value=\"%s\"><br>\n",
			prefix,i, fmt.Sprintf("%v", value{{ .FieldNameWithFirstPoint }}))
	}

	//
	// Array script of : {{ .FieldName }} 
	//
	out += "<script>\n"
	out += fmt.Sprintf("var initVal{{ .FieldName }} = %d;\n",len(value{{ .FieldNameWithFirstPoint}}))
	out += "\n"
	out += "function insertAfter{{ .FieldName }}(elem, refElem) { \n"
	out += "  var parent = refElem.parentNode; \n"
	out += "  var next = refElem.nextSibling; \n"
	out += "  if (next) { \n"
	out += "    return parent.insertBefore(elem, next); \n"
	out += "  } else { \n"
	out += "    return parent.appendChild(elem); \n"
	out += "  } \n"
	out += "} \n"

	out += "function createEl{{ .FieldName }}(context) { \n"
	out += "	// create input\n"
	out += "	var el = document.createElement(\"input\"); \n"
	out += "	el.type = \"text\"; \n"
	out += "	el.name = \"{{ .StructName }}{{ .FieldName }}[\"+initVal{{ .FieldName }}+\"]\"; \n"
	out += "	initVal{{ .FieldName }}++; \n"
	out += "	insertAfter{{ .FieldName }}(el,context); \n"
	out += "	// create label\n"
	out += "	var label = document.createElement(\"br\");\n"
	out += "	label.id = \"breakLine\" + initVal{{ .FieldName }} + \"{{ .FieldName }}\";\n";
	out += "	insertAfter{{ .FieldName }} (label, context);\n"
	out += " } \n"

	out += "function add{{ .FieldName }}() { \n"
	out += "\t var name = 'breakLine' + initVal{{ .FieldName }} + '{{ .FieldName }}'; \n"
	out += "\t createEl{{ .FieldName }}(document.getElementById(name)); \n"
	out += "\t console.log(name);\n "
	out += "} \n"
	out += "</script>\n"

	out += "<button type=\"button\" OnClick=\"add{{ .FieldName }}()\">+</button>\n"
	out += fmt.Sprintf("<br id=\"breakLine%d{{ .FieldName }}\">\n",len(value{{ .FieldNameWithFirstPoint }}))

	`

			t := template.New("Ident template")
			if t, err = t.Parse(tmpl); err != nil {
				return
			}

			if err = t.Execute(&buf, f); err != nil {
				return
			}

		default:
			// ast.Print(token.NewFileSet(), v)
			// TODO : Uncomment : err = fmt.Errorf("Type is not supported of array: %T. %#v", v, v.Elt)
			Parameter.Source.WriteString(fmt.Sprintf("\n\n// Type is not supported of array: %T. %#v", v, v.Elt.(*ast.Ident).Name))
			return
		}

	default:
		// TODO : Uncomment : err = fmt.Errorf("Type is not supported: %T", v)
		Parameter.Source.WriteString(fmt.Sprintf("\n\n// Type is not supported: %T\n\n", v))
		return
	}

	Parameter.Source.WriteString(buf.String())
	Parameter.Source.WriteString("\n\n\n")

	return
}
