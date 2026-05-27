package dto

type CreateUnitRequest struct {
	UnitNumber string `json:"unit_number" binding:"required"`
	Floor      int    `json:"floor" binding:"required"`
}

type UpdateUnitRequest struct {
	UnitNumber string `json:"unit_number"`
	Floor      int    `json:"floor"`
}
