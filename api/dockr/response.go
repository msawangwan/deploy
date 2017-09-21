package dockr

// Success200 is a success response
type Success200 struct {
	ID              string
	Created         string
	Path            string
	Args            []string
	State           State
	Image           string
	ResolvConfPath  string
	HostnamePath    string
	HostsPath       string
	LogPath         string
	Node            Node
	Name            string
	RestartCount    float64
	Driver          string
	MountLabel      string
	ProcessLabel    string
	AppArmorProfile string
	ExecIDs         string
	HostConfig      HostConfig
	// GraphDriver
	SizeRW     float64
	SizeRootFS float64
	Mounts     Mounts
	// Config
	NetworkSettings NetworkSettings
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
