package ciio

// import (
// 	"testing"

// 	"github.com/msawangwan/ci.io/lib/jsonutil"
// )

// const (
// 	buildfileFilepath = "../../test/mock/Buildfile.mock.json"
// )

// var tests = map[string]string{
// 	"loadBuildFile": "load_build_file",
// 	"runStages":     "run_stages",
// }

// func TestLoadingBuildFile(t *testing.T) {
// 	t.Logf("starting test: %s ..", tests["loadBuildFile"])

// 	var (
// 		payload Buildfile
// 	)

// 	if e := jsonutil.FromFilepath(buildfileFilepath, &payload); e != nil {
// 		t.Fatalf("%s", e)
// 	}

// 	if e := jsonutil.PrettyPrintStruct(&payload); e != nil {
// 		t.Fatalf("%s", e)
// 	}

// 	t.Logf("completed test: %s ..", tests["loadBuildFile"])
// }

// func TestRunStagesInBuildFile(t *testing.T) {
// 	t.Logf("starting test: %s ..", tests["runStages"])

// 	var (
// 		payload Buildfile
// 	)

// 	if e := jsonutil.FromFilepath(buildfileFilepath, &payload); e != nil {
// 		t.Fatalf("%s", e)
// 	}

// 	for _, v := range payload.Stages {
// 		if e := RunStage(v); e != nil {
// 			t.Fatalf("%s", e)
// 		}
// 	}

// 	t.Logf("completed test: %s ..", tests["runStages"])
// }
