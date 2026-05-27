package controller

import (
	"apartment-manager-backend/internal/application/dto"
	"apartment-manager-backend/internal/application/service"
	service_error "apartment-manager-backend/internal/application/service/error"
	"apartment-manager-backend/pkg/response"
	"apartment-manager-backend/pkg/validator"
	"errors"
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

	data, err := c.unitService.Create(ctx, apartmentID, &req)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed_to_create_unit", err)
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

	data, err := c.unitService.Update(ctx, unitID, &req)
	if err != nil {
		if errors.Is(err, service_error.ErrUnitNotFound) {
			response.Error(ctx, http.StatusNotFound, "not_found_unit", err)
			return
		}
		response.Error(ctx, http.StatusInternalServerError, "failed_to_update_unit", err)
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

	data, err := c.unitService.PopUser(ctx, unitID)
	if err != nil {
		if errors.Is(err, service_error.ErrUnitNotFound) {
			response.Error(ctx, http.StatusNotFound, "not_found_unit", err)
			return
		}
		response.Error(ctx, http.StatusInternalServerError, "failed_to_pop_user_from_unit", err)
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

	if err := c.unitService.Delete(ctx, unitID); err != nil {
		if errors.Is(err, service_error.ErrUnitNotFound) {
			response.Error(ctx, http.StatusNotFound, "not_found_unit", err)
			return
		}
		response.Error(ctx, http.StatusInternalServerError, "failed_to_delete_unit", err)
		return
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

	unit, err := c.unitService.GetByID(ctx, unitID)
	if err != nil {
		if errors.Is(err, service_error.ErrUnitNotFound) {
			response.Error(ctx, http.StatusNotFound, "not_found_unit", err)
			return
		}
		response.Error(ctx, http.StatusInternalServerError, "failed_to_get_unit", err)
		return
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

	data, err := c.unitService.PushUser(ctx, unitID, &req)
	if err != nil {
		if errors.Is(err, service_error.ErrUnitNotFound) {
			response.Error(ctx, http.StatusNotFound, "not_found_unit", err)
			return
		}
		if errors.Is(err, service_error.ErrUserNotFound) {
			response.Error(ctx, http.StatusNotFound, "not_found_user", err)
			return
		}
		response.Error(ctx, http.StatusInternalServerError, "failed_to_push_user_to_unit", err)
		return
	}

	response.Success(ctx, http.StatusOK, "user_assigned_to_unit_successfully", data)
}
