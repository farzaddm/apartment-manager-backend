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

func CreateOrOverWriteManagersAndAdminAndResident(db *gorm.DB, passwordHasher *hasher.BcryptHasher) {
	password := "123456"
	h, _ := passwordHasher.Hash(password)

	// ============================================================================
	// 1. UUID DEFINITIONS
	// ============================================================================

	// User UUIDs
	u1, _ := uuid.Parse("77777777-7777-7777-7777-777777777777") // Resident 1 (Apartment A)
	u2, _ := uuid.Parse("88888887-8888-8888-8888-888888888888")
	u3, _ := uuid.Parse("99999999-9999-9999-9999-999999999999")
	u4, _ := uuid.Parse("99999999-9999-5555-9999-999999999999")
	u5, _ := uuid.Parse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa") // Resident 2 (Apartment C)
	u6, _ := uuid.Parse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb")

	// Apartment UUIDs
	aptID, _ := uuid.Parse("a0000000-0000-0000-0000-000000000000")  // Apartment A
	aptID2, _ := uuid.Parse("c0000000-0000-0000-0000-000000000000") // Apartment C

	// Unit UUIDs for Resident Voldemorts
	unitID1, _ := uuid.Parse("77777777-eeee-eeee-eeee-777777777777") // Unit for Resident 1
	unitID2, _ := uuid.Parse("aaaaaaaa-eeee-eeee-eeee-aaaaaaaaaaaa") // Unit for Resident 2

	// ============================================================================
	// 2. USER SEED DATA
	// ============================================================================
	users := []entity.User{
		{
			BaseModel:   entity.BaseModel{ID: u1},
			ApartmentID: &aptID,
			FirstName:   "ولدمورت ساکن",
			LastName:    "پلید زاده",
			Username:    "voldemort_res",
			Email:       "voldemort_res@gmail.com",
			Phone:       "+989193456792",
			Password:    h,
			Role:        entity.RoleResident,
			Gender:      GenderPtr(entity.GenderMale),
		},
		{
			BaseModel:   entity.BaseModel{ID: u2},
			ApartmentID: &aptID,
			FirstName:   "ولدمورت مدیر",
			LastName:    "پلید زاده",
			Username:    "voldemort_man",
			Email:       "voldemort_man@gmail.com",
			Phone:       "+989123406792",
			Password:    h,
			Role:        entity.RoleManager,
			Gender:      GenderPtr(entity.GenderMale),
		},
		{
			BaseModel:   entity.BaseModel{ID: u3},
			ApartmentID: nil,
			FirstName:   "ولدمورت ادمین",
			LastName:    "پلید زاده",
			Username:    "voldemort_ad",
			Email:       "voldemort_ad@gmail.com",
			Phone:       "+989123456797",
			Password:    h,
			Role:        entity.RoleAdmin,
			Gender:      GenderPtr(entity.GenderMale),
		},
		{
			BaseModel:   entity.BaseModel{ID: u4},
			ApartmentID: &aptID2,
			FirstName:   "ولدمورت پلاس مدیر",
			LastName:    "پلید پلید زاده",
			Username:    "voldemort_man2",
			Email:       "voldemort_man2@gmail.com",
			Phone:       "+989123456997",
			Password:    h,
			Role:        entity.RoleManager,
			Gender:      GenderPtr(entity.GenderMale),
		},
		{
			BaseModel:   entity.BaseModel{ID: u5},
			ApartmentID: &aptID2,
			FirstName:   "ولدمورت پلاس ساکن",
			LastName:    "پلید پلید زاده",
			Username:    "voldemort_res2",
			Email:       "voldemort_res2@gmail.com",
			Phone:       "+989123456007",
			Password:    h,
			Role:        entity.RoleResident,
			Gender:      GenderPtr(entity.GenderMale),
		},
		{
			BaseModel:   entity.BaseModel{ID: u6},
			ApartmentID: nil,
			FirstName:   "ولدمورت پلاس ادمین",
			LastName:    "پلید پلید زاده",
			Username:    "voldemort_ad2",
			Email:       "voldemort_ad2@gmail.com",
			Phone:       "+989123450123",
			Password:    h,
			Role:        entity.RoleAdmin,
			Gender:      GenderPtr(entity.GenderMale),
		},
	}

	// ============================================================================
	// 3. UNIT SEED DATA
	// ============================================================================
	units := []entity.Unit{
		{
			BaseModel:   entity.BaseModel{ID: unitID1},
			ApartmentID: aptID,
			UserID:      &u1, // Map to Voldemort Resident 1
			UnitNumber:  "303",
			Floor:       3,
		},
		{
			BaseModel:   entity.BaseModel{ID: unitID2},
			ApartmentID: aptID2,
			UserID:      &u5, // Map to Voldemort Resident 2
			UnitNumber:  "303",
			Floor:       3,
		},
	}

	// ============================================================================
	// 4. DATABASE SEED EXECUTION (with Upsert logic)
	// ============================================================================

	// Seed Users First
	for _, u := range users {
		err := db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true,
		}).Create(&u).Error

		if err != nil {
			log.Printf("Error seeding user %s: %v", u.Username, err)
		}
	}

	// Seed Units Second (Prevents foreign key violations)
	for _, unit := range units {
		err := db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true,
		}).Create(&unit).Error

		if err != nil {
			log.Printf("Error seeding unit %s for user %s: %v", unit.UnitNumber, *unit.UserID, err)
		}
	}
}
