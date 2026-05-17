package service

import (
	"apartment-manager-backend/internal/application/dto"
	"apartment-manager-backend/internal/domain/entity"
	"apartment-manager-backend/internal/domain/repository/postgres"
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type InviteCodeService struct {
	inviteCodeRepo postgres.InviteCodeInterface
	unitRepo       postgres.UnitInterface
	apartmentRepo  postgres.ApartmentInterface
}

func NewInviteCodeService(inviteCodeRepo postgres.InviteCodeInterface, unitRepo postgres.UnitInterface, apartmentRepo postgres.ApartmentInterface) *InviteCodeService {
	return &InviteCodeService{
		inviteCodeRepo: inviteCodeRepo,
		unitRepo:       unitRepo,
		apartmentRepo:  apartmentRepo,
	}
}

func (s *InviteCodeService) Create(ctx context.Context, req dto.CreateInviteRequest, actorUserID string) (*entity.InviteCode, error) {
	if req.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("expiration date must be a future timestamp")
	}

	existingInvite, err := s.inviteCodeRepo.GetActiveInviteCode(ctx, req.ApartmentID, req.UnitID)
	if err != nil {
		return nil, fmt.Errorf("failed checking existing invites: %w", err)
	}
	if existingInvite != nil {
		return nil, fmt.Errorf("an active invite code already exists for this unit (Code: %s) expiring at %s",
			existingInvite.Code, existingInvite.ExpiresAt.Format("2006-01-02 15:04:05"))
	}

	// 1. Fetch unit details to execute cross-verifications
	unit, err := s.unitRepo.GetByID(ctx, req.UnitID)
	if err != nil {
		return nil, errors.New("unit not found")
	}

	// Rule Check A: Verify that the unit actually belongs to the specified building
	if unit.ApartmentID.String() != req.ApartmentID {
		return nil, errors.New("the specified unit does not belong to this apartment building")
	}

	// Rule Check B: Verify that the requested unit is currently vacant (UserID is nil)
	if unit.UserID != nil && *unit.UserID != uuid.Nil {
		return nil, errors.New("cannot generate invite code: this unit is already occupied")
	}

	// Rule Check C: Verify that the requesting user is the building's designated manager
	realManagerID, err := s.apartmentRepo.GetApartmentManagerID(ctx, req.ApartmentID)
	if err != nil {
		return nil, err
	}
	if realManagerID != actorUserID {
		return nil, errors.New("permission denied: you are not the manager of this apartment building")
	}

	bytes := make([]byte, 6)
	if _, err := rand.Read(bytes); err != nil {
		return nil, err
	}
	generatedCode := hex.EncodeToString(bytes)

	aptUUID, _ := uuid.Parse(req.ApartmentID)
	unitUUID, _ := uuid.Parse(req.UnitID)

	inviteCode := &entity.InviteCode{
		ApartmentID: aptUUID,
		UnitID:      unitUUID,
		Code:        generatedCode,
		ExpiresAt:   *req.ExpiresAt,
	}

	// 3. Save to database
	err = s.inviteCodeRepo.Create(ctx, inviteCode)
	if err != nil {
		return nil, err
	}

	return inviteCode, nil
}
