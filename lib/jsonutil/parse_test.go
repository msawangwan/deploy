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

var cd = func() error { return os.Chdir("../../") }
var pwd = func() string { wd, _ := os.Getwd(); return wd }

func TestDecoderFromFilepath(t *testing.T) {
	e := cd()
	if e != nil {
		t.Errorf("%s", e)
	}

	wd := pwd()
	t.Logf("%s", wd)

	var data jsonData

	if e = DecodeFromFilepath("test/Testfile.json", &data); e != nil {
		t.Errorf("%s", e)
	}

	t.Logf("%+v", data)
}

func TestFromFilepath(t *testing.T) {
	wd := pwd()
	t.Logf("%s", wd)

	var (
		data jsonData
	)

	e := FromFilepath("test/Testfile.json", &data)
	if e != nil {
		t.Error(e)
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
