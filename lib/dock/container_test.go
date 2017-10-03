package dock

import (
	"testing"
)

func TestASimpleCeateContainerPayload(t *testing.T) {
	payload := CreateContainerPayload{Image: "some_img"}

	result, err := payload.Build()
	if err != nil {
		t.Fatalf("%s", err)
	}

	t.Logf("%s", result)
}

func TestContainerWithExposedPorts(t *testing.T) {
	payload := CreateContainerPayload{Image: "some_image", Port: "8080"}

	result, err := payload.Build()
	if err != nil {
		t.Fatalf("%s", err)
	}

	t.Logf("%s", result)
}
