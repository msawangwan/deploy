package payload

// Pusher is a github webhook object
type Pusher struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
