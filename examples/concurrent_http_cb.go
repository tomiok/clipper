package main

import (
	"fmt"
	"github.com/tomiok/clipper"
	"net/http"
)

func main() {
	//out := make(chan bool, 1)
	var res *http.Response
	_ = clipper.Do("my_command", func() error {
		r, err := http.Get("http://www.google.com/robots.txt")
		res = r
		//out <- true
		return err
	}, nil)

	/*	select {
		case <-out:
			value := <- valChan
			if value == 0{
				fmt.Println("no errors")
			} else {
				fmt.Println("some errors here")
			}
		}*/
	fmt.Println(res)
	clipper.FillStats("my_command", true)
}
