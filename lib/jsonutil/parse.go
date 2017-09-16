package jsonutil

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
)

// FromFile reads a file f into a type struct value v
func FromFile(f string, v interface{}) error {
	raw, err := ioutil.ReadFile(f)

	if err != nil {
		return err
	}

	json.Unmarshal(raw, v)

	return nil
}

// BufPretty formats json into a buffer
func BufPretty(r io.Reader, delim, indent string) (out bytes.Buffer, err error) {
	src, err := ioutil.ReadAll(r)

	if err != nil {
		return
	}

	err = json.Indent(&out, []byte(src), delim, indent)

	return
}
