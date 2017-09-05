package util

import (
	"bytes"
	"encoding/json"
)

// PrettyJSON wraps json.Indent, return value can be cast to a string by calling .String()
func PrettyJSON(buf []byte, delimiter, indent string) (bytes.Buffer, error) {
	var (
		out bytes.Buffer
		err error
	)

	err = json.Indent(&out, buf, delimiter, indent)

	return out, err
}
