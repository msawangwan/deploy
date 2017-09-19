package dock

import (
	"testing"

	"github.com/msawangwan/ci.io/lib/dockutil"
)

func TestBuildAPIURLs(t *testing.T) {
	var (
		apiurl APIURL
		res    string
		err    error
	)

	didPass := func() {
		if err != nil {
			t.Errorf("%s", err)
		} else {
			t.Logf("%s", res)
		}
	}

	apiurl = APIURL{
		Command: "containers",
		Option:  "json",
	}

	res, err = dockutil.ResolveAPIEndpoint(&List{apiurl})

	didPass()

	apiurl = APIURL{
		Command: "containers",
		Option:  "create",
		Parameters: map[string]string{
			"name": "container_name",
		},
	}

	res, err = dockutil.ResolveAPIEndpoint(&Create{apiurl})

	didPass()

	apiurl = APIURL{
		Command: "containers",
		Option:  "json",
		ID:      "1234598765abcdefg",
	}

	res, err = dockutil.ResolveAPIEndpoint(&Inspect{apiurl})

	didPass()
}
