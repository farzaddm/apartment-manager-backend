package dto

import (
	"apartment-manager-backend/internal/domain/constant"
	"apartment-manager-backend/internal/domain/entity"
	"context"
	"errors"

	"github.com/google/uuid"
)

type TokenKeys struct {
	userID      uuid.UUID
	role        entity.UserRole
	apartmentID uuid.UUID
}

var (
	ErrUserIDNotFoundInContext      = errors.New("user identifier missing from request context")
	ErrUserRoleNotFoundInContext    = errors.New("user role missing from request context")
	ErrApartmentIDNotFoundInContext = errors.New("apartment identifier missing from request context")

	ErrUserIDCantParseToUUID      = errors.New("failed to parse user ID from request context")
	ErrApartmentIDCantParseToUUID = errors.New("failed to parse apartment ID from request context")
)

func NewTokenKeys(ctx context.Context) (*TokenKeys, error) {
	userID, err := getUserIDKetFromContext(ctx)
	if userID == nil {
		return nil, ErrUserIDNotFoundInContext
	} else if err != nil {
		return nil, ErrUserIDCantParseToUUID
	}

	role, err := getRoleKetFromContext(ctx)
	if role == nil {
		return nil, ErrUserRoleNotFoundInContext
	}

	t := TokenKeys{
		userID: *userID,
		role:   *role,
	}

	apartmentID, ex, err := getApartmentIDKetFromContext(ctx)
	if ex == false {
		return nil, ErrApartmentIDNotFoundInContext
	} else if err == nil {
		if apartmentID == nil {
			t.apartmentID = constant.NilApartmentIDKeyToken
		} else {
			t.apartmentID = *apartmentID
		}
	} else /*err != nil*/ {
		return nil, ErrApartmentIDCantParseToUUID
	}

	return &t, nil
}

func (t *TokenKeys) GetUserID() uuid.UUID {
	return t.userID
}

func (t *TokenKeys) GetApartmentID() uuid.UUID {
	return t.apartmentID
}

func (t *TokenKeys) GetRole() entity.UserRole {
	return t.role
}

func getUserIDKetFromContext(ctx context.Context) (*uuid.UUID, error) {
	rawBaseUserID := ctx.Value(constant.UserIDKeyToken)
	if rawBaseUserID == nil {
		return nil, nil
	}
	str_baseUserID := rawBaseUserID.(string)
	baseUserID, err := uuid.Parse(str_baseUserID)
	if err != nil {
		return nil, err
	}
	return &baseUserID, nil
}

func getRoleKetFromContext(ctx context.Context) (*entity.UserRole, error) {
	rawRole := ctx.Value(constant.RoleKeyToken)
	if rawRole == nil {
		return nil, nil
	}
	str_role := rawRole.(string)
	role := entity.UserRole(str_role)

	return &role, nil
}

func getApartmentIDKetFromContext(ctx context.Context) (*uuid.UUID, bool, error) { // TODO : auth midd need to be edited!!!
	hasApar := ctx.Value(constant.HasApartment)
	if hasApar == nil { // auth midd bug
		return nil, false, nil
	}
	rawBaseApartmentID := ctx.Value(constant.ApartmentIDKeyToken)
	if rawBaseApartmentID == nil { // this user has not any apartment
		return nil, true, nil
	}
	str_baseApartmentID := rawBaseApartmentID.(string)
	baseApartmentID, err := uuid.Parse(str_baseApartmentID)
	if err != nil {
		return nil, true, err
	}
	return &baseApartmentID, true, nil
}
