package types

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
)

type jsonFile struct {
	Path string
	Data jsonData
}

type jsonData struct {
	K1 string `json:"k1"`
	K2 string `json:"k2"`
	K3 string `json:"k3"`
	K4 string `json:"k4"`
}

func (j *jsonFile) Read() error {
	raw, err := ioutil.ReadFile(j.Path)

	if err != nil {
		return err
	}

	json.Unmarshal(raw, &j.Data)

	return nil
}

func TestReadJSONFile(t *testing.T) {
	er := os.Chdir("../")

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

	file := &jsonFile{"test/Buildfile.test", data}
	er = file.Read()

	if er != nil {
		t.Error(er)
	}

	t.Logf("%+v", file)
}
