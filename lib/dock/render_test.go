package dock

import (
	"testing"

	"github.com/msawangwan/ci.io/lib/symbol"
)

func TestTemplateBuilders(t *testing.T) {
	payloads := []Templater{
		CreateContainerPayload{Image: "some_image"},
		CreateContainerPayload{Image: "some_image", Port: "8080"},
		StartContainerPayload{ContainerID: "1234"},
		StartContainerPayload{ContainerID: "4312", ContainerPort: "9090"},
		StartContainerPayload{
			ContainerID:   "98765",
			ContainerPort: "8090",
			HostIP:        "10.0.0.1",
		},
		StartContainerPayload{
			ContainerID:   "82384",
			ContainerPort: "9040",
			HostPort:      "80",
		},
		StartContainerPayload{
			ContainerID:   "233456",
			ContainerPort: "9080",
			HostIP:        "192.168.0.1",
			HostPort:      "80",
		},
		CreateContainerAPICall{},
		CreateContainerAPICall{
			Parameters: map[string]string{"some_param": "some_value"},
		},
		CreateContainerAPICall{
			Parameters: map[string]string{"some_param": "some_value", "another_param": "another_value"},
		},
		StartContainerAPICall{
			ContainerID: "1273683",
		},
		BuildImageAPICall{},
		BuildImageAPICall{
			Parameters: map[string]string{"some_key": "some_val"},
		},
		BuildImageAPICall{
			Parameters: map[string]string{"some_key": "some_val", "another_key": "another_val"},
		},
	}

	for _, p := range payloads {
		r, e := p.Build()
		if e != nil {
			t.Fatalf("%s", e)
		}
		t.Logf("%s", r)
	}

	t.Logf("%s", symbol.PassMark)
}
