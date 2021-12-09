package clipper

import (
	"errors"
	"time"
)

// DoInSync is the same as do, but you should wait until the function finishes.
func DoInSync(name string, fn func() error, callbackFn func() error) error {
	cb := getClipper(&Configs{Name: name})

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
	statusResult := <-ch

	if statusResult != finishOk {
		return errors.New("finish with errors, please checkout previous logs")
	}

	return nil
}
