package build

import (
	"bytes"
	"os/exec"
)

// BareRepoer ...
type BareRepoer interface {
	Init(name string, params ...string) (result string, err error)
}

// BareRepo ...
type BareRepo struct {
	WorkspaceDir string
	ScratchDir   string
	Repository   string
}

// Init ...
func (br BareRepo) Init(name string, params ...string) (r string, e error) {
	var stdout, stderr bytes.Buffer

	args := []string{
		br.WorkspaceDir,
		br.ScratchDir,
		br.Repository,
	}

	args = append(args, params...)

	cmd := exec.Command("makebarer", args...)

	cmd.Dir = "__ws"
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if e = cmd.Run(); e != nil {
		return stderr.String(), e
	}

	return stdout.String(), nil
}
