package main

import (
	"github.com/tomiok/clipper"
	"net/http"
)

func main() {
	var res *http.Response
	clipper.Do(&clipper.Configs{Name: "my_command"}, func() error {
		r, err := http.Get("http://www.google.com/robots.txt")
		res = r
		return err
	}, nil)

	clipper.FillStats("my_command", true)
}
