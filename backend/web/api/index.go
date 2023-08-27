package api

import (
	"net/http"
)

const html = `
<!DOCTYPE html>
<html lang="en">

<head>
	<title>Dripfile</title>

	<meta charset="utf-8" />
	<meta name="description" content="Dripfile - File transfers made easy" />
	<meta name="viewport" content="initial-scale=1, width=device-width" />

	<script type="module" src="https://unpkg.com/rapidoc/dist/rapidoc-min.js"></script>
</head>

<body>
	<rapi-doc
		spec-url="/static/etc/openapi.yaml"
		render-style="read"
	></rapi-doc>
</body>
</html>
`

func (app *Application) handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(html))
}
