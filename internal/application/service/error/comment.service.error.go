package error

import "errors"

var ErrCommentNotFound = errors.New("comment not found")
var ErrCommentUnauthorizedAccess = errors.New("comment does not have authorized access")
var ErrTicketOfCommentOrCommentNotFound = errors.New("ticket of comment or comment not found")
var ErrUserOfCommentOrCommentNotFound = errors.New("user of comment or user not found")
var ErrTicketOfCommentNotFound = errors.New("ticket of comment not found")
