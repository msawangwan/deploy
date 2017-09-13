package util

import (
	"os"
	"testing"
)

func TestReadJSONFile(t *testing.T) {
	wd, _ := os.Getwd()

	t.Log(wd)
}
