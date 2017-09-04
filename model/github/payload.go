package github

// Payload is a webhook object
type Payload struct {
	Ref    string
	Before string
	After  string
}
