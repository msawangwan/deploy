package jsonutil

import (
	"os"
	"strings"
	"testing"
)

type jsonData struct {
	K1 string `json:"k1"`
	K2 string `json:"k2"`
	K3 string `json:"k3"`
	K4 string `json:"k4"`
}

func TestFromFile(t *testing.T) {
	er := os.Chdir("../../")
	if er != nil {
		t.Error(er)
	}

	wd, er := os.Getwd()
	if er != nil {
		t.Error(er)
	}

	t.Log(wd)

	var (
		data jsonData
	)

	er = FromFile("test/Buildfile.test", data)
	if er != nil {
		t.Error(er)
	}

	t.Logf("data: %+v", data)
}

func TestBufPretty(t *testing.T) {
	var r = strings.NewReader(`{ "key1": "value1", "key2": "value2", "key3": "value3" }`)

	buf, er := BufPretty(r, "", "  ")
	if er != nil {
		t.Error(er)
	}

	t.Log(buf.String())
}
