package initsamples

import (
	"apartment-manager-backend/internal/domain/entity"
	"apartment-manager-backend/pkg/hasher"
	"log"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func StringPtr(s string) *string                       { return &s }
func RolePtr(r entity.UserRole) *entity.UserRole       { return &r }
func GenderPtr(g entity.GenderType) *entity.GenderType { return &g }

// TODO : Move this Func to right directories and REFACTOR IT
func CreateOrOverWriteManagersAndAdminAndResident(db *gorm.DB, passwordHasher *hasher.BcryptHasher) {
	password := "123456"

	h, _ := passwordHasher.Hash(password)

	u1, _ := uuid.Parse("77777777-7777-7777-7777-777777777777")
	u2, _ := uuid.Parse("88888887-8888-8888-8888-888888888888")
	u3, _ := uuid.Parse("99999999-9999-9999-9999-999999999999")
	aptID, _ := uuid.Parse("a0000000-0000-0000-0000-000000000000")

	users := []entity.User{
		{
			BaseModel:   entity.BaseModel{ID: u1},
			ApartmentID: &aptID,
			FirstName:   "Voldemort",
			LastName:    "VoldemortianResident",
			Username:    "voldemort_res",
			Email:       "voldemort_res@gmail.com",
			Phone:       "+14165550097",
			Password:    h,
			Role:        entity.RoleResident,
			Gender:      GenderPtr(entity.GenderMale),
		},
		{
			BaseModel:   entity.BaseModel{ID: u2},
			ApartmentID: &aptID,
			FirstName:   "Voldemort",
			LastName:    "VoldemortianManager",
			Username:    "voldemort_man",
			Email:       "voldemort_man@gmail.com",
			Phone:       "+14165550999",
			Password:    h,
			Role:        entity.RoleManager,
			Gender:      GenderPtr(entity.GenderMale),
		},
		{
			BaseModel:   entity.BaseModel{ID: u3},
			ApartmentID: &aptID,
			FirstName:   "Voldemort",
			LastName:    "VoldemortianAdmin",
			Username:    "voldemort_ad",
			Email:       "voldemort_ad@gmail.com",
			Phone:       "+14165559797",
			Password:    h,
			Role:        entity.RoleAdmin,
			Gender:      GenderPtr(entity.GenderMale),
		},
	}

	for _, u := range users {
		err := db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true,
		}).Create(&u).Error

		if err != nil {
			log.Printf("Error seeding user %s: %v", u.Username, err)
		}
	}
}
