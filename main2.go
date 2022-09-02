package main

import (
	"fmt"
	"os"

	"github.com/theandrew168/dripfile/internal/html"
//	"github.com/theandrew168/dripfile/internal/html/api"
//	"github.com/theandrew168/dripfile/internal/html/app"
	"github.com/theandrew168/dripfile/internal/html/site"
)

func main() {
	reload := false
	if os.Getenv("DEBUG") != "" {
		reload = true
	}

	tmpl := html.New(reload)

	err := tmpl.Site.Index(os.Stdout, site.IndexParams{})
	if err != nil {
		fmt.Println(err)
		return
	}

//	err = tmpl.Site.AuthLogin(os.Stdout, site.AuthLoginParams{})
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//
//	err = tmpl.Site.AuthRegister(os.Stdout, site.AuthRegisterParams{})
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//
//	err = tmpl.API.Index(os.Stdout, api.IndexParams{})
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//
//	err = tmpl.App.Dashboard(os.Stdout, app.DashboardParams{})
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
}
