package codegen

var _ = Error("")

// Error is a simple type for classifying errors
type Error string

const (
	// EAlreadyExist indicates that file we want to generate exist already
	// and we will not append to it
	EAlreadyExist Error = "file already exist"
)

// Error satisfies the error interface
func (e Error) Error() string {
	return string(e)
}
