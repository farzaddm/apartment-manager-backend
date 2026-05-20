package controller

import (
	"apartment-manager-backend/internal/application/dto"
	"apartment-manager-backend/internal/application/service"
	"apartment-manager-backend/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ApartmentController struct {
	apartmentService service.ApartmentService
}

func NewApartmentController(apartmentService service.ApartmentService) *ApartmentController {
	return &ApartmentController{
		apartmentService: apartmentService,
	}
}

func (c *ApartmentController) Create(ctx *gin.Context) {
	var req dto.CreateApartmentRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid_request_body", err)
		return
	}

	apartment := &dto.CreateApartmentRequest{
		Name:       req.Name,
		Province:   req.Province,
		City:       req.City,
		Address:    req.Address,
		PostalCode: req.PostalCode,
	}

	if err := c.apartmentService.Create(ctx, apartment); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed_to_create_apartment", err)
		return
	}

	response.Success(ctx, http.StatusCreated, "apartment_created_successfully", apartment)
}

func (c *ApartmentController) Update(ctx *gin.Context) {
	idParam := ctx.Param("apartment_id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid_apartment_id", err)
		return
	}

	var req dto.UpdateApartmentRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid_request_body", err)
		return
	}

	apartment := &dto.UpdateApartmentRequest{
		Name:       req.Name,
		Province:   req.Province,
		City:       req.City,
		Address:    req.Address,
		PostalCode: req.PostalCode,
	}

	if err := c.apartmentService.Update(ctx,id, apartment); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed_to_update_apartment", err)
		return
	}

	response.Success(ctx, http.StatusOK, "apartment_updated_successfully", apartment)
}

func (c *ApartmentController) Delete(ctx *gin.Context) {
	idParam := ctx.Param("apartment_id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid_apartment_id", err)
		return
	}

	if err := c.apartmentService.Delete(ctx, id); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed_to_delete_apartment", err)
		return
	}

	response.Success(ctx, http.StatusOK, "apartment_deleted_successfully", nil)
}

func (c *ApartmentController) GetByID(ctx *gin.Context) {
	idParam := ctx.Param("apartment_id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid_apartment_id", err)
		return
	}

	apartment, err := c.apartmentService.GetByID(ctx, id)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed_to_get_apartment", err)
		return
	}

	response.Success(ctx, http.StatusOK, "apartment_fetched_successfully", apartment)
}

func (c *ApartmentController) GetByIDWithUsers(ctx *gin.Context) {
	idParam := ctx.Param("apartment_id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid_apartment_id", err)
		return
	}

	apartment, err := c.apartmentService.GetWithUsers(ctx, id)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed_to_get_apartment_users", err)
		return
	}

	response.Success(ctx, http.StatusOK, "apartment_users_fetched_successfully", apartment)
}

func (c *ApartmentController) GetByIDWithRules(ctx *gin.Context) {
	idParam := ctx.Param("apartment_id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid_apartment_id", err)
		return
	}

	apartment, err := c.apartmentService.GetWithRules(ctx, id)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed_to_get_apartment_rules", err)
		return
	}

	response.Success(ctx, http.StatusOK, "apartment_rules_fetched_successfully", apartment)
}

func (c *ApartmentController) GetByIDWithAnnouncements(ctx *gin.Context) {
	idParam := ctx.Param("apartment_id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid_apartment_id", err)
		return
	}

	apartment, err := c.apartmentService.GetWithAnnouncements(ctx, id)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed_to_get_apartment_announcements", err)
		return
	}

	response.Success(ctx, http.StatusOK, "apartment_announcements_fetched_successfully", apartment)
}

func (c *ApartmentController) GetByIDWithInviteCodes(ctx *gin.Context) {
	idParam := ctx.Param("apartment_id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid_apartment_id", err)
		return
	}

	apartment, err := c.apartmentService.GetWithInviteCodes(ctx, id)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed_to_get_apartment_invite_codes", err)
		return
	}

	response.Success(ctx, http.StatusOK, "apartment_invite_codes_fetched_successfully", apartment)
}
