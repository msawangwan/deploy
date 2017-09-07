package webhook

// RequestHeaders represent webhook request headers
type RequestHeaders struct {
	EventName      string
	EventGUID      string
	EventSignature string
}