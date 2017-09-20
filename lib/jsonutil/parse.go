package jsonutil

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
)

// FromReader is a helper function that decodes a reader (assumed to be valid JSON) into a struct of the same object type
func FromReader(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}

// FromFilepath wraps FromReader, making it easy to read JSON files given a filepath string
func FromFilepath(p string, v interface{}) error {
	r, e := os.Open(p)
	if e != nil {
		return e
	}

	return FromReader(r, v)
}

// BufPretty pretty-print formats a reader into a bytes.Buffer, assumes the reader is JSON
func BufPretty(r io.Reader, delim, indent string) (out bytes.Buffer, e error) {
	src, e := ioutil.ReadAll(r)
	if e != nil {
		return
	}

	e = json.Indent(&out, []byte(src), delim, indent)

	return
}
