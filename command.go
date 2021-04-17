package clipper

import (
	"log"
	"time"
)

type command struct {
	cb               *Clipper
	start            time.Time
	duration         int
	runFunction      func() error
	fallbackFunction func() error
}

func Do(name string, fn func() error, fallbackFn func() error) {
	cb := getClipper(name)
	cmd := &command{
		cb:               cb,
		start:            time.Now(),
		runFunction:      fn,
		fallbackFunction: fallbackFn,
	}
	run(cmd)
}

func run(cmd *command) {
	cb := cmd.cb

	if cb.isOpen() {
		// fail fast here
		return
	}

	err := cmd.runFunction()
	cb.numOfRuns++
	if err != nil {
		cb.update(err)
		if cmd.fallbackFunction != nil {
			log.Println(cmd.fallbackFunction())
			return
		}
	}
}
