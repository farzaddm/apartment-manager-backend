package error

import "errors"

var ErrUserNotFound = errors.New("user not found")
var ErrUserUnauthorizedAccess = errors.New("user does not have authorized access")
