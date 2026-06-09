package error

import "errors"

var ErrApartmentNotFound = errors.New("apartment not found")
var ErrApartmentUnauthorizedAccess = errors.New("apartment does not have authorized access")
