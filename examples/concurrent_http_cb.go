package main

import (
	"fmt"
	"github.com/tomiok/clipper"
	"net/http"
)

func main() {
	out := make(chan bool, 1)
	cerr := clipper.Do("my_command", func() error {
		_, err := http.Get("httpdasdas@adasd")
		out <- true
		return err
	}, nil)

	select {
	case <-out:
		if e, ok := <-cerr; ok {
			fmt.Println(e)
		}
	}

	clipper.FillStats("my_command", true)
}
