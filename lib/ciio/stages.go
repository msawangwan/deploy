package ciio

import (
	"bytes"
	"io"
	"os"
	"os/exec"
)

// CompleteStage executes the commands in a stage
func CompleteStage(s Stage) error {
	for _, c := range s.Commands {
		o, e := ExecuteCommand(c)
		if e != nil {
			return e
		}

		io.Copy(os.Stdout, &o)
	}
	return nil
}

// ExecuteCommand executs a single command
func ExecuteCommand(c Command) (o bytes.Buffer, e error) {
	p := exec.Command(c.Exec, c.Args...)
	p.Stdout = &o
	e = p.Run()
	return
}
