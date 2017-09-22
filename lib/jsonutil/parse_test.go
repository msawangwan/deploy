package jsonutil

import (
	"io"
	"strings"
	"testing"
)

const (
	jsonFile = "../../test/mock/Object.mock.json"
)

var jsonReader io.Reader

type jsonData struct {
	K1 string `json:"k1"`
	K2 string `json:"k2"`
	K3 string `json:"k3"`
	K4 string `json:"k4"`
}

func TestDecodeReader(t *testing.T) {
	var data jsonData

	jsonReader = strings.NewReader(
		`{ 
			"k1": "v1",
			"k2": "v2",
			"k3": "v3",
			"k4": "v4"
		}`,
	)

	if e := FromReader(jsonReader, &data); e != nil {
		t.Errorf("%s", e)
	}

	t.Logf("%+v", data)
}

func TestFromFilepath(t *testing.T) {
	var data jsonData

	if e := FromFilepath(jsonFile, &data); e != nil {
		t.Errorf("%s", e)
	}

	t.Logf("%+v", data)
}

func TestExtractBufferFormatted(t *testing.T) {
	jsonReader = strings.NewReader(
		`{ 
			"k1": "v1",
			"k2": "v2",
			"k3": "v3",
			"k4": "v4"
		}`,
	)

	buf, er := ExtractBufferFormatted(jsonReader, "", "  ")
	if er != nil {
		t.Error(er)
	}

	t.Log(buf.String())
}
