package main

import "fmt"

func createForm(structName string) (err error) {
	AddImport("fmt")
	Parameter.Source.WriteString(fmt.Sprintf(
		`
func (value %s) FormDefault(handlerName string) (out string){
	out += "<!DOCTYPE html>\n"
	out += "<html>\n"
	out += "<body>\n"`, structName))
	Parameter.Source.WriteString(`
	out += fmt.Sprintf("<form action=\"%s\" target=\"_blank\" method=\"GET\">\n", handlerName)
	out += value.ToHtml()
	out += "<input type=\"submit\" value=\"Submit\">"
	out += "</form>"
	out += "<br>\n"
	out += "</body>\n"
	out += "</html>\n"

	return
}
`)
	return nil
}
