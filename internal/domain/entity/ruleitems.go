package entity

import "github.com/google/uuid"

type RuleItem struct {
	BaseModel

	RuleID uuid.UUID

	Body string

	Rule Rule
}
