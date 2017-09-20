package jsonutil

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
)

// FromReader Decodes json in a stream format
func FromReader(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}

// FromFilepath is a wrapper for FromReader
func FromFilepath(p string, v interface{}) error {
	r, e := os.Open(p)
	if e != nil {
		return e
	}

	return FromReader(r, v)
}

// BufPretty formats json into a buffer
func BufPretty(r io.Reader, delim, indent string) (out bytes.Buffer, e error) {
	src, e := ioutil.ReadAll(r)
	if e != nil {
		return
	}

	e = json.Indent(&out, []byte(src), delim, indent)

	return
}
