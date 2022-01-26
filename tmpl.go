package main

import (
	"html/template"
	"net/http"
)

// base template that includes CSS / JS, etc
var base string = `
{{define "base"}}
<!DOCTYPE html>
<html lang="en">

<head>
	<title>Nested Template Example</title>

	<meta charset="utf-8" />
	<meta name="viewport" content="width=device-width, initial-scale=1.0" />

	<!-- CSS / JS links would go here -->
</head>

<body>
	{{template "body" .}}
</body>

</html>
{{end}}`

// shared nav content, extends "base"
var app string = `
{{template "base" .}}

{{define "body"}}
<nav>
	<a href="#">Home</a>
	<a href="#">About</a>
	<a href="#">Contact</a>
</nav>
<main>
	{{template "main" .}}
</main>
{{end}}
`

// actual page to render
var page string = `
{{template "body" .}}

{{define "main"}}
<h1>Hello main!</h1>
{{end}}`

func main() {
	http.HandleFunc("/", handleIndex)
	http.ListenAndServe("127.0.0.1:5000", nil)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	t := template.New("example")

	// parse page template
	t, err := t.Parse(page)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// parse app template
	t, err = t.Parse(app)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// parse base template
	t, err = t.Parse(base)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = t.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}
