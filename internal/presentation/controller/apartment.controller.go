package controller

import (
	"apartment-manager-backend/internal/application/dto"
	"apartment-manager-backend/internal/application/service"
	"apartment-manager-backend/internal/domain/entity"
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
		response.Error(ctx, http.StatusBadRequest, "invalid request body", err)
		return
	}

	apartment := &entity.Apartment{
		Name:       req.Name,
		Province:   req.Province,
		City:       req.City,
		Address:    req.Address,
		PostalCode: req.PostalCode,
	}

	if err := c.apartmentService.Create(ctx, apartment); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to create apartment", err)
		return
	}

	response.Success(ctx, http.StatusCreated, "apartment created successfully", apartment)
}

func (c *ApartmentController) Update(ctx *gin.Context) {
	idParam := ctx.Param("apartment_id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid apartment id", err)
		return
	}

	var req dto.UpdateApartmentRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid request body", err)
		return
	}

	apartment := &entity.Apartment{
		BaseModel: entity.BaseModel{
			ID: id,
		},
		Name:       req.Name,
		Province:   req.Province,
		City:       req.City,
		Address:    req.Address,
		PostalCode: req.PostalCode,
	}

	if err := c.apartmentService.Update(ctx, apartment); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to update apartment", err)
		return
	}

	response.Success(ctx, http.StatusOK, "apartment updated successfully", apartment)
}

func (c *ApartmentController) Delete(ctx *gin.Context) {
	idParam := ctx.Param("apartment_id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid apartment id", err)
		return
	}

	if err := c.apartmentService.Delete(ctx, id); err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to delete apartment", err)
		return
	}

	response.Success(ctx, http.StatusOK, "apartment deleted successfully", nil)
}

func (c *ApartmentController) GetByID(ctx *gin.Context) {
	idParam := ctx.Param("apartment_id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid apartment id", err)
		return
	}

	apartment, err := c.apartmentService.GetByID(ctx, id)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to get apartment", err)
		return
	}

	response.Success(ctx, http.StatusOK, "apartment fetched successfully", apartment)
}

func (c *ApartmentController) GetByIDWithUsers(ctx *gin.Context) {
	idParam := ctx.Param("apartment_id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid apartment id", err)
		return
	}

	apartment, err := c.apartmentService.GetWithUsers(ctx, id)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to get apartment users", err)
		return
	}

	response.Success(ctx, http.StatusOK, "apartment users fetched successfully", apartment)
}

func (c *ApartmentController) GetByIDWithRules(ctx *gin.Context) {
	idParam := ctx.Param("apartment_id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid apartment id", err)
		return
	}

	apartment, err := c.apartmentService.GetWithRules(ctx, id)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to get apartment rules", err)
		return
	}

	response.Success(ctx, http.StatusOK, "apartment rules fetched successfully", apartment)
}

func (c *ApartmentController) GetByIDWithAnnouncements(ctx *gin.Context) {
	idParam := ctx.Param("apartment_id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid apartment id", err)
		return
	}

	apartment, err := c.apartmentService.GetWithAnnouncements(ctx, id)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to get apartment announcements", err)
		return
	}

	response.Success(ctx, http.StatusOK, "apartment announcements fetched successfully", apartment)
}

func (c *ApartmentController) GetByIDWithInviteCodes(ctx *gin.Context) {
	idParam := ctx.Param("apartment_id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid apartment id", err)
		return
	}

	apartment, err := c.apartmentService.GetWithInviteCodes(ctx, id)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed to get apartment invite codes", err)
		return
	}

	response.Success(ctx, http.StatusOK, "apartment invite codes fetched successfully", apartment)
}
