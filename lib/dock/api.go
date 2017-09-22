package dock

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

// Port is a json field for the containers payload
type Port struct {
	IP          string  `json:"IP"`
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

// State is a json field
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

// GraphDriver is a json field
type GraphDriver struct {
	Name string   `json:"Name"`
	Data struct{} `json:"Data"`
}

// HealthConfig is a json field
type HealthConfig struct {
	Test        []string `json:"Test"`
	Interval    float64  `json:"Interval"`
	Timeout     float64  `json:"Timeout"`
	Retries     float64  `json:"Retries"`
	StartPeriod float64  `json:"StartPeriod"`
}

// Config is a json field
type Config struct {
	HostName        string       `json:"Hostname"`
	DomainName      string       `json:"Domainname"`
	User            string       `json:"User"`
	AttachStdin     bool         `json:"AttachStdin"`
	AttachStdout    bool         `json:"AttachStdout"`
	AttachStderr    bool         `json:"AttachStderr"`
	ExposedPorts    struct{}     `json:"ExposedPorts"`
	TTY             bool         `json:"Tty"`
	OpenStdin       bool         `json:"OpenStdin"`
	StdinOnce       bool         `json:"StdinOnce"`
	Env             []string     `json:"Env"`
	Cmd             []string     `json:"Cmd"`
	HealthCheck     HealthConfig `json:"HealthCheck"`
	ArgsEscaped     bool         `json:"ArgsEscaped"`
	Image           string       `json:"Image"`
	Volumes         struct{}     `json:"Volumes"`
	WorkingDir      string       `json:"WorkingDir"`
	EntryPoint      []string     `json:"Entrypoint"`
	NetworkDisabled bool         `json:"NetworkDisabled"`
	MacAddress      string       `json:"MacAddress"`
	OnBuild         []string     `json:"OnBuild"`
	Labels          struct{}     `json:"Labels"`
	StopSignal      string       `json:"StopSignal"`
	StopTimeout     float64      `json:"StopTimeout"`
	Shell           []string     `json:"Shell"`
}

// NetworkConfig is a json field
type NetworkConfig struct {
	Bridge      string  `json:"Bridge"`
	Gateway     string  `json:"Gateway"`
	Address     string  `json:"Address"`
	IPPrefixLen float64 `json:"IPPrefixLen"`
	MacAddress  string  `json:"MacAddress"`
	PortMapping string  `json:"PortMapping"`
	Ports       Port    `json:"Ports"`
}

// Container is a docker api json payload
type Container struct {
	ID              string          `json:"Id"`
	Names           []string        `json:"Names"`
	Image           string          `json:"Image"`
	ImageID         string          `json:"ImageID"`
	Command         string          `json:"Command"`
	Created         float64         `json:"Created"`
	State           string          `json:"State"`
	Status          string          `json:"Status"`
	Ports           []Port          `json:"Ports"`
	Labels          Labels          `json:"Labels"`
	SizeRw          float64         `json:"SizeRw"`
	SizeRootFS      float64         `json:"SizeRootFs"`
	HostConfig      HostConfig      `json:"HostConfig"`
	NetworkSettings NetworkSettings `json:"NetworkSettings"`
	Mounts          []Mounts        `json:"Mounts"`
}

// Image is a docker api json payload
type Image struct {
	ID          string   `json:"Id"`
	ParentID    string   `json:"ParentId"`
	RepoTags    []string `json:"RepoTags"`
	RepoDigests []string `json:"RepoDigests"`
	Created     float64  `json:"Created"`
	Size        float64  `json:"Size"`
	VirtualSize float64  `json:"VirtualSize"`
	SharedSize  float64  `json:"SharedSize"`
	Labels      Labels   `json:"Labels"`
	Containers  float64  `json:"Containers"`
}
