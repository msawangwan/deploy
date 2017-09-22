package ciio

// Buildfile is a json object that is similar to a Dockerfile
type Buildfile struct {
	ContainerName string  `json:"containerName"`
	Image         string  `json:"image"`
	Addr          Addr    `json:"addr"`
	Cmd           Command `json:"cmd"`
}

// // APITemplate is for templating a docker api call
// type APITemplate struct {
// 	Endpoint         string            `json:"endpoint"`
// 	QueryStringPairs map[string]string `json:"queryStringPairs"`
// }
