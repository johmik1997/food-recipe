package utils
import "fmt"

// NewValidationError creates a new error with the given validation message.
func NewValidationError(msg string) error {
    return fmt.Errorf("validation error: %s", msg)
}

type ValidationError struct {
    Message string
}

func (e *ValidationError) Error() string {
    return e.Message
}