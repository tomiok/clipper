package clipper

import (
	"errors"
	"time"
)

type command struct {
	cb               *Clipper
	start            time.Time
	duration         int
	runFunction      func() error
	fallbackFunction func() error
	end              chan bool
	err              chan error
}

func Do(name string, fn func() error, fallbackFn func() error) chan error {
	cb := getClipper(name)
	cmd := &command{
		cb:               cb,
		start:            time.Now(),
		runFunction:      fn,
		fallbackFunction: fallbackFn,
		err:              make(chan error, 1),
	}
	return run(cmd)
}

func run(cmd *command) chan error {
	cb := cmd.cb
	cb.mutex.Lock()

	defer cb.mutex.Unlock()
	if cb.isOpen() {
		// fail fast here
		cmd.err <- errors.New("circuit is open")
		return cmd.err
	}
	var err error
	go func() {
		defer func() {
			cmd.end <- true
		}()

		err = cmd.runFunction()

		cb.numOfRuns++
		if err != nil {
			cb.update(err)
			if cmd.fallbackFunction != nil {
				err = cmd.fallbackFunction()
				if err != nil {
					cmd.err <- err
					return
				}
				return
			}
			cmd.err <- err
			return
		}
		return
	}()

	go func() {
		timer := time.NewTimer(5 * time.Second)
		defer timer.Stop()
		select {
		case <-cmd.end:
			return
		case <-timer.C:
			return
		}
	}()

	return cmd.err
}
