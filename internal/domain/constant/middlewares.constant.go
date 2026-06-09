package constant

import "github.com/google/uuid"

const (
	UserIDKeyToken      string = "user_id"
	RoleKeyToken        string = "role"
	ApartmentIDKeyToken string = "apartment_id"
	HasApartment        string = "has_apartment_id"
)

var (
	NilApartmentIDKeyToken uuid.UUID = uuid.Nil
)
