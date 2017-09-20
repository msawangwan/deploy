package jsonutil

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
)

// DecodeFromReader Decodes json in a stream format
func DecodeFromReader(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}

// DecodeFromFilepath is a wrapper for DecodeFromReader
func DecodeFromFilepath(p string, v interface{}) error {
	r, e := os.Open(p)
	if e != nil {
		return e
	}

	return DecodeFromReader(r, v)
}

/* TODO: delete these deprecated functions */

// FromFilepath reads a file f into a type struct value v
func FromFilepath(f string, v interface{}) error {
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
