package dock

type Success interface {
	Code() float64
}

type Error interface {
	Message() string
}

// Success200 is a success response
type Success200 struct {
	ID              string        `json:"Id"`
	Created         string        `json:"Created"`
	Path            string        `json:"Path"`
	Args            []string      `json:"Args"`
	State           State         `json:"State"`
	Image           string        `json:"Image"`
	ResolvConfPath  string        `json:"ResolvConfPath"`
	HostnamePath    string        `json:"HostnamePath"`
	HostsPath       string        `json:"HostsPath"`
	LogPath         string        `json:"LogPath"`
	Node            struct{}      `json:"Node"`
	Name            string        `json:"Name"`
	RestartCount    float64       `json:"RestartCount"`
	Driver          string        `json:"Driver"`
	MountLabel      string        `json:"MountLabel"`
	ProcessLabel    string        `json:"ProcessLabel"`
	AppArmorProfile string        `json:"AppArmorProfile"`
	ExecIDs         string        `json:"ExecIDs"`
	HostConfig      HostConfig    `json:"HostConfig"`
	GraphDriver     GraphDriver   `json:"GraphDriver"`
	SizeRW          float64       `json:"SizeRw"`
	SizeRootFS      float64       `json:"SizeRootFs"`
	Mounts          Mounts        `json:"Mounts"`
	Config          Config        `json:"Config"`
	NetworkSettings NetworkConfig `json:"NetworkSettings"`
}

// Success201 is the schema for docker response success
type Success201 struct {
	ID       string   `json:"Id"`
	Warnings []string `json:"Warnings"`
}

// Success204 is a no error success response
type Success204 struct{}

// Error400 is a bad parameter error
type Error400 struct {
	Message string `json:"message"`
}

// Error404 is a no such container error
type Error404 struct {
	Message string `json:"message"`
}

// Error409 is a conflict error
type Error409 struct {
	Message string `json:"message"`
}

// Error500 is an internal server error
type Error500 struct {
	Message string `json:"message"`
}
