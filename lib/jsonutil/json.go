package jsonutil

import (
	"encoding/json"
	"io/ioutil"
)

// FromFile reads a file f into a type struct value v
func FromFile(f string, v interface{}) error {
	raw, err := ioutil.ReadFile(f)

	if err != nil {
		return err
	}

	json.Unmarshal(raw, &v)

	return nil
}
