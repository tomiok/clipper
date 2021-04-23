package main

import (
	"fmt"
	"github.com/tomiok/clipper"
	"net/http"
)

func main() {
	out := make(chan bool, 1)
	var res *http.Response
	 cerr := clipper.Do("my_command", func() error {
		r, err := http.Get("hsdadasdasdasdttp://www.google.com/robots.txt")
		res = r
		out <- true
		return err
	}, nil)

	select {
	case v:= <-out:
		fmt.Println(v)
	case e := <- cerr:
		fmt.Println(e)

	}
	clipper.FillStats("my_command", true)
}
