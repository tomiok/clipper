package clipper

import "time"

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

	if cmd.cb.open {

		return
	}
	run(cmd)
}

func run(cmd *command) {
	cb := cmd.cb

	if !cb.isOpen() {
		// cannot run

		return
	}

	cb.update(cmd.runFunction())
}
