package webhook

import "github.com/msawangwan/ci.io/lib/internal/webhook/payload"

// PushEvent describes a github webhook json data structure
type PushEvent struct {
	Ref        string             `json:"ref"`
	Before     string             `json:"before"`
	After      string             `json:"after"`
	Created    bool               `json:"created"`
	Deleted    bool               `json:"deleted"`
	Forced     bool               `json:"forced"`
	BaseRef    string             `json:"base_ref"`
	Compare    string             `json:"compare"`
	Commits    []payload.Commit   `json:"commits"`
	HeadCommit payload.Commit     `json:"head_commit"`
	Repository payload.Repository `json:"repository"`
	Pusher     payload.Pusher     `json:"pusher"`
	Sender     payload.Sender     `json:"sender"`
}
