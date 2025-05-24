// auth/errors.go
package auth
import "errors"
var (
	ErrUserExists       = errors.New("user already exists")
	ErrInvalidEmail     = errors.New("invalid email format")
	ErrWeakPassword     = errors.New("password too weak")
)