package webhook

// RequestHeaders represent webhook request headers
type RequestHeaders struct {
	EventName      string `json:"event_name"`
	EventGUID      string `json:"event_guid"`
	EventSignature string `json:"event_sig"`
}
