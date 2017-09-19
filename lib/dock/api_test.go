package dock

import (
	"testing"

	"github.com/msawangwan/ci.io/lib/dockutil"
)

func TestBuildAPIURLs(t *testing.T) {
	apiurls := []dockutil.APIEndpointResolver{
		&ListContainers{
			APIURL{
				Command: "containers",
				Option:  "json",
			},
		},
		&CreateContainer{
			APIURL{
				Command: "containers",
				Option:  "create",
				Parameters: map[string]string{
					"name": "container_name",
				},
			},
		},
		&InspectContainer{
			APIURL{
				Command: "containers",
				Option:  "json",
				ID:      "1234598765abcdefg",
			},
		},
	}

	for _, apiurl := range apiurls {
		res, err := dockutil.ResolveAPIEndpoint(apiurl)
		if err != nil {
			t.Errorf("%s", err)
		} else {
			t.Logf("%s", res)
		}
	}
}
