package webhook

import "github.com/msawangwan/ci.io/model/webhook/key"

// PushEvent describes a github webhook json data structure
type PushEvent struct {
	Ref        string         `json:"ref"`
	Before     string         `json:"before"`
	After      string         `json:"after"`
	Created    bool           `json:"created"`
	Deleted    bool           `json:"deleted"`
	Forced     bool           `json:"forced"`
	BaseRef    string         `json:"base_ref"`
	Compare    string         `json:"compare"`
	Commits    []key.Commit   `json:"commits"`
	HeadCommit key.Commit     `json:"head_commit"`
	Repository key.Repository `json:"repository"`
	Pusher     key.Pusher     `json:"pusher"`
	Sender     key.Sender     `json:"sender"`
}
