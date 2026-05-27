package error

import "errors"

var (
	ErrUnitNotFound      = errors.New("unit_not_found")
	ErrUnitUpdateFailed  = errors.New("unit_update_failed")
	ErrUnitAlreadyExists = errors.New("unit_number_already_exists_in_floor")
)
