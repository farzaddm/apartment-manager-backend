package dto

type CreateRuleRequest struct {
	Title       string `json:"title" binding:"required,min=3,max=255"`
	Description string `json:"description" binding:"required"`
	Category    string `json:"category" binding:"required"`
}

type UpdateRuleRequest struct {
	Title       string `json:"title" binding:"omitempty,min=3,max=255"`
	Description string `json:"description" binding:"omitempty"`
	Category    string `json:"category" binding:"omitempty"`
}

type RuleResponse struct {
	ID          string `json:"id"`
	ApartmentID string `json:"apartment_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Category    string `json:"category"`
	CreatedAt   string `json:"created_at"`
}
