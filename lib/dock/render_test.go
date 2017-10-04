package dock

import (
	"testing"

	"github.com/msawangwan/ci.io/lib/symbol"
)

func TestTemplateBuilders(t *testing.T) {
	payloads := []Templater{
		CreateContainerPayload{Image: "some_image"},
		NewCreateContainerPayload("some_image", ""),
		NewCreateContainerPayload("some_image", "80/tcp"),
		NewStartContainerPayload("123456", "", "", ""),
		NewStartContainerPayload("123456", "80/tcp", "", ""),
		NewStartContainerPayload("123456", "80/tcp", "", "9091"),
		NewStartContainerPayload("123456", "80/tcp", "192.168.0.1", ""),
		NewStartContainerPayload("123456", "80/tcp", "192.168.0.1", "9091"),
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
