package controller

import (
	"apartment-manager-backend/internal/application/dto"
	"apartment-manager-backend/internal/application/service"
	service_error "apartment-manager-backend/internal/application/service/error"
	"apartment-manager-backend/pkg/response"
	"apartment-manager-backend/pkg/validator"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UnitController struct {
	unitService service.UnitService
}

func NewUnitController(unitService service.UnitService) *UnitController {
	return &UnitController{
		unitService: unitService,
	}
}

// Create: POST /apartment/:apartment_id/unit
func (c *UnitController) Create(ctx *gin.Context) {
	aptIDParam := ctx.Param("apartment_id")
	apartmentID, err := uuid.Parse(aptIDParam)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid_apartment_id", err)
		return
	}

	var req dto.CreateUnitRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errList := validator.ParseValidationErrors(err)
		res := &response.StandardResponse{
			Success:    false,
			StatusCode: http.StatusBadRequest,
			Message:    "invalid_request_body",
			Errors:     append(errList, err.Error()),
		}
		res.SendResponse(ctx)
		return
	}

	tokenKeys, err := dto.NewTokenKeys(ctx)
	if err != nil {
		switch err {
		case dto.ErrUserIDNotFoundInContext:
			response.Error(ctx, http.StatusUnauthorized, "user_id_not_found_in_context", err)
			return
		case dto.ErrUserRoleNotFoundInContext:
			response.Error(ctx, http.StatusUnauthorized, "user_role_not_found_in_context", err)
			return
		case dto.ErrApartmentIDNotFoundInContext:
			response.Error(ctx, http.StatusUnauthorized, "apartment_id_not_found_in_context", err)
			return
		case dto.ErrUserIDCantParseToUUID:
			response.Error(ctx, http.StatusBadRequest, "invalid_user_id_format", err)
			return
		case dto.ErrApartmentIDCantParseToUUID:
			response.Error(ctx, http.StatusBadRequest, "invalid_apartment_id_format", err)
			return
		default:
			response.Error(ctx, http.StatusInternalServerError, "unexpected_error", err)
			log.Println(err)
			return
		}
	}

	data, err := c.unitService.Create(ctx, tokenKeys, apartmentID, &req)
	if err != nil {
		switch err {
		case service_error.ErrUnitUnauthorizedAccess:
			response.Error(ctx, http.StatusUnauthorized, "unauthorized_access_to_unit", err)
		default:
			response.Error(ctx, http.StatusInternalServerError, "failed_to_create_unit", err)
		}
		return
	}

	response.Success(ctx, http.StatusCreated, "unit_created_successfully", data)
}

// Update: PUT /apartment/:apartment_id/unit/:unit_id
func (c *UnitController) Update(ctx *gin.Context) {
	unitIDParam := ctx.Param("unit_id")
	unitID, err := uuid.Parse(unitIDParam)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid_unit_id", err)
		return
	}

	var req dto.UpdateUnitRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errList := validator.ParseValidationErrors(err)
		res := &response.StandardResponse{
			Success:    false,
			StatusCode: http.StatusBadRequest,
			Message:    "invalid_request_body",
			Errors:     append(errList, err.Error()),
		}
		res.SendResponse(ctx)
		return
	}

	tokenKeys, err := dto.NewTokenKeys(ctx)
	if err != nil {
		switch err {
		case dto.ErrUserIDNotFoundInContext:
			response.Error(ctx, http.StatusUnauthorized, "user_id_not_found_in_context", err)
			return
		case dto.ErrUserRoleNotFoundInContext:
			response.Error(ctx, http.StatusUnauthorized, "user_role_not_found_in_context", err)
			return
		case dto.ErrApartmentIDNotFoundInContext:
			response.Error(ctx, http.StatusUnauthorized, "apartment_id_not_found_in_context", err)
			return
		case dto.ErrUserIDCantParseToUUID:
			response.Error(ctx, http.StatusBadRequest, "invalid_user_id_format", err)
			return
		case dto.ErrApartmentIDCantParseToUUID:
			response.Error(ctx, http.StatusBadRequest, "invalid_apartment_id_format", err)
			return
		default:
			response.Error(ctx, http.StatusInternalServerError, "unexpected_error", err)
			log.Println(err)
			return
		}
	}

	data, err := c.unitService.Update(ctx, tokenKeys, unitID, &req)
	if err != nil {
		switch err {
		case service_error.ErrUnitUnauthorizedAccess:
			response.Error(ctx, http.StatusUnauthorized, "unauthorized_access_to_unit", err)

		case service_error.ErrUnitNotFound:
			response.Error(ctx, http.StatusNotFound, "not_found_unit", err)

		default:
			response.Error(ctx, http.StatusInternalServerError, "failed_to_update_unit", err)
		}
		return
	}

	response.Success(ctx, http.StatusOK, "unit_updated_successfully", data)
}

// PopUser: PATCH /apartment/:apartment_id/unit/:unit_id
func (c *UnitController) PopUser(ctx *gin.Context) {
	unitIDParam := ctx.Param("unit_id")
	unitID, err := uuid.Parse(unitIDParam)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid_unit_id", err)
		return
	}

	tokenKeys, err := dto.NewTokenKeys(ctx)
	if err != nil {
		switch err {
		case dto.ErrUserIDNotFoundInContext:
			response.Error(ctx, http.StatusUnauthorized, "user_id_not_found_in_context", err)
			return
		case dto.ErrUserRoleNotFoundInContext:
			response.Error(ctx, http.StatusUnauthorized, "user_role_not_found_in_context", err)
			return
		case dto.ErrApartmentIDNotFoundInContext:
			response.Error(ctx, http.StatusUnauthorized, "apartment_id_not_found_in_context", err)
			return
		case dto.ErrUserIDCantParseToUUID:
			response.Error(ctx, http.StatusBadRequest, "invalid_user_id_format", err)
			return
		case dto.ErrApartmentIDCantParseToUUID:
			response.Error(ctx, http.StatusBadRequest, "invalid_apartment_id_format", err)
			return
		default:
			response.Error(ctx, http.StatusInternalServerError, "unexpected_error", err)
			log.Println(err)
			return
		}
	}

	data, err := c.unitService.PopUser(ctx, tokenKeys, unitID)
	if err != nil {
		switch err {
		case service_error.ErrUnitUnauthorizedAccess:
			response.Error(ctx, http.StatusUnauthorized, "unauthorized_access_to_unit", err)

		case service_error.ErrUnitNotFound:
			response.Error(ctx, http.StatusNotFound, "not_found_unit", err)

		default:
			response.Error(ctx, http.StatusInternalServerError, "failed_to_pop_user_from_unit", err)
		}
		return
	}

	response.Success(ctx, http.StatusOK, "user_popped_from_unit_successfully", data)
}

// Delete: DELETE /apartment/:apartment_id/unit/:unit_id
func (c *UnitController) Delete(ctx *gin.Context) {
	unitIDParam := ctx.Param("unit_id")
	unitID, err := uuid.Parse(unitIDParam)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid_unit_id", err)
		return
	}

	tokenKeys, err := dto.NewTokenKeys(ctx)
	if err != nil {
		switch err {
		case dto.ErrUserIDNotFoundInContext:
			response.Error(ctx, http.StatusUnauthorized, "user_id_not_found_in_context", err)
			return
		case dto.ErrUserRoleNotFoundInContext:
			response.Error(ctx, http.StatusUnauthorized, "user_role_not_found_in_context", err)
			return
		case dto.ErrApartmentIDNotFoundInContext:
			response.Error(ctx, http.StatusUnauthorized, "apartment_id_not_found_in_context", err)
			return
		case dto.ErrUserIDCantParseToUUID:
			response.Error(ctx, http.StatusBadRequest, "invalid_user_id_format", err)
			return
		case dto.ErrApartmentIDCantParseToUUID:
			response.Error(ctx, http.StatusBadRequest, "invalid_apartment_id_format", err)
			return
		default:
			response.Error(ctx, http.StatusInternalServerError, "unexpected_error", err)
			log.Println(err)
			return
		}
	}

	if err := c.unitService.Delete(ctx, tokenKeys, unitID); err != nil {
		switch err {

		case service_error.ErrUnitNotFound:
			response.Error(ctx, http.StatusNotFound, "unit_not_found", err)
			return

		case service_error.ErrUnitUnauthorizedAccess:
			response.Error(ctx, http.StatusForbidden, "unit_unauthorized_access", err)
			return

		default:
			response.Error(ctx, http.StatusInternalServerError, "delete_unit_failed", err)
			return
		}
	}

	response.Success(ctx, http.StatusOK, "unit_deleted_successfully", nil)
}

// GetByID: GET /apartment/:apartment_id/unit/:unit_id
func (c *UnitController) GetByID(ctx *gin.Context) {
	unitIDParam := ctx.Param("unit_id")
	unitID, err := uuid.Parse(unitIDParam)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid_unit_id", err)
		return
	}

	tokenKeys, err := dto.NewTokenKeys(ctx)
	if err != nil {
		switch err {
		case dto.ErrUserIDNotFoundInContext:
			response.Error(ctx, http.StatusUnauthorized, "user_id_not_found_in_context", err)
			return
		case dto.ErrUserRoleNotFoundInContext:
			response.Error(ctx, http.StatusUnauthorized, "user_role_not_found_in_context", err)
			return
		case dto.ErrApartmentIDNotFoundInContext:
			response.Error(ctx, http.StatusUnauthorized, "apartment_id_not_found_in_context", err)
			return
		case dto.ErrUserIDCantParseToUUID:
			response.Error(ctx, http.StatusBadRequest, "invalid_user_id_format", err)
			return
		case dto.ErrApartmentIDCantParseToUUID:
			response.Error(ctx, http.StatusBadRequest, "invalid_apartment_id_format", err)
			return
		default:
			response.Error(ctx, http.StatusInternalServerError, "unexpected_error", err)
			log.Println(err)
			return
		}
	}

	unit, err := c.unitService.GetByID(ctx, tokenKeys, unitID)
	if err != nil {
		switch err {

		case service_error.ErrUnitNotFound:
			response.Error(ctx, http.StatusNotFound, "unit_not_found", err)
			return

		case service_error.ErrUnitUnauthorizedAccess:
			response.Error(ctx, http.StatusForbidden, "forbidden_unit_access", err)
			return

		default:
			response.Error(ctx, http.StatusInternalServerError, "get_unit_failed", err)
			return
		}
	}

	response.Success(ctx, http.StatusOK, "unit_fetched_successfully", unit)
}

// PushUser: POST /apartments/:apartment_id/units/:unit_id/users
func (c *UnitController) PushUser(ctx *gin.Context) {
	unitIDParam := ctx.Param("unit_id")
	unitID, err := uuid.Parse(unitIDParam)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid_unit_id", err)
		return
	}

	var req dto.PushUserToUnitRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		errList := validator.ParseValidationErrors(err)
		res := &response.StandardResponse{
			Success:    false,
			StatusCode: http.StatusBadRequest,
			Message:    "invalid_request_body",
			Errors:     errList,
		}
		res.SendResponse(ctx)
		return
	}

	tokenKeys, err := dto.NewTokenKeys(ctx)
	if err != nil {
		switch err {
		case dto.ErrUserIDNotFoundInContext:
			response.Error(ctx, http.StatusUnauthorized, "user_id_not_found_in_context", err)
			return
		case dto.ErrUserRoleNotFoundInContext:
			response.Error(ctx, http.StatusUnauthorized, "user_role_not_found_in_context", err)
			return
		case dto.ErrApartmentIDNotFoundInContext:
			response.Error(ctx, http.StatusUnauthorized, "apartment_id_not_found_in_context", err)
			return
		case dto.ErrUserIDCantParseToUUID:
			response.Error(ctx, http.StatusBadRequest, "invalid_user_id_format", err)
			return
		case dto.ErrApartmentIDCantParseToUUID:
			response.Error(ctx, http.StatusBadRequest, "invalid_apartment_id_format", err)
			return
		default:
			response.Error(ctx, http.StatusInternalServerError, "unexpected_error", err)
			log.Println(err)
			return
		}
	}

	data, err := c.unitService.PushUser(ctx, tokenKeys, unitID, &req)
	if err != nil {
		switch err {

		case service_error.ErrUnitNotFound:
			response.Error(ctx, http.StatusNotFound, "unit_not_found", err)
			return

		case service_error.ErrUserNotFound:
			response.Error(ctx, http.StatusNotFound, "user_not_found", err)
			return

		case service_error.ErrUnitUnauthorizedAccess:
			response.Error(ctx, http.StatusForbidden, "forbidden_unit_access", err)
			return

		default:
			response.Error(ctx, http.StatusInternalServerError, "push_user_to_unit_failed", err)
			return
		}
	}

	response.Success(ctx, http.StatusOK, "user_assigned_to_unit_successfully", data)
}
