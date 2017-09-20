package jsonutil

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
)

// FromFile reads a file f into a type struct value v (TODO: rename to 'FromFilePath(...)')
func FromFile(f string, v interface{}) error {
	b, e := ioutil.ReadFile(f)
	if e != nil {
		return e
	}

	if e = json.Unmarshal(b, v); e != nil {
		return e
	}

	return nil
}

// FromReader wraps ioutil.ReadAll() and json.Unmarshal()
func FromReader(r io.Reader, v interface{}) error {
	b, e := ioutil.ReadAll(r)
	if e != nil {
		return e
	}

	if e = json.Unmarshal(b, v); e != nil {
		return e
	}

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
