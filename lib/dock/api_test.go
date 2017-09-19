package dock

import (
	"testing"

	"github.com/msawangwan/ci.io/lib/dockutil"
)

func TestAPIURL(t *testing.T) {
	apiurl := APIURL{
		Command: "containers",
		Option:  "create",
		Parameters: map[string]string{
			"name": "container_name",
		},
	}

	create := &Create{apiurl}

	result, err := dockutil.ResolveAPIEndpoint(create)
	if err != nil {
		t.Errorf("%s", err)
	} else {
		t.Logf("%s", result)
	}
}
