package dockr

// Port is a json field for the containers payload
type Port struct {
	PrivatePort float64 `json:"PrivatePort"`
	PublicPort  float64 `json:"PublicPort"`
	Type        string  `json:"Type"`
}

// Labels is a json field for the containers payload
type Labels struct {
	ComExampleVendor  string `json:"com.example.vendor"`
	ComExampleLicense string `json:"com.example.license"`
	ComExampleVersion string `json:"com.example.version"`
}

// HostConfig is a json field for the containers payload
type HostConfig struct {
	NetworkMode string `json:"NetworkMode"`
}

// Bridge is a json field for the containers payload
type Bridge struct {
	NetworkID           string  `json:"NetworkID"`
	EndpointID          string  `json:"EndpointID"`
	Gateway             string  `json:"Gateway"`
	IPAddress           string  `json:"IPAddress"`
	IPPrefixLen         float64 `json:"IPPrefixLen"`
	IPv6Gateway         string  `json:"IPv6Gateway"`
	GlobalIPv6Address   string  `json:"GlobalIPv6Address"`
	GlobalIPv6PrefixLen float64 `json:"GlobalIPv6PrefixLen"`
	MacAddress          string  `json:"MacAddress"`
}

// Networks is a json field for the containers payload
type Networks struct {
	Bridge Bridge `json:"bridge"`
}

// NetworkSettings is a json field for the containers payload
type NetworkSettings struct {
	Networks Networks `json:"Networks"`
}

// Mounts is a json field for the containers payload
type Mounts struct {
	Name        string `json:"Name"`
	Source      string `json:"Source"`
	Destination string `json:"Destination"`
	Driver      string `json:"Driver"`
	Mode        string `json:"Mode"`
	RW          bool   `json:"RW"`
	Propagation string `json:"Propagation"`
}

type State struct {
	Status     string  `json:"Status"`
	Running    bool    `json:"Running"`
	Paused     bool    `json:"Paused"`
	Restarting bool    `json:"Restarting"`
	OOMKilled  bool    `json:"OOMKilled"`
	Dead       bool    `json:"Dead"`
	PID        float64 `json:"Pid"`
	ExitCode   float64 `json:"ExitCode"`
	Error      string  `json:"Error"`
	StartedAt  string  `json:"StartedAt"`
	FinishedAt string  `json:"FinishedAt"`
}

type Node struct{}
