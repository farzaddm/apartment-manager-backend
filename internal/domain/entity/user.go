package entity

import (
	"github.com/google/uuid"
)

type User struct {
	BaseModel

	ApartmentID *uuid.UUID `gorm:"column:apartment_id"`

	FirstName string `gorm:"column:first_name"`
	LastName  string `gorm:"column:last_name"`

	Username string `gorm:"column:username"`
	Email    string `gorm:"column:email"`
	Phone    string `gorm:"column:phone"`

	Password string `gorm:"column:password"`

	Role   UserRole    `gorm:"column:role"`
	Gender *GenderType `gorm:"column:gender"`

	ProfileImageURL *string `gorm:"column:profile_image_url"`

	Apartment *Apartment `gorm:"foreignKey:ApartmentID"`

	Tickets  []Ticket  `gorm:"foreignKey:user_id"`
	Comments []Comment `gorm:"foreignKey:user_id"`
	Unit     *Unit
}
