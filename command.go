package clipper

import (
	"time"
)

type status int

var (
	ok          status = 0
	withErr     status = 1
	withTimeout status = 2
)

type command struct {
	cb               *Clipper
	start            int64
	maxDuration      int
	minDuration      int
	runFunction      func() error
	fallbackFunction func() error
	end              chan bool
	status           chan status
	cmdType          string
}

func Do(name string, fn func() error, fallbackFn func() error) chan status {
	cb := getClipper(name)
	cmd := &command{
		cb:               cb,
		start:            time.Now().Unix(),
		runFunction:      fn,
		fallbackFunction: fallbackFn,
		cmdType:          "async",
		status:           make(chan status, 1),
		end:              make(chan bool, 1),
	}
	return run(cmd)
}

func run(cmd *command) chan status {
	cb := cmd.cb
	cb.mutex.Lock()

	defer cb.mutex.Unlock()
	if cb.isOpen() {
		cmd.status <- 1
		return cmd.status
	}

	go func() {
		defer func() {
			cmd.end <- true
		}()

		err := cmd.runFunction()

		cb.numOfRuns++
		if err != nil {
			cb.update(err)
			if cmd.fallbackFunction != nil {
				err = cmd.fallbackFunction()
				if err != nil {
					cmd.status <- withErr
					return
				} else {
					cmd.status <- ok
					return
				}
			} else {
				cmd.status <- withErr
				return
			}
		}
		cmd.status <- ok
		return
	}()

	go func() {
		timer := time.NewTimer(5 * time.Second)
		defer timer.Stop()
		select {
		case <-cmd.end:
			return
		case <-timer.C:
			cmd.status <- withTimeout
			return
		}

	}()

	return cmd.status
}
