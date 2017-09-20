package err

import "fmt"

type NotEnoughInfoToProceedError struct {
	InternalError  error
	AdditionalInfo string
}

func (e NotEnoughInfoToProceedError) Error() string {
	return fmt.Sprintf("cannot continue with operation, missing information [%s]: %g", e.AdditionalInfo, e.InternalError)
}

type BuildfileNotFoundError struct {
	InternalError  error
	AdditionalInfo string
}

func (e BuildfileNotFoundError) Error() string {
	return fmt.Sprintf("failed to find a buildfile [%s] %g", e.AdditionalInfo, e.InternalError)
}
