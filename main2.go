package main

import (
	"bytes"
	"fmt"

	"github.com/theandrew168/dripfile/internal/html/site"
)

func main() {
	var b bytes.Buffer
	err := site.Index(&b, site.IndexParams{})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(b.String())

	var c bytes.Buffer
	err = site.AuthLogin(&c, site.AuthLoginParams{})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(c.String())
}
