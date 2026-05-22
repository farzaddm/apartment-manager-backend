package controller

import (
	"apartment-manager-backend/internal/application/dto"
	"apartment-manager-backend/internal/application/service"
	"apartment-manager-backend/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RuleController struct {
	ruleService *service.RuleService
}

func NewRuleController(rs *service.RuleService) *RuleController {
	return &RuleController{ruleService: rs}
}

func (c *RuleController) Create(ctx *gin.Context) {
	aptID, err := uuid.Parse(ctx.Param("apartment_id"))
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid_apartment_id", err)
		return
	}

	var req dto.CreateRuleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid_request_body", err)
		return
	}

	res, err := c.ruleService.Create(ctx.Request.Context(), aptID, req)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.Success(ctx, http.StatusCreated, "rule_created_successfully", res)
}

func (c *RuleController) Get(ctx *gin.Context) {
	aptID, err := uuid.Parse(ctx.Param("apartment_id"))
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid_apartment_id", err)
		return
	}
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid_rule_id", err)
		return
	}

	res, err := c.ruleService.GetByID(ctx.Request.Context(), id, aptID)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(ctx, http.StatusOK, "rule_fetch_successfully", res)
}

func (c *RuleController) List(ctx *gin.Context) {
	aptID, err := uuid.Parse(ctx.Param("apartment_id"))
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid_apartment_id", err)
		return
	}

	category := ctx.Query("category")
	var res []dto.RuleResponse

	if category != "" {
		res, err = c.ruleService.ListByApartmentAndCategory(ctx.Request.Context(), aptID, category)
	} else {
		res, err = c.ruleService.ListByApartmentID(ctx.Request.Context(), aptID)
	}

	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}

	response.Success(ctx, http.StatusOK, "rules_fetch_successfully", res)
}

func (c *RuleController) Update(ctx *gin.Context) {
	aptID, err := uuid.Parse(ctx.Param("apartment_id"))
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid_apartment_id", err)
		return
	}
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid_rule_id", err)
		return
	}

	var req dto.UpdateRuleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid_request_body", err)
		return
	}

	res, err := c.ruleService.Update(ctx.Request.Context(), id, aptID, req)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.Success(ctx, http.StatusOK, "rule_update_successfully", res)
}

func (c *RuleController) Delete(ctx *gin.Context) {
	aptID, err := uuid.Parse(ctx.Param("apartment_id"))
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid_apartment_id", err)
		return
	}
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid_rule_id", err)
		return
	}

	if err := c.ruleService.Delete(ctx.Request.Context(), id, aptID); err != nil {
		response.Error(ctx, http.StatusInternalServerError, err.Error(), nil)
		return
	}
	response.Success(ctx, http.StatusOK, "rule_delete_successfully", nil)
}
