package main

import (
	"fmt"
	"go/ast"
)

func structToHtml(a *ast.Field, structName string) (str string, err error) {
	var f field
	f.Parse(a, structName)

	// header
	str += fmt.Sprintf("func (value %s) ToHtml() (out string) {\n", structName)
	// footer
	defer func() {
		str += "	return\n"
		str += "}\n"
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

		// 	case *ast.ArrayType:

	default:
		err = fmt.Errorf("Type is not supported: %T", v)
		return
	}

	return
}
