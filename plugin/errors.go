package plugin

// fieldBuildError represents field error reflection error
type fieldBuildError struct {
	Message string
}

func newBuildError(message string) *fieldBuildError {
	return &fieldBuildError{Message: message}
}

// Error returns error message
func (e *fieldBuildError) Error() string {
	return e.Message
}
