package main

import (
	"fmt"
	"os"

	"github.com/theandrew168/dripfile/internal/html"
	"github.com/theandrew168/dripfile/internal/html/site"
)

func main() {
	tmpl := html.New()

	err := tmpl.Site.Index(os.Stdout, site.IndexParams{})
	if err != nil {
		fmt.Println(err)
		return
	}

	err = tmpl.Site.AuthLogin(os.Stdout, site.AuthLoginParams{})
	if err != nil {
		fmt.Println(err)
		return
	}

	err = tmpl.Site.AuthRegister(os.Stdout, site.AuthRegisterParams{})
	if err != nil {
		fmt.Println(err)
		return
	}
}
