package dock

import (
	"bytes"
	"os/exec"
)

// RepositoryBuilder ...
type RepositoryBuilder interface {
	Build() (result string, err error)
}

// BuildRepository comment
func BuildRepository(bare, tmp, repo, user, oauth string) (result string, er error) {
	var stdout, stderr bytes.Buffer

	args := []string{bare, tmp, repo, user, oauth}

	cmd := exec.Command("echo", args...)

	cmd.Dir = ""
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if er := cmd.Run(); er != nil {
		return stderr.String(), er
	}

	return stdout.String(), nil
}
