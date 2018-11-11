package main

func HtmlToStruct(a *ast.Field, structName string) (str string, err error) {
	// 	par = fmt.Sprintf("func (value *%s) FromHtml(r *http.Request) (err error) {\n", name) +
	// 		"	et := errors.New(\"Errors of convert\")\n" +
	// 		par

	// // -----------------
	// buf.Reset()
	// switch v.Name {
	//
	// // float64
	// case "float64":
	// 	tmpl = `
	// 	{
	// 		// check if exist
	// 		if str, ok := r.Form["{{ .Name }}"]; ok{
	// 			if len(str) == 1 {
	// 				v, err := strconv.ParseFloat(str[0],64)
	// 				if err != nil {
	// 					et.Add(err)
	// 				} else {
	// 					{{ .ValueName }} = v
	// 				}
	// 			}
	// 		}
	// 	}
	// `
	// 	t := template.New("Person template")
	//
	// 	t, err := t.Parse(tmpl)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	//
	// 	err = t.Execute(&buf, y)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	//
	// // string
	// case "string":
	// 	tmpl = `
	// 	{
	// 		if str, ok := r.Form["{{ .Name }}"]; ok{
	// 			if len(str) == 1 {
	// 				{{ .ValueName }} = str[0]
	// 			}
	// 		}
	// 	}
	// `
	// 	t := template.New("Person template")
	//
	// 	t, err := t.Parse(tmpl)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	//
	// 	err = t.Execute(&buf, y)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	//
	// // int
	// case "int":
	// 	tmpl = `
	// 	{
	// 		// check if exist
	// 		if str, ok := r.Form["{{ .Name }}"]; ok{
	// 			if len(str) == 1 {
	// 				v, err := strconv.ParseInt(str[0],10,64)
	// 				if err != nil {
	// 					et.Add(err)
	// 				} else {
	// 					{{ .ValueName }} = int(v)
	// 				}
	// 			}
	// 		}
	// 	}
	// `
	// 	t := template.New("Person template")
	//
	// 	t, err := t.Parse(tmpl)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	//
	// 	err = t.Execute(&buf, y)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	//
	// default:
	// 	fmt.Println("> Type >", v.Name)
	// }
	//
	// 		par += "\n"
	// 		par += fmt.Sprintf("	/"+"/ Name : %v\n", y.Name)
	// 		par += fmt.Sprintf("	/"+"/ Type : %v\n", v.Name)
	// 		par += buf.String()
	// 		par += "\n"
	//
	// 	default:
	// 		// debug
	// 		ast.Print(token.NewFileSet(), fs.Type)
	// 	}
	// }

	//
	// 	par += "	if (et.IsError()){\n"
	// 	par += "		return et\n"
	// 	par += "	}\n"
	// 	par += "	return\n"
	// 	par += "}\n"

	return
}
