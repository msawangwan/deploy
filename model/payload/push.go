package payload

import "github.com/msawangwan/ci.io/model/payload/structs"

// PushEvent describes a webhook
type PushEvent struct {
	Repository structs.Repository
}
