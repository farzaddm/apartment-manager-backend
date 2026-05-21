package error

import "errors"

var ErrTicketNotFound = errors.New("ticket not found")
var ErrTicketIsPrivate = errors.New("ticket is private")
var ErrTicketUnauthorizedAccess = errors.New("ticket does not have authorized access")
