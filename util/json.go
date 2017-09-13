package util

import (
	"io/ioutil"
	"encoding/json"
)

// ReadFileJSON reads a file f into a type struct value v
func ReadFileJSON(f string, v interface{}) error {
	raw, err := ioutil.ReadFile(f)

	if err != nil {
		return err
	}

	json.Unmarshal(raw, &v)

	return nil
}
