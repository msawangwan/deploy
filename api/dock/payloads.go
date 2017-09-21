package dock

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
