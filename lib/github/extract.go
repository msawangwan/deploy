package github

import (
	"encoding/json"
	"io"
)

func ExtractRepositoryName(r io.Reader) (n string, e error) {
	var p *PushEvent

	if e = json.NewDecoder(r).Decode(&p); e != nil {
		return
	}

	n = p.Repository.Name

	return
}
