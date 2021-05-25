package clipper

import (
	"sync"
	"time"
)

type status int

var (
	ok          status = 0
	withErr     status = 1
	withTimeout status = 1
)

type command struct {
	cb               *Clipper
	start            int64
	duration         int
	runFunction      func() error
	fallbackFunction func() error
	end              chan bool
	status           chan status
}

func Do(name string, fn func() error, fallbackFn func() error) chan status {
	cb := getClipper(name)
	cmd := &command{
		cb:               cb,
		start:            time.Now().Unix(),
		runFunction:      fn,
		fallbackFunction: fallbackFn,
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

		var err error
		once := &sync.Once{}
		once.Do(func() {
			err = cmd.runFunction()
		})

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
