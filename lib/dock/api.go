package dock

import "fmt"

// APIResponseError ...
type APIResponseError struct {
	Message      string `json:"message,omitempty"`
	ExpectedCode int    `json:"-"`
	ActualCode   int    `json:"-"`
}

// Error ...
func (are APIResponseError) Error() string {
	return fmt.Sprintf(
		"[api_response_error][status_code_mismatch][expected: %d][actual: %d][message: %s]",
		are.ExpectedCode,
		are.ActualCode,
		are.Message,
	)
}
