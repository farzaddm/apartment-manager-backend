package controller

import (
	"apartment-manager-backend/internal/application/service"

	"github.com/gin-gonic/gin"
)

type TicketController struct {
	ticketService service.TicketService
}

func NewTicketController(ticketService service.TicketService) *TicketController {
	return &TicketController{
		ticketService: ticketService,
	}
}

func (c *TicketController) Create(ctx *gin.Context) {}

func (c *TicketController) GetByID(ctx *gin.Context) {}

func (c *TicketController) GetFully(ctx *gin.Context) {}

func (c *TicketController) List(ctx *gin.Context) {}

func (c *TicketController) Update(ctx *gin.Context) {}

func (c *TicketController) UpdateTicketStatus(ctx *gin.Context) {}

func (c *TicketController) Delete(ctx *gin.Context) {}
