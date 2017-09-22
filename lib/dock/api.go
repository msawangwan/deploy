package dock

// Success200 is the schema for 'inspect' success response
// TODO: rename to 'InspectResponse'
type Success200 struct {
	ID              string          `json:"Id"`
	Created         string          `json:"Created"`
	Path            string          `json:"Path"`
	Args            []string        `json:"Args"`
	State           State           `json:"State"`
	Image           string          `json:"Image"`
	ResolvConfPath  string          `json:"ResolvConfPath"`
	HostnamePath    string          `json:"HostnamePath"`
	HostsPath       string          `json:"HostsPath"`
	LogPath         string          `json:"LogPath"`
	Node            struct{}        `json:"Node"`
	Name            string          `json:"Name"`
	RestartCount    float64         `json:"RestartCount"`
	Driver          string          `json:"Driver"`
	MountLabel      string          `json:"MountLabel"`
	ProcessLabel    string          `json:"ProcessLabel"`
	AppArmorProfile string          `json:"AppArmorProfile"`
	ExecIDs         string          `json:"ExecIDs"`
	HostConfig      HostConfig      `json:"HostConfig"`
	GraphDriver     GraphDriverData `json:"GraphDriver"`
	SizeRW          float64         `json:"SizeRw"`
	SizeRootFS      float64         `json:"SizeRootFs"`
	Mounts          []Mount         `json:"Mounts"`
	Config          ContainerConfig `json:"Config"`
	NetworkSettings NetworkConfig   `json:"NetworkSettings"`
}

// Success201 is the schema for 'create' response success
// TODO: Rename To 'CreateResponse'
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

// Mounts is a json field that represents a mountpoint
type Mount struct {
	Name          string        `json:"Name"`
	Target        string        `json:"Target"`
	Source        string        `json:"Source"`
	Type          string        `json:"Type"`
	Destination   string        `json:"Destination"`
	Driver        string        `json:"Driver"`
	Mode          string        `json:"Mode"`
	ReadOnly      bool          `json:"ReadOnly"`
	RW            bool          `json:"RW"`
	Propagation   string        `json:"Propagation"`
	Consistency   string        `json:"Consistency"`
	BindOptions   BindOptions   `json:"BindOptions"`
	VolumeOptions VolumeOptions `json:"VolumeOptions"`
	TMPFSOptions  TMPFSOptions  `json:"TmpfsOptions"`
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
type GraphDriverData struct {
	Name string `json:"Name"`
	Data struct {
		AdditionalProperties `json:"AdditionalProperties"`
	} `json:"Data"`
}

// ContainerConfig is a json field
type ContainerConfig struct {
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

// HostConfig is an payload object that represents the container config depending on the host
type HostConfig struct {
	CPUShares            float64        `json:"CpuShares"`
	Memory               float64        `json:""`
	CGroupParent         string         `json:""`
	BLKIOWeight          float64        `json:""`
	BLKIOWeightDevice    ThrottleDevice `json:""`
	BLKIODeviceReadBPS   ThrottleDevice `json:""`
	BLKIODeviceWriteBPS  ThrottleDevice `json:""`
	BLKIODeviceReadLOps  ThrottleDevice `json:""`
	BLKIODeviceWriteLOps ThrottleDevice `json:""`
	CPUPeriod            float64        `json:""`
	CPUQuota             float64        `json:""`
	CPURealtimePeriod    float64        `json:""`
	CPURealtimeRuntime   float64        `json:""`
	CPUSetCPUs           string         `json:""`
	CPUSetMems           string         `json:""`
	Devices              struct{}       `json:""`
	DeviceCGroupRules    []string       `json:""`
	DiskQuota            float64        `json:""`
	KernelMemory         float64        `json:""`
	MemoryReservation    float64        `json:""`
	MemorySwap           float64        `json:""`
	MemorySwappiness     float64        `json:""`
	NanoCPUs             float64        `json:""`
	OOMKillDisable       bool           `json:""`
	PIDsLimit            float64        `json:""`
	ULimits              struct{}       `json:""`
	CPUCount             float64        `json:""`
	CPUPercent           float64        `json:""`
	IOMaximumBandwidth   float64        `json:""`
	Binds                []string       `json:""`
	ContainerIDFile      string         `json:""`
	LogConfig            struct{}       `json:""`
	NetworkMode          string         `json:""`
	PortBindings         struct{}       `json:""`
	AutoRemove           bool           `json:""`
	VolumeDriver         string         `json:""`
	VolumesFrom          []string       `json:""`
	Mounts               []Mount        `json:""`
	CapAdd               []string       `json:""`
	CapDrop              []string       `json:""`
	DNS                  []string       `json:""`
	DNSOptions           []string       `json:""`
	DNSSearch            []string       `json:""`
	ExtraHosts           []string       `json:""`
	GroupAdd             []string       `json:""`
	IPCMode              string         `json:""`
	CGroup               string         `json:""`
	Links                []string       `json:""`
	OOMScoreAdj          float64        `json:""`
	PIDMode              string         `json:""`
	Privileged           bool           `json:""`
	PublishAllPorts      bool           `json:""`
	ReadOnlyRootFS       bool           `json:""`
	SecurityOpt          []string       `json:""`
	StorageOpt           struct{}       `json:""`
	TMPFS                struct{}       `json:""`
	UTSMode              string         `json:""`
	UserNSMode           string         `json:""`
	SHMSize              float64        `json:""`
	SYSCTLs              struct{}       `json:""`
	Runtime              string         `json:""`
	ConsoleSize          []float64      `json:""`
	Isolation            string         `json:""`
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

// HealthConfig is a json field
type HealthConfig struct {
	Test        []string `json:"Test"`
	Interval    float64  `json:"Interval"`
	Timeout     float64  `json:"Timeout"`
	Retries     float64  `json:"Retries"`
	StartPeriod float64  `json:"StartPeriod"`
}

type LogConfig struct {
	Type   string   `json:"Type"`
	Config struct{} `json:"Config"`
}

type DriverConfig struct {
	Name    string  `json:"Name"`
	Options Options `json:"Options"`
}

// Container is a docker api json payload
// TODO: Rename to 'ListResponse'
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
	Mounts          []Mount         `json:"Mounts"`
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

type ThrottleDevice struct {
	Path string  `json:"Path"`
	Rate float64 `json:"Rate"`
}

type DeviceMapping struct {
	PathOnHost        string `json:"PathOnHost"`
	PathInContainer   string `json:"PathInContainer"`
	CGroupPermissions string `json:"CgroupPermissions"`
}

type ULimits struct {
	Name string  `json:"Name"`
	Soft float64 `json:"Soft"`
	Hard float64 `json:"Hard"`
}

type PortBindings struct {
	HostIP   string `json:"HostIp"`
	HostPort string `json:"HostPort"`
}

type RestartPolicy struct {
	Name              string  `json:"Name"`
	MaximumRetryCount float64 `json:"MaximumRetryCount"`
}

type Options struct {
	AdditionalProperties AdditionalProperties `json:"AdditionalProperties"`
}

type BindOptions struct {
	Propagation interface{} `json:"Propagation"`
}

type VolumeOptions struct {
	NoCopy       bool         `json:"NoCopy"`
	Labels       []Labels     `json:"Labels"`
	DriverConfig DriverConfig `json:"DriverConfig"`
}

type TMPFSOptions struct {
	SizeBytes float64 `json:"SizeBytes"`
	Mode      float64 `json:"Mode"`
}

type StorageOpt string
type TMPFS string
type SYSCTLs string
type AdditionalProperties string
