package dock

// InspectResponse is the schema for docker api command 'inspect' on 200 ok
type InspectResponse struct {
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

// CreateResponse is the schema for 'create' response success
type CreateResponse struct {
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

// Mount is a json field that represents a mountpoint
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

// GraphDriverData is a json field
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
	CPUShares            float64        `json:"CpuShares,omitempty"`
	Memory               float64        `json:"Memory,omitempty"`
	CGroupParent         string         `json:"CgroupParent,omitempty"`
	BLKIOWeight          float64        `json:"BlkioWeight,omitempty"`
	BLKIOWeightDevice    ThrottleDevice `json:"BlkioWeightDevice,omitempty"`
	BLKIODeviceReadBPS   ThrottleDevice `json:"BlkioDeviceReadBPS,omitempty"`
	BLKIODeviceWriteBPS  ThrottleDevice `json:"BlkioDeviceWriteBPS,omitempty"`
	BLKIODeviceReadLOps  ThrottleDevice `json:"BlkioDeviceReadLOps,omitempty"`
	BLKIODeviceWriteLOps ThrottleDevice `json:"BlkioDeviceWriteLOps,omitempty"`
	CPUPeriod            float64        `json:"CpuPeriod,omitempty"`
	CPUQuota             float64        `json:"CpuQuota,omitempty"`
	CPURealtimePeriod    float64        `json:"CpuRealtimePeriod,omitempty"`
	CPURealtimeRuntime   float64        `json:"CpuRealtimeRuntime,omitempty"`
	CPUSetCPUs           string         `json:"CpusetCPUs,omitempty"`
	CPUSetMems           string         `json:"CpuSetMems,omitempty"`
	Devices              struct{}       `json:"Devices,omitempty"`
	DeviceCGroupRules    []string       `json:"DeviceCgroupRules,omitempty"`
	DiskQuota            float64        `json:"DiskQuota,omitempty"`
	KernelMemory         float64        `json:"KernelMemory,omitempty"`
	MemoryReservation    float64        `json:"MemoryReservation,omitempty"`
	MemorySwap           float64        `json:"MemorySwap,omitempty"`
	MemorySwappiness     float64        `json:"MemorySwappiness,omitempty"`
	NanoCPUs             float64        `json:"NanoCPUs,omitempty"`
	OOMKillDisable       bool           `json:"OomKillDisable,omitempty"`
	PIDsLimit            float64        `json:"PidsLimit,omitempty"`
	ULimits              struct{}       `json:"Ulimits,omitempty"`
	CPUCount             float64        `json:"CpuCount,omitempty"`
	CPUPercent           float64        `json:"CpuPercent,omitempty"`
	IOMaximumIOps        float64        `json:"IOMaximumIops,omitempty"`
	IOMaximumBandwidth   float64        `json:"IOMaximumBandwidth,omitempty"`
	Binds                []string       `json:"Binds,omitempty"`
	ContainerIDFile      string         `json:"ContainerIDFile,omitempty"`
	LogConfig            struct{}       `json:"LogConfig,omitempty"`
	NetworkMode          string         `json:"NetworkMode,omitempty"`
	PortBindings         struct{}       `json:"PortBindings,omitempty"`
	RestartPolicy        RestartPolicy  `json:"RestartPolicy,omitempty"`
	AutoRemove           bool           `json:"AutoRemove,omitempty"`
	VolumeDriver         string         `json:"VolumeDriver,omitempty"`
	VolumesFrom          []string       `json:"VolumesFrom,omitempty"`
	Mounts               []Mount        `json:"Mounts,omitempty"`
	CapAdd               []string       `json:"CapAdd,omitempty"`
	CapDrop              []string       `json:"CapDrop,omitempty"`
	DNS                  []string       `json:"Dns,omitempty"`
	DNSOptions           []string       `json:"DnsOptions,omitempty"`
	DNSSearch            []string       `json:"DnsSearch,omitempty"`
	ExtraHosts           []string       `json:"ExtraHosts,omitempty"`
	GroupAdd             []string       `json:"GroupAdd,omitempty"`
	IPCMode              string         `json:"IpcMode,omitempty"`
	CGroup               string         `json:"Cgroup,omitempty"`
	Links                []string       `json:"Links,omitempty"`
	OOMScoreAdj          float64        `json:"OomScoreAdj,omitempty"`
	PIDMode              string         `json:"PidMode,omitempty"`
	Privileged           bool           `json:"Privileged,omitempty"`
	PublishAllPorts      bool           `json:"PublishAllPorts,omitempty"`
	ReadOnlyRootFS       bool           `json:"ReadOnlyRootFS,omitempty"`
	SecurityOpt          []string       `json:"SecurityOpt,omitempty"`
	StorageOpt           struct{}       `json:"StorageOpt,omitempty"`
	TMPFS                struct{}       `json:"Tmpfs,omitempty"`
	UTSMode              string         `json:"UTSMode,omitempty"`
	UserNSMode           string         `json:"UsernsMode,omitempty"`
	SHMSize              float64        `json:"ShmSize,omitempty"`
	SYSCTLs              struct{}       `json:"Sysctls,omitempty"`
	Runtime              string         `json:"Runtime,omitempty"`
	ConsoleSize          []float64      `json:"ConsoleSize,omitempty"`
	Isolation            string         `json:"Isolation,omitempty"`
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

// HealthConfig is a json field in ContainerConfig
type HealthConfig struct {
	Test        []string `json:"Test"`
	Interval    float64  `json:"Interval"`
	Timeout     float64  `json:"Timeout"`
	Retries     float64  `json:"Retries"`
	StartPeriod float64  `json:"StartPeriod"`
}

// LogConfig is a json field in HostConfig
type LogConfig struct {
	Type   string   `json:"Type"`
	Config struct{} `json:"Config"`
}

// DriverConfig is a json field in DriverOptions
type DriverConfig struct {
	Name    string  `json:"Name"`
	Options Options `json:"Options"`
}

// ListResponse is a docker api json payload
type ListResponse struct {
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

// ThrottleDevice is a json field in HostConfig
type ThrottleDevice struct {
	Path string  `json:"Path"`
	Rate float64 `json:"Rate"`
}

// DeviceMapping is a json field in HostConfig
type DeviceMapping struct {
	PathOnHost        string `json:"PathOnHost"`
	PathInContainer   string `json:"PathInContainer"`
	CGroupPermissions string `json:"CgroupPermissions"`
}

// ULimits is a json field in HostConfig
type ULimits struct {
	Name string  `json:"Name"`
	Soft float64 `json:"Soft"`
	Hard float64 `json:"Hard"`
}

// PortBindings is a json field in HostConfig
type PortBindings struct {
	HostIP   string `json:"HostIp"`
	HostPort string `json:"HostPort"`
}

// RestartPolicy is a json field in HostConfig
type RestartPolicy struct {
	Name              string  `json:"Name"`
	MaximumRetryCount float64 `json:"MaximumRetryCount"`
}

// Options is a json field in HostConfig
type Options struct {
	AdditionalProperties AdditionalProperties `json:"AdditionalProperties"`
}

// BindOptions is a json field in HostConfig
type BindOptions struct {
	Propagation interface{} `json:"Propagation"`
}

// VolumeOptions is a json field in HostConfig
type VolumeOptions struct {
	NoCopy       bool         `json:"NoCopy"`
	Labels       []Labels     `json:"Labels"`
	DriverConfig DriverConfig `json:"DriverConfig"`
}

// TMPFSOptions is a json field in HostConfig
type TMPFSOptions struct {
	SizeBytes float64 `json:"SizeBytes"`
	Mode      float64 `json:"Mode"`
}

// StorageOpt is a json field in HostConfig
type StorageOpt string

// TMPFS is a json field in HostConfig
type TMPFS string

// SYSCTLs is a json field in HostConfig
type SYSCTLs string

// AdditionalProperties is a json field
type AdditionalProperties string
