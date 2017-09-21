// package dock

// // Port is a json field for the containers payload
// type Port struct {
// 	IP          string  `json:"IP"`
// 	PrivatePort float64 `json:"PrivatePort"`
// 	PublicPort  float64 `json:"PublicPort"`
// 	Type        string  `json:"Type"`
// }

// // Labels is a json field for the containers payload
// type Labels struct {
// 	ComExampleVendor  string `json:"com.example.vendor"`
// 	ComExampleLicense string `json:"com.example.license"`
// 	ComExampleVersion string `json:"com.example.version"`
// }

// // HostConfig is a json field for the containers payload
// type HostConfig struct {
// 	NetworkMode string `json:"NetworkMode"`
// }

// // Bridge is a json field for the containers payload
// type Bridge struct {
// 	NetworkID           string  `json:"NetworkID"`
// 	EndpointID          string  `json:"EndpointID"`
// 	Gateway             string  `json:"Gateway"`
// 	IPAddress           string  `json:"IPAddress"`
// 	IPPrefixLen         float64 `json:"IPPrefixLen"`
// 	IPv6Gateway         string  `json:"IPv6Gateway"`
// 	GlobalIPv6Address   string  `json:"GlobalIPv6Address"`
// 	GlobalIPv6PrefixLen float64 `json:"GlobalIPv6PrefixLen"`
// 	MacAddress          string  `json:"MacAddress"`
// }

// // Networks is a json field for the containers payload
// type Networks struct {
// 	Bridge Bridge `json:"bridge"`
// }

// // NetworkSettings is a json field for the containers payload
// type NetworkSettings struct {
// 	Networks Networks `json:"Networks"`
// }

// // Mounts is a json field for the containers payload
// type Mounts struct {
// 	Name        string `json:"Name"`
// 	Source      string `json:"Source"`
// 	Destination string `json:"Destination"`
// 	Driver      string `json:"Driver"`
// 	Mode        string `json:"Mode"`
// 	RW          bool   `json:"RW"`
// 	Propagation string `json:"Propagation"`
// }

// // State is a json field
// type State struct {
// 	Status     string  `json:"Status"`
// 	Running    bool    `json:"Running"`
// 	Paused     bool    `json:"Paused"`
// 	Restarting bool    `json:"Restarting"`
// 	OOMKilled  bool    `json:"OOMKilled"`
// 	Dead       bool    `json:"Dead"`
// 	PID        float64 `json:"Pid"`
// 	ExitCode   float64 `json:"ExitCode"`
// 	Error      string  `json:"Error"`
// 	StartedAt  string  `json:"StartedAt"`
// 	FinishedAt string  `json:"FinishedAt"`
// }

// // GraphDriver is a json field
// type GraphDriver struct {
// 	Name string   `json:"Name"`
// 	Data struct{} `json:"Data"`
// }

// // HealthConfig is a json field
// type HealthConfig struct {
// 	Test        []string `json:"Test"`
// 	Interval    float64  `json:"Interval"`
// 	Timeout     float64  `json:"Timeout"`
// 	Retries     float64  `json:"Retries"`
// 	StartPeriod float64  `json:"StartPeriod"`
// }

// // Config is a json field
// type Config struct {
// 	HostName        string       `json:"Hostname"`
// 	DomainName      string       `json:"Domainname"`
// 	User            string       `json:"User"`
// 	AttachStdin     bool         `json:"AttachStdin"`
// 	AttachStdout    bool         `json:"AttachStdout"`
// 	AttachStderr    bool         `json:"AttachStderr"`
// 	ExposedPorts    struct{}     `json:"ExposedPorts"`
// 	TTY             bool         `json:"Tty"`
// 	OpenStdin       bool         `json:"OpenStdin"`
// 	StdinOnce       bool         `json:"StdinOnce"`
// 	Env             []string     `json:"Env"`
// 	Cmd             []string     `json:"Cmd"`
// 	HealthCheck     HealthConfig `json:"HealthCheck"`
// 	ArgsEscaped     bool         `json:"ArgsEscaped"`
// 	Image           string       `json:"Image"`
// 	Volumes         struct{}     `json:"Volumes"`
// 	WorkingDir      string       `json:"WorkingDir"`
// 	EntryPoint      []string     `json:"Entrypoint"`
// 	NetworkDisabled bool         `json:"NetworkDisabled"`
// 	MacAddress      string       `json:"MacAddress"`
// 	OnBuild         []string     `json:"OnBuild"`
// 	Labels          struct{}     `json:"Labels"`
// 	StopSignal      string       `json:"StopSignal"`
// 	StopTimeout     float64      `json:"StopTimeout"`
// 	Shell           []string     `json:"Shell"`
// }

// // NetworkConfig is a json field
// type NetworkConfig struct {
// 	Bridge      string  `json:"Bridge"`
// 	Gateway     string  `json:"Gateway"`
// 	Address     string  `json:"Address"`
// 	IPPrefixLen float64 `json:"IPPrefixLen"`
// 	MacAddress  string  `json:"MacAddress"`
// 	PortMapping string  `json:"PortMapping"`
// 	Ports       Port    `json:"Ports"`
// }
