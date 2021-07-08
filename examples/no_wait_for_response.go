package main

import (
	"github.com/tomiok/clipper"
	"net/http"
)

func main3() {
	var res *http.Response
	clipper.Do("my_command", func() error {
		r, err := http.Get("http://www.google.com/robots.txt")
		res = r
		return err
	}, nil)

	clipper.FillStats("my_command", true)
}
