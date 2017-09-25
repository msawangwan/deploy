package dock

// ListResponse is a docker api json payload
type ListResponse struct {
	ID              string          `json:"Id,omitempty"`
	Names           []string        `json:"Names,omitempty"`
	Image           string          `json:"Image,omitempty"`
	ImageID         string          `json:"ImageID,omitempty"`
	Command         string          `json:"Command,omitempty"`
	Created         float64         `json:"Created,omitempty"`
	State           string          `json:"State,omitempty"`
	Status          string          `json:"Status,omitempty"`
	Ports           []Port          `json:"Ports,omitempty"`
	Labels          Labels          `json:"Labels,omitempty"`
	SizeRw          float64         `json:"SizeRw,omitempty"`
	SizeRootFS      float64         `json:"SizeRootFs,omitempty"`
	HostConfig      HostConfig      `json:"HostConfig,omitempty"`
	NetworkSettings NetworkSettings `json:"NetworkSettings,omitempty"`
	Mounts          []Mount         `json:"Mounts,omitempty"`
	Message         string          `json:"message,omitempty"`
}

// InspectResponse is the schema for docker api command 'inspect' on 200 ok
type InspectResponse struct {
	ID              string          `json:"Id,omitempty"`
	Created         string          `json:"Created,omitempty"`
	Path            string          `json:"Path,omitempty"`
	Args            []string        `json:"Args,omitempty"`
	State           State           `json:"State,omitempty"`
	Image           string          `json:"Image,omitempty"`
	ResolvConfPath  string          `json:"ResolvConfPath,omitempty"`
	HostnamePath    string          `json:"HostnamePath,omitempty"`
	HostsPath       string          `json:"HostsPath,omitempty"`
	LogPath         string          `json:"LogPath,omitempty"`
	Node            struct{}        `json:"Node,omitempty"`
	Name            string          `json:"Name,omitempty"`
	RestartCount    float64         `json:"RestartCount,omitempty"`
	Driver          string          `json:"Driver,omitempty"`
	MountLabel      string          `json:"MountLabel,omitempty"`
	ProcessLabel    string          `json:"ProcessLabel,omitempty"`
	AppArmorProfile string          `json:"AppArmorProfile,omitempty"`
	ExecIDs         string          `json:"ExecIDs,omitempty"`
	HostConfig      HostConfig      `json:"HostConfig,omitempty"`
	GraphDriver     GraphDriverData `json:"GraphDriver,omitempty"`
	SizeRW          float64         `json:"SizeRw,omitempty"`
	SizeRootFS      float64         `json:"SizeRootFs,omitempty"`
	Mounts          []Mount         `json:"Mounts,omitempty"`
	Config          ContainerConfig `json:"Config,omitempty"`
	NetworkSettings NetworkConfig   `json:"NetworkSettings,omitempty"`
	Message         string          `json:"message,omitempty"`
}

// CreateRequest is the request body to create a container
type CreateRequest struct {
	HostName         string           `json:"Hostname,omitempty"`
	DomainName       string           `json:"Domainname,omitempty"`
	User             string           `json:"User,omitempty"`
	AttachStdin      bool             `json:"AttachStdin,omitempty"`
	AttachStdout     bool             `json:"AttachStdout,omitempty"`
	AttachStderr     bool             `json:"AttachStderr,omitempty"`
	ExposedPorts     struct{}         `json:"ExposedPorts,omitempty"`
	TTY              bool             `json:"Tty,omitempty"`
	OpenStdin        bool             `json:"OpenStdin,omitempty"`
	StdinOnce        bool             `json:"StdinOnce,omitempty"`
	Env              []string         `json:"Env,omitempty"`
	Cmd              []string         `json:"Cmd,omitempty"`
	HealthCheck      HealthConfig     `json:"HealthCheck,omitempty"`
	ArgsEscaped      bool             `json:"ArgsEscaped,omitempty"`
	Image            string           `json:"Image,omitempty"`
	Volumes          struct{}         `json:"Volumes,omitempty"`
	WorkingDir       string           `json:"WorkingDir,omitempty"`
	EntryPoint       []string         `json:"Entrypoint,omitempty"`
	NetworkDisabled  bool             `json:"NetworkDisabled,omitempty"`
	MacAddress       string           `json:"MacAddress,omitempty"`
	OnBuild          []string         `json:"OnBuild,omitempty"`
	Labels           struct{}         `json:"Labels,omitempty"`
	StopSignal       string           `json:"StopSignal,omitempty"`
	StopTimeout      float64          `json:"StopTimeout,omitempty"`
	Shell            []string         `json:"Shell,omitempty"`
	HostConfig       HostConfig       `json:"HostConfig,omitempty"`
	NetworkingConfig NetworkingConfig `json:"NetworkingConfig,omitempty"`
}

// CreateResponse is the schema for 'create' response success
type CreateResponse struct {
	ID       string   `json:"Id,omitempty"`
	Warnings []string `json:"Warnings,omitempty"`
	Message  string   `json:"message,omitempty"`
}

// StartResponse is for creating container
type StartResponse struct {
	Message string `json:"message,omitempty"`
}

// EmptyResponse is a no error success response
type EmptyResponse struct{}

// Image is a docker api json payload
type Image struct {
	ID          string   `json:"Id,omitempty"`
	ParentID    string   `json:"ParentId,omitempty"`
	RepoTags    []string `json:"RepoTags,omitempty"`
	RepoDigests []string `json:"RepoDigests,omitempty"`
	Created     float64  `json:"Created,omitempty"`
	Size        float64  `json:"Size,omitempty"`
	VirtualSize float64  `json:"VirtualSize,omitempty"`
	SharedSize  float64  `json:"SharedSize,omitempty"`
	Labels      Labels   `json:"Labels,omitempty"`
	Containers  float64  `json:"Containers,omitempty"`
}

// // Error400 is a bad parameter error
// type Error400 struct {
// 	Message string `json:"message"`
// }

// // Error404 is a no such container error
// type Error404 struct {
// 	Message string `json:"message"`
// }

// // Error409 is a conflict error
// type Error409 struct {
// 	Message string `json:"message"`
// }

// // Error500 is an internal server error
// type Error500 struct {
// 	Message string `json:"message"`
// }

// Port is a json field for the containers payload
type Port struct {
	IP          string  `json:"IP,omitempty"`
	PrivatePort float64 `json:"PrivatePort,omitempty"`
	PublicPort  float64 `json:"PublicPort,omitempty"`
	Type        string  `json:"Type,omitempty"`
}

// Labels is a json field for the containers payload
type Labels struct {
	ComExampleVendor  string `json:"com.example.vendor"`
	ComExampleLicense string `json:"com.example.license"`
	ComExampleVersion string `json:"com.example.version"`
}

// Bridge is a json field for the containers payload
type Bridge struct {
	NetworkID           string  `json:"NetworkID,omitempty"`
	EndpointID          string  `json:"EndpointID,omitempty"`
	Gateway             string  `json:"Gateway,omitempty"`
	IPAddress           string  `json:"IPAddress,omitempty"`
	IPPrefixLen         float64 `json:"IPPrefixLen,omitempty"`
	IPv6Gateway         string  `json:"IPv6Gateway,omitempty"`
	GlobalIPv6Address   string  `json:"GlobalIPv6Address,omitempty"`
	GlobalIPv6PrefixLen float64 `json:"GlobalIPv6PrefixLen,omitempty"`
	MacAddress          string  `json:"MacAddress,omitempty"`
}

// Networks is a json field for the containers payload
type Networks struct {
	Bridge Bridge `json:"bridge,omitempty"`
}

// NetworkSettings is a json field for the containers payload
type NetworkSettings struct {
	Networks Networks `json:"Networks,omitempty"`
}

// EndpointSettings is a json field
type EndpointSettings struct {
	IPAMConfig          IPAMConfig `json:"IPAMConfig,omitempty"`
	Links               []string   `json:"Links,omitempty"`
	Aliases             []string   `json:"Aliases,omitempty"`
	NetworkID           string     `json:"NetworkID,omitempty"`
	EndpointID          string     `json:"EndpointID,omitempty"`
	Gateway             string     `json:"Gateway,omitempty"`
	IPAddress           string     `json:"IpAddress,omitempty"`
	IPPrefixLen         float64    `json:"IpPrefixLen,omitempty"`
	IPv6Gateway         string     `json:"IPv6Gateway,omitempty"`
	GlobalIPv6Address   string     `json:"GlobalIPv6Address,omitempty"`
	GlobalIPv6PrefixLen float64    `json:"GlobalIPv6PrefixLen,omitempty"`
	MacAddress          string     `json:"MacAddress,omitempty"`
}

// Mount is a json field that represents a mountpoint
type Mount struct {
	Name          string        `json:"Name,omitempty"`
	Target        string        `json:"Target,omitempty"`
	Source        string        `json:"Source,omitempty"`
	Type          string        `json:"Type,omitempty"`
	Destination   string        `json:"Destination,omitempty"`
	Driver        string        `json:"Driver,omitempty"`
	Mode          string        `json:"Mode,omitempty"`
	ReadOnly      bool          `json:"ReadOnly,omitempty"`
	RW            bool          `json:"RW,omitempty"`
	Propagation   string        `json:"Propagation,omitempty"`
	Consistency   string        `json:"Consistency,omitempty"`
	BindOptions   BindOptions   `json:"BindOptions,omitempty"`
	VolumeOptions VolumeOptions `json:"VolumeOptions,omitempty"`
	TMPFSOptions  TMPFSOptions  `json:"TmpfsOptions,omitempty"`
}

// State is a json field
type State struct {
	Status     string  `json:"Status,omitempty"`
	Running    bool    `json:"Running,omitempty"`
	Paused     bool    `json:"Paused,omitempty"`
	Restarting bool    `json:"Restarting,omitempty"`
	OOMKilled  bool    `json:"OOMKilled,omitempty"`
	Dead       bool    `json:"Dead,omitempty"`
	PID        float64 `json:"Pid,omitempty"`
	ExitCode   float64 `json:"ExitCode,omitempty"`
	Error      string  `json:"Error,omitempty"`
	StartedAt  string  `json:"StartedAt,omitempty"`
	FinishedAt string  `json:"FinishedAt,omitempty"`
}

// GraphDriverData is a json field
type GraphDriverData struct {
	Name string `json:"Name,omitempty"`
	Data struct {
		AdditionalProperties `json:"AdditionalProperties,omitempty"`
	} `json:"Data,omitempty"`
}

// ContainerConfig is a json field
type ContainerConfig struct {
	HostName        string       `json:"Hostname,omitempty"`
	DomainName      string       `json:"Domainname,omitempty"`
	User            string       `json:"User,omitempty"`
	AttachStdin     bool         `json:"AttachStdin,omitempty"`
	AttachStdout    bool         `json:"AttachStdout,omitempty"`
	AttachStderr    bool         `json:"AttachStderr,omitempty"`
	ExposedPorts    struct{}     `json:"ExposedPorts,omitempty"`
	TTY             bool         `json:"Tty,omitempty"`
	OpenStdin       bool         `json:"OpenStdin,omitempty"`
	StdinOnce       bool         `json:"StdinOnce,omitempty"`
	Env             []string     `json:"Env,omitempty"`
	Cmd             []string     `json:"Cmd,omitempty"`
	HealthCheck     HealthConfig `json:"HealthCheck,omitempty"`
	ArgsEscaped     bool         `json:"ArgsEscaped,omitempty"`
	Image           string       `json:"Image,omitempty"`
	Volumes         struct{}     `json:"Volumes,omitempty"`
	WorkingDir      string       `json:"WorkingDir,omitempty"`
	EntryPoint      []string     `json:"Entrypoint,omitempty"`
	NetworkDisabled bool         `json:"NetworkDisabled,omitempty"`
	MacAddress      string       `json:"MacAddress,omitempty"`
	OnBuild         []string     `json:"OnBuild,omitempty"`
	Labels          struct{}     `json:"Labels,omitempty"`
	StopSignal      string       `json:"StopSignal,omitempty"`
	StopTimeout     float64      `json:"StopTimeout,omitempty"`
	Shell           []string     `json:"Shell,omitempty"`
}

// HostConfig is an payload object that represents the container config depending on the host
type HostConfig struct {
	CPUShares            float64          `json:"CpuShares,omitempty"`
	Memory               float64          `json:"Memory,omitempty"`
	CGroupParent         string           `json:"CgroupParent,omitempty"`
	BLKIOWeight          float64          `json:"BlkioWeight,omitempty"`
	BLKIOWeightDevice    []ThrottleDevice `json:"BlkioWeightDevice,omitempty"`
	BLKIODeviceReadBPS   []ThrottleDevice `json:"BlkioDeviceReadBPS,omitempty"`
	BLKIODeviceWriteBPS  []ThrottleDevice `json:"BlkioDeviceWriteBPS,omitempty"`
	BLKIODeviceReadLOps  []ThrottleDevice `json:"BlkioDeviceReadLOps,omitempty"`
	BLKIODeviceWriteLOps []ThrottleDevice `json:"BlkioDeviceWriteLOps,omitempty"`
	CPUPeriod            float64          `json:"CpuPeriod,omitempty"`
	CPUQuota             float64          `json:"CpuQuota,omitempty"`
	CPURealtimePeriod    float64          `json:"CpuRealtimePeriod,omitempty"`
	CPURealtimeRuntime   float64          `json:"CpuRealtimeRuntime,omitempty"`
	CPUSetCPUs           string           `json:"CpusetCPUs,omitempty"`
	CPUSetMems           string           `json:"CpuSetMems,omitempty"`
	Devices              struct{}         `json:"Devices,omitempty"`
	DeviceCGroupRules    []string         `json:"DeviceCgroupRules,omitempty"`
	DiskQuota            float64          `json:"DiskQuota,omitempty"`
	KernelMemory         float64          `json:"KernelMemory,omitempty"`
	MemoryReservation    float64          `json:"MemoryReservation,omitempty"`
	MemorySwap           float64          `json:"MemorySwap,omitempty"`
	MemorySwappiness     float64          `json:"MemorySwappiness,omitempty"`
	NanoCPUs             float64          `json:"NanoCPUs,omitempty"`
	OOMKillDisable       bool             `json:"OomKillDisable,omitempty"`
	PIDsLimit            float64          `json:"PidsLimit,omitempty"`
	ULimits              struct{}         `json:"Ulimits,omitempty"`
	CPUCount             float64          `json:"CpuCount,omitempty"`
	CPUPercent           float64          `json:"CpuPercent,omitempty"`
	IOMaximumIOps        float64          `json:"IOMaximumIops,omitempty"`
	IOMaximumBandwidth   float64          `json:"IOMaximumBandwidth,omitempty"`
	Binds                []string         `json:"Binds,omitempty"`
	ContainerIDFile      string           `json:"ContainerIDFile,omitempty"`
	LogConfig            struct{}         `json:"LogConfig,omitempty"`
	NetworkMode          string           `json:"NetworkMode,omitempty"`
	PortBindings         struct{}         `json:"PortBindings,omitempty"`
	RestartPolicy        RestartPolicy    `json:"RestartPolicy,omitempty"`
	AutoRemove           bool             `json:"AutoRemove,omitempty"`
	VolumeDriver         string           `json:"VolumeDriver,omitempty"`
	VolumesFrom          []string         `json:"VolumesFrom,omitempty"`
	Mounts               []Mount          `json:"Mounts,omitempty"`
	CapAdd               []string         `json:"CapAdd,omitempty"`
	CapDrop              []string         `json:"CapDrop,omitempty"`
	DNS                  []string         `json:"Dns,omitempty"`
	DNSOptions           []string         `json:"DnsOptions,omitempty"`
	DNSSearch            []string         `json:"DnsSearch,omitempty"`
	ExtraHosts           []string         `json:"ExtraHosts,omitempty"`
	GroupAdd             []string         `json:"GroupAdd,omitempty"`
	IPCMode              string           `json:"IpcMode,omitempty"`
	CGroup               string           `json:"Cgroup,omitempty"`
	Links                []string         `json:"Links,omitempty"`
	OOMScoreAdj          float64          `json:"OomScoreAdj,omitempty"`
	PIDMode              string           `json:"PidMode,omitempty"`
	Privileged           bool             `json:"Privileged,omitempty"`
	PublishAllPorts      bool             `json:"PublishAllPorts,omitempty"`
	ReadOnlyRootFS       bool             `json:"ReadOnlyRootFS,omitempty"`
	SecurityOpt          []string         `json:"SecurityOpt,omitempty"`
	StorageOpt           struct{}         `json:"StorageOpt,omitempty"`
	TMPFS                struct{}         `json:"Tmpfs,omitempty"`
	UTSMode              string           `json:"UTSMode,omitempty"`
	UserNSMode           string           `json:"UsernsMode,omitempty"`
	SHMSize              float64          `json:"ShmSize,omitempty"`
	SYSCTLs              struct{}         `json:"Sysctls,omitempty"`
	Runtime              string           `json:"Runtime,omitempty"`
	ConsoleSize          []float64        `json:"ConsoleSize,omitempty"`
	Isolation            string           `json:"Isolation,omitempty"`
}

// NetworkConfig is a json field
type NetworkConfig struct {
	Bridge      string  `json:"Bridge,omitempty"`
	Gateway     string  `json:"Gateway,omitempty"`
	Address     string  `json:"Address,omitempty"`
	IPPrefixLen float64 `json:"IPPrefixLen,omitempty"`
	MacAddress  string  `json:"MacAddress,omitempty"`
	PortMapping string  `json:"PortMapping,omitempty"`
	Ports       Port    `json:"Ports,omitempty"`
}

// NetworkingConfig is a json field
type NetworkingConfig struct {
	EndpointsConfig EndpointsConfig `json:"EndpointsConfig,omitempty"`
}

// EndpointsConfig is a json field
type EndpointsConfig struct {
	AdditionalProperties EndpointSettings `json:"additionalProperties,omitempty"`
}

// IPAMConfig is a json field
type IPAMConfig struct {
	IPv4Address  string   `json:"IPv4Address,omitempty"`
	IPv6Address  string   `json:"IPv6Address,omitempty"`
	LinkLocalIPs []string `json:"LinkLocalIPs,omitempty"`
}

// HealthConfig is a json field in ContainerConfig
type HealthConfig struct {
	Test        []string `json:"Test,omitempty"`
	Interval    float64  `json:"Interval,omitempty"`
	Timeout     float64  `json:"Timeout,omitempty"`
	Retries     float64  `json:"Retries,omitempty"`
	StartPeriod float64  `json:"StartPeriod,omitempty"`
}

// LogConfig is a json field in HostConfig
type LogConfig struct {
	Type   string   `json:"Type,omitempty"`
	Config struct{} `json:"Config,omitempty"`
}

// DriverConfig is a json field in DriverOptions
type DriverConfig struct {
	Name    string  `json:"Name,omitempty"`
	Options Options `json:"Options,omitempty"`
}

// ThrottleDevice is a json field in HostConfig
type ThrottleDevice struct {
	Path string  `json:"Path,omitempty"`
	Rate float64 `json:"Rate,omitempty"`
}

// DeviceMapping is a json field in HostConfig
type DeviceMapping struct {
	PathOnHost        string `json:"PathOnHost,omitempty"`
	PathInContainer   string `json:"PathInContainer,omitempty"`
	CGroupPermissions string `json:"CgroupPermissions,omitempty"`
}

// ULimits is a json field in HostConfig
type ULimits struct {
	Name string  `json:"Name,omitempty"`
	Soft float64 `json:"Soft,omitempty"`
	Hard float64 `json:"Hard,omitempty"`
}

// PortBindings is a json field in HostConfig
type PortBindings struct {
	HostIP   string `json:"HostIp,omitempty"`
	HostPort string `json:"HostPort,omitempty"`
}

// RestartPolicy is a json field in HostConfig
type RestartPolicy struct {
	Name              string  `json:"Name,omitempty"`
	MaximumRetryCount float64 `json:"MaximumRetryCount,omitempty"`
}

// Options is a json field in HostConfig
type Options struct {
	AdditionalProperties AdditionalProperties `json:"AdditionalProperties,omitempty"`
}

// BindOptions is a json field in HostConfig
type BindOptions struct {
	Propagation interface{} `json:"Propagation,omitempty"`
}

// VolumeOptions is a json field in HostConfig
type VolumeOptions struct {
	NoCopy       bool         `json:"NoCopy,omitempty"`
	Labels       []Labels     `json:"Labels,omitempty"`
	DriverConfig DriverConfig `json:"DriverConfig,omitempty"`
}

// TMPFSOptions is a json field in HostConfig
type TMPFSOptions struct {
	SizeBytes float64 `json:"SizeBytes,omitempty"`
	Mode      float64 `json:"Mode,omitempty"`
}

// StorageOpt is a json field in HostConfig
type StorageOpt string

// TMPFS is a json field in HostConfig
type TMPFS string

// SYSCTLs is a json field in HostConfig
type SYSCTLs string

// AdditionalProperties is a json field
type AdditionalProperties string
