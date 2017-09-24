package github

// Owner is a github webhook object
type Owner struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Author is a github webhook object
type Author struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

// Committer is a github webhook object
type Committer struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

// Commit is a webhook object
type Commit struct {
	ID        string    `json:"id"`
	TreeID    string    `json:"tree_id"`
	Distinct  bool      `json:"distinct"`
	Message   string    `json:"message"`
	Timestamp string    `json:"timestamp"`
	URL       string    `json:"url"`
	Author    Author    `json:"author"`
	Committer Committer `json:"committer"`
	Added     []string  `json:"added"`
	Removed   []string  `json:"removed"`
	Modified  []string  `json:"modified"`
}

// Repository is a github webhook object
type Repository struct {
	ID               float64 `json:"id"`
	Name             string  `json:"name"`
	FullName         string  `json:"full_name"`
	Owner            Owner   `json:"owner"`
	Private          bool    `json:"private"`
	HTMLURL          string  `json:"html_url"`
	Description      string  `json:"description"`
	Fork             bool    `json:"fork"`
	URL              string  `json:"url"`
	ForksURL         string  `json:"forks_url"`
	KeysURL          string  `json:"keys_url"`
	CollaboratorsURL string  `json:"collaborators_url"`
	TeamsURL         string  `json:"teams_url"`
	HooksURL         string  `json:"hooks_url"`
	IssueEventsURL   string  `json:"issue_events_url"`
	EventsURL        string  `json:"events_url"`
	AssigneesURL     string  `json:"assignees_url"`
	BranchesURL      string  `json:"branches_url"`
	TagsURL          string  `json:"tags_url"`
	BlobsURL         string  `json:"blobs_url"`
	GitTagsURL       string  `json:"git_tags_url"`
	GitRefsURL       string  `json:"git_refs_url"`
	TreesURL         string  `json:"trees_url"`
	StatusesURL      string  `json:"statuses_url"`
	LanguagesURL     string  `json:"languages_url"`
	StargazersURL    string  `json:"stargazers_url"`
	ContributorsURL  string  `json:"contributors_url"`
	SubscribersURL   string  `json:"subscribers_url"`
	SubscriptionURL  string  `json:"subscription_url"`
	CommitsURL       string  `json:"commits_url"`
	GitCommitsURL    string  `json:"git_commits_url"`
	CommentsURL      string  `json:"comments_url"`
	IssueCommentURL  string  `json:"issue_comment_url"`
	ContentsURL      string  `json:"contents_url"`
	CompareURL       string  `json:"compare_url"`
	MergesURL        string  `json:"merges_url"`
	ArchiveURL       string  `json:"archive_url"`
	DownloadsURL     string  `json:"downloads_url"`
	IssuesURL        string  `json:"issues_url"`
	PullsURL         string  `json:"pulls_url"`
	MilestonesURL    string  `json:"milestones_url"`
	NotificationsURL string  `json:"notifications_url"`
	LabelsURL        string  `json:"labels_url"`
	ReleasesURL      string  `json:"releases_url"`
	CreatedAt        float64 `json:"created_at"`
	UpdatedAt        string  `json:"updated_at"`
	PushedAt         float64 `json:"pushed_at"`
	GitURL           string  `json:"git_url"`
	SSHURL           string  `json:"ssh_url"`
	CloneURL         string  `json:"clone_url"`
	SVNURL           string  `json:"svn_url"`
	Homepage         string  `json:"homepage"`
	Size             float64 `json:"size"`
	StargazersCount  float64 `json:"stargazers_count"`
	WatchersCount    float64 `json:"watchers_count"`
	Language         string  `json:"language"`
	HasIssues        bool    `json:"has_issues"`
	HasDownloads     bool    `json:"has_downloads"`
	HasWiki          bool    `json:"has_wiki"`
	HasPages         bool    `json:"has_pages"`
	ForksCount       float64 `json:"forks_count"`
	MirrorURL        string  `json:"mirror_url"`
	OpenIssuesCount  float64 `json:"open_issues_count"`
	Forks            float64 `json:"forks"`
	OpenIssues       float64 `json:"open_issues"`
	Watchers         float64 `json:"watchers"`
	DefaultBranch    string  `json:"default_branch"`
	Stargazers       float64 `json:"stargazers"`
	MasterBranch     string  `json:"master_branch"`
}

// Pusher is a github webhook object
type Pusher struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Sender is a github webhook object
type Sender struct {
	Login             string  `json:"login"`
	ID                float64 `json:"id"`
	AvatarURL         string  `json:"avatar_url"`
	GravatarID        string  `json:"gravatar_id"`
	URL               string  `json:"url"`
	HTMLURL           string  `json:"html_url"`
	FollowersURL      string  `json:"followers_url"`
	FollowingURL      string  `json:"following_url"`
	GistsURL          string  `json:"gists_url"`
	StarredURL        string  `json:"starred_url"`
	SubscriptionsURL  string  `json:"subscriptions_url"`
	OrganizationsURL  string  `json:"organizations_url"`
	ReposURL          string  `json:"repos_url"`
	EventsURL         string  `json:"events_url"`
	ReceivedEventsURL string  `json:"received_events_url"`
	Type              string  `json:"type"`
	SiteAdmin         bool    `json:"site_admin"`
}