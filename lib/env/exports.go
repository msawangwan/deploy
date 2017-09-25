package env

import "log"

type host string
type port string

// Exports is shared global env state
type Exports struct {
	Host host
	Port port
	*log.Logger
}
