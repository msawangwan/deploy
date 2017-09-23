package ciio

import (
	"bytes"
	"io"
	"os"
	"os/exec"
)

// RunStage executes the commands in a stage
func RunStage(s Stage) error {
	for _, c := range s.Commands {
		o, e := ExecuteCommand(c)
		if e != nil {
			return e
		}

		io.Copy(os.Stdout, &o)
	}
	return nil
}

// ExecuteCommand executes a single command
func ExecuteCommand(c Command) (o bytes.Buffer, e error) {
	p := exec.Command(c.Exec, c.Args...)
	p.Stdout = &o
	e = p.Run()
	return
}
