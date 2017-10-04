package dock

import (
	"testing"

	"github.com/msawangwan/ci.io/lib/symbol"
)

func TestTemplateBuilders(t *testing.T) {
	var (
		img           = "some_img"
		id            = "1234567890"
		containerport = "80/tcp"
		hostport      = "8080"
		hostip        = "127.0.0.1"
	)
	payloads := []Templater{
		NewCreateContainerPayload(img, containerport, hostip, hostport),
		NewCreateContainerPayload(img, containerport, hostip, ""),
		NewCreateContainerPayload(img, containerport, "", hostport),
		NewCreateContainerPayload(img, containerport, "", ""),
		NewCreateContainerPayload(img, "", "", ""),
		CreateContainerAPICall{},
		CreateContainerAPICall{
			Parameters: map[string]string{"some_param": "some_value"},
		},
		CreateContainerAPICall{
			Parameters: map[string]string{"some_param": "some_value", "another_param": "another_value"},
		},
		StartContainerAPICall{
			ContainerID: id,
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
