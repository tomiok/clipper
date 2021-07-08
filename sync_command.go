package clipper

import (
	"errors"
	"time"
)

func DoAsync(name string, fn func() error, callbackFn func() error) error {
	cb := getClipper(name)

	cmd := &command{
		cb:               cb,
		start:            time.Now().Unix(),
		runFunction:      fn,
		fallbackFunction: callbackFn,
		end:              make(chan bool, 1),
		status:           make(chan status, 1),
		cmdType:          "sync",
	}

	ch := run(cmd)
	s := <-ch

	if s == 1 {
		return errors.New("error")
	}

	return nil
}
