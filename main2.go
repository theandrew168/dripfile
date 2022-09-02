package main

import (
	"bytes"
	"fmt"

	"github.com/theandrew168/dripfile/internal/html"
	"github.com/theandrew168/dripfile/internal/html/site"
)

func main() {
	tmpl := html.New()

	var b bytes.Buffer
	err := tmpl.Site.Index(&b, site.IndexParams{})
	err := site.Index(&b, 
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(b.String())
}
