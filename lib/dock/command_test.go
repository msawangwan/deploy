package dock

import (
	"testing"
)

func TestBuildAPIURLs(t *testing.T) {
	apiurls := []APIStringBuilder{
		&ContainerCommand{
			URLComponents{
				Command: "containers",
				Option:  "json",
			},
		},
		&ContainerCommand{
			URLComponents{
				Command: "containers",
				Option:  "create",
				Parameters: map[string]string{
					"name": "container_name",
				},
			},
		},
		&ContainerCommandByID{
			URLComponents{
				Command: "containers",
				Option:  "json",
			},
			"1234598765abcdefg",
		},
		&ContainerCommandByID{
			URLComponents{
				Command: "containers",
				Option:  "start",
			},
			"1234598765abcdefg",
		},
		&ContainerCommandByID{
			URLComponents{
				Command: "containers",
				Option:  "stop",
				Parameters: map[string]string{
					"some_param":    "some_value",
					"another_param": "another_value",
				},
			},
			"1234598765abcdefg",
		},
	}

	for _, apiurl := range apiurls {
		// res, err := dockutil.ResolveAPIEndpoint(apiurl)
		res, err := BuildAPIURLString(apiurl)
		if err != nil {
			t.Errorf("%s", err)
		} else {
			t.Logf("%s", res)
		}
	}
}
