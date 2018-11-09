package main

import "go/ast"

type S2html string

func (s S2html) String() (out string) {

	// 		fmt.Sprintf("func (value %s) ToHtml() (out string) {\n", name) +
	// 		out
	//
	// 	// footer
	// 	out += "	return\n"
	// 	out += "}\n"

	return
}

func (s *S2html) Parse(a *ast.Field, structName string) {
}
