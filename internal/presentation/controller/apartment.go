package controller

import (
	"apartment-manager-backend/internal/application/service"

	"github.com/gin-gonic/gin"
)

type ApartmentController struct {
	apartmentService service.ApartmentService
}

func NewApartmentController(apartmentService service.ApartmentService) *ApartmentController {
	return &ApartmentController{
		apartmentService: apartmentService,
	}
}

func (c *ApartmentController) Create(ctx *gin.Context) {}

func (c *ApartmentController) GetByID(ctx *gin.Context) {}

func (c *ApartmentController) GetByIDWithUsers(ctx *gin.Context) {}

func (c *ApartmentController) GetByIDWithRules(ctx *gin.Context) {}

func (c *ApartmentController) GetByIDWithAnnouncements(ctx *gin.Context) {}

func (c *ApartmentController) GetByIDWithInviteCodes(ctx *gin.Context) {}

func (c *ApartmentController) Update(ctx *gin.Context) {}

func (c *ApartmentController) Delete(ctx *gin.Context) {}
