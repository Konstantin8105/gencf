package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"strings"
	"text/template"
)

func structToHtml(a *ast.Field, structName string) (str string, err error) {
	var f field
	f.Parse(a, structName)

	// header
	str += fmt.Sprintf("func (value %s) ToHtml() (out string) {\n", structName)
	// footer
	defer func() {
		str += "\treturn\n"
		str += "}\n\n"
	}()

	// convert types
	switch v := a.Type.(type) {
	case *ast.StructType:
		str += fmt.Sprintf(
			"\tout += fmt.Sprintf(\"\\n<br><strong>%s</strong><br>\\n\")\n", f.Docs)
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

	default:
		err = fmt.Errorf("Type is not supported: %T", v)
		return
	}

	return
}
