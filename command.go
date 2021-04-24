package clipper

import (
	"sync"
	"time"
)

type command struct {
	cb               *Clipper
	start            time.Time
	duration         int
	runFunction      func() error
	fallbackFunction func() error
	end              chan bool
	status           chan int
}

func Do(name string, fn func() error, fallbackFn func() error) chan int {
	cb := getClipper(name)
	cmd := &command{
		cb:               cb,
		start:            time.Now(),
		runFunction:      fn,
		fallbackFunction: fallbackFn,
		status:           make(chan int, 1),
		end:              make(chan bool, 1),
	}
	return run(cmd)
}

func run(cmd *command) chan int {
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
					cmd.status <- 1
					return
				} else {
					cmd.status <- 0
					return
				}
			} else {
				cmd.status <- 1
				return
			}
		}
		cmd.status <- 0
		return
	}()

	go func() {
		timer := time.NewTimer(5 * time.Second)
		defer timer.Stop()
		select {
		case <-cmd.end:
			return
		case <-timer.C:
			cmd.status <- 1
			return
		}

	}()

	return cmd.status
}
