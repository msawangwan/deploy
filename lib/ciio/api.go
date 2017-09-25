package ciio

// Buildfile is a json object that is similar to a Dockerfile
type Buildfile struct {
	ContainerName     string            `json:"containerName,omitempty"`
	Image             string            `json:"image,omitempty"`
	WorkingDir        string            `json:"workingDir,omitempty"`
	Cmd               Command           `json:"cmd,omitempty"`
	EntryPoint        Command           `json:"entryPoint,omitempty"`
	NetworkParameters NetworkParameters `json:"networkParameters,omitempty"`
	Stages            []Stage           `json:"stages,omitempty"`
}

// NetworkParameters are important network settings
type NetworkParameters struct {
	IP      string `json:"ip,omitempty"`
	PortOut string `json:"portOut,omitempty"`
	PortIn  string `json:"portIn,omitempty"`
}

// Stage represents the commands to run during a stage
type Stage struct {
	Label    string    `json:"label,omitempty"`
	Commands []Command `json:"commands,omitempty"`
}

// Command is a command and an array of args to pass into it
type Command struct {
	Exec string   `json:"exec,omitempty"`
	Args []string `json:"args,omitempty"`
}