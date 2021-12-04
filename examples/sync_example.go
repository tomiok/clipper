package main

import (
	"fmt"
	"github.com/tomiok/clipper"
	"net/http"
)

func main() {
	err := clipper.DoInSync("my_command", func() error {
		_, err := http.Get("bad-protocol://www.badgoogle.com/robots.txt")
		return err
	}, nil)

	fmt.Println(fmt.Sprintf("%v", err))
}
