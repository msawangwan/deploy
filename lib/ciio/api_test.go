package ciio

import (
	"testing"

	"github.com/msawangwan/ci.io/lib/jsonutil"
)

const (
	buildfileFilepath = "../../test/mock/Buildfile.mock.json"
)

func TestLoadingBuildFile(t *testing.T) {
	var (
		payload Buildfile
	)

	if e := jsonutil.FromFilepath(buildfileFilepath, &payload); e != nil {
		t.Fatalf("%s", e)
	}

	if e := jsonutil.PrettyPrintStruct(&payload); e != nil {
		t.Fatalf("%s", e)
	}

	t.Logf("test completed with success")
}

func TestExecuteStagesInBuildFile(t *testing.T) {
	var (
		payload Buildfile
	)

	if e := jsonutil.FromFilepath(buildfileFilepath, &payload); e != nil {
		t.Fatalf("%s", e)
	}

	for _, v := range payload.Stages {
		if e := CompleteStage(v); e != nil {
			t.Fatalf("%s", e)
		}
	}

	t.Logf("test completed with success")
}
