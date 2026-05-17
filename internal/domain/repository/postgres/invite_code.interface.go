package postgres

import (
	"apartment-manager-backend/internal/domain/entity"
	"context"
)

type InviteCodeInterface interface {
	Create(ctx context.Context, inviteCode *entity.InviteCode) error

	GetActiveInviteCode(ctx context.Context, apartmentID string, unitID string) (*entity.InviteCode, error)
	//ValidateInviteCode(ctx context.Context) error

}
