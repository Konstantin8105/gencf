package main

import "go/ast"

type S2html string

func (s S2html) String() (out string) {
	return string(s)
}

func (s *S2html) Parse(a *ast.Field, structName string) {
	var f field
	f.Parse(a, structName)

	out := string(*s)
	defer func() {
		s = &S2html(out)
	}()

	// header
	out += fmt.Sprintf("func (value %s) ToHtml() (out string) {\n", structName)
	// footer
	defer func() {
		out += "	return\n"
		out += "}\n"
	}()

	// convert types
	switch v := a.Type.(type) {
	case *ast.StructType:
		out += fmt.Sprintf(
			"\tout += fmt.Sprintf(\"\\n<br><strong>%s</strong><br>\\n\")\n", f.Docs)
		for _, fss := range v.Fields.List {
			var temp S2html
			temp.Parse(fss, y.Name)
			out += string(temp)
		}
	}
}
