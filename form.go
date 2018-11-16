package main

func newHandler(w http.ResponseWriter, r *http.Request) {
	d := Nu18007_H101()
	tmpl := `
<!DOCTYPE html>
<html>
<body>

<h1>New projects:</h1>

<form action="/result" target="_blank" method="GET">

` + d.ToHtml() + `

  <br>
  <input type="submit" value="Submit">
</form>

</body>
</html>
`
	w.Write([]byte(tmpl))
}
