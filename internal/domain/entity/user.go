package entity

import "github.com/google/uuid"

type User struct {
	BaseModel

	ApartmentID *uuid.UUID

	FirstName string
	LastName  string

	Username string
	Email    string
	Phone    string

	Password string

	Role UserRole

	Gender GenderType

	ProfileImageURL string

	Apartment Apartment

	Tickets  []Ticket
	Comments []Comment
}
