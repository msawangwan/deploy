package payload

// Owner is a github webhook object
type Owner struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
