package clipper

import (
	"time"
)

type status int

var (
	finishOk          status = 0
	finishWithErr     status = 1
	finishWithTimeout status = 2
)

type command struct {
	cb               *Clipper
	start            int64
	runFunction      func() error
	fallbackFunction func() error
	end              chan bool
	status           chan status
	cmdType          string
}

// Do will perform the command operation, return a status with
// 0 == no error
// 1 == error (app)
// 2 == timeout
// everything != 0 is counted as error
// you can provide the config
func Do(name string, fn func() error, fallbackFn func() error, cfg *Configs) chan status {
	if cfg != nil {

	}

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

		cb.statistics.numOfRuns++
		if err != nil {
			cb.update(err)
			if cmd.fallbackFunction != nil {
				err = cmd.fallbackFunction()
				if err != nil {
					cmd.status <- finishWithErr
					return
				} else {
					cmd.status <- finishOk
					return
				}
			} else {
				cmd.status <- finishWithErr
				return
			}
		}
		cmd.status <- finishOk
		return
	}()

	go func() {
		timer := time.NewTimer(5 * time.Second)
		defer timer.Stop()
		select {
		case <-cmd.end:
			return
		case <-timer.C:
			cmd.status <- finishWithTimeout
			return
		}
	}()

	return cmd.status
}
