package main

import (
	"github.com/tomiok/clipper"
	"net/http"
)

func main() {
	clipper.Do(&clipper.Configs{Name: "my_command"}, func() error {
		_, err := http.Get("http://www.google.com/robots.txt")
		return err
	}, nil)

	clipper.Do(&clipper.Configs{Name: "my_command"}, func() error {
		_, err := http.Get("http://www.google.com/robots.txt")
		return err
	}, nil)

	clipper.FillStats("my_command", true)
}
