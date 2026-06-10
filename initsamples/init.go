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
	u4, _ := uuid.Parse("99999999-9999-5555-9999-999999999999")
	u5, _ := uuid.Parse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
	u6, _ := uuid.Parse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb")

	aptID, _ := uuid.Parse("a0000000-0000-0000-0000-000000000000")
	aptID2, _ := uuid.Parse("c0000000-0000-0000-0000-000000000000")

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
			ApartmentID: nil,
			FirstName:   "Voldemort",
			LastName:    "VoldemortianAdmin",
			Username:    "voldemort_ad",
			Email:       "voldemort_ad@gmail.com",
			Phone:       "+14165559797",
			Password:    h,
			Role:        entity.RoleAdmin,
			Gender:      GenderPtr(entity.GenderMale),
		},
		{
			BaseModel:   entity.BaseModel{ID: u4},
			ApartmentID: &aptID2,
			FirstName:   "Voldemort2",
			LastName:    "VoldemortianManager2",
			Username:    "voldemort_man2",
			Email:       "voldemort_man2@gmail.com",
			Phone:       "+14165550992",
			Password:    h,
			Role:        entity.RoleManager,
			Gender:      GenderPtr(entity.GenderMale),
		},
		{
			BaseModel:   entity.BaseModel{ID: u5},
			ApartmentID: &aptID2,
			FirstName:   "Voldemort2",
			LastName:    "VoldemortianResident2",
			Username:    "voldemort_res2",
			Email:       "voldemort_res2@gmail.com",
			Phone:       "+14165550098",
			Password:    h,
			Role:        entity.RoleResident,
			Gender:      GenderPtr(entity.GenderMale),
		},
		{
			BaseModel:   entity.BaseModel{ID: u6},
			ApartmentID: nil,
			FirstName:   "Voldemort2",
			LastName:    "VoldemortianAdmin2",
			Username:    "voldemort_ad2",
			Email:       "voldemort_ad2@gmail.com",
			Phone:       "+14165550099",
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
