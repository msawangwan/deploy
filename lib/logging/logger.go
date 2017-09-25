package logging

import "log"

type Printer interface {
	Log(s ...string)
}

type Empty struct{}

func (e Empty) Log(s ...string) {}

type Debug struct {
	*log.Logger
}

func (d Debug) Log(s ...string) {
	for _, m := range s {
		d.Printf("%s", m)
	}
}
