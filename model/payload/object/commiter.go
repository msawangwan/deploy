package object

// Commiter is a github webhook object
type Commiter struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Username string `json:"username"`
}