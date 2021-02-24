package plugin

// buildError represents field error reflection error
type buildError struct {
	Message string
}

// newBuildError constructs new build error
func newBuildError(message string) *buildError {
	return &buildError{Message: message}
}

// Error returns error message
func (e *buildError) Error() string {
	return e.Message
}
