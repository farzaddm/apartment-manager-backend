package dto

import (
	"apartment-manager-backend/internal/domain/entity"
)

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

func MapRuleToResponse(rule *entity.Rule) *RuleResponse {
	if rule == nil {
		return nil
	}
	return &RuleResponse{
		ID:          rule.ID.String(),
		ApartmentID: rule.ApartmentID.String(),
		Title:       rule.Title,
		Description: rule.Description,
		Category:    string(rule.Category),
		CreatedAt:   rule.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}

func MapRulesToSliceResponse(rules []entity.Rule) []RuleResponse {
	if len(rules) == 0 {
		return []RuleResponse{}
	}
	res := make([]RuleResponse, len(rules))
	for i, r := range rules {
		mapped := MapRuleToResponse(&r)
		if mapped != nil {
			res[i] = *mapped
		}
	}
	return res
}
