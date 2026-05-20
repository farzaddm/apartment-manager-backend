package error

import "errors"

var ErrCommentNotFound = errors.New("comment not found")
var ErrCommentUnauthorizedAccess = errors.New("comment does not have authorized access")
