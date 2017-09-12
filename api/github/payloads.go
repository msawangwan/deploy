package github

// PushEvent describes a github webhook json data structure
type PushEvent struct {
	Ref        string     `json:"ref"`
	Before     string     `json:"before"`
	After      string     `json:"after"`
	Created    bool       `json:"created"`
	Deleted    bool       `json:"deleted"`
	Forced     bool       `json:"forced"`
	BaseRef    string     `json:"base_ref"`
	Compare    string     `json:"compare"`
	Commits    []Commit   `json:"commits"`
	HeadCommit Commit     `json:"head_commit"`
	Repository Repository `json:"repository"`
	Pusher     Pusher     `json:"pusher"`
	Sender     Sender     `json:"sender"`
}
