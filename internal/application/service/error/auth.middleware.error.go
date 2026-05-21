package error

import "errors"

var (
	ErrUserIDNotFoundInContext   = errors.New("user identifier missing from request context")
	ErrUserRoleNotFoundInContext = errors.New("user role missing from request context")
)
