package controller

import (
	"apartment-manager-backend/internal/application/dto"
	"apartment-manager-backend/internal/application/service"
	service_error "apartment-manager-backend/internal/application/service/error"
	"apartment-manager-backend/pkg/response"
	"apartment-manager-backend/pkg/validator"
	"errors"
	"log"
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
	data, err := c.apartmentService.Create(ctx, tokenKeys, &req)
	if err != nil {
		response.Error(ctx, http.StatusInternalServerError, "failed_to_create_apartment", err)
		return
	}

	response.Success(ctx, http.StatusCreated, "apartment_created_successfully", data)
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

	var data *dto.ApartmentResponse
	if data, err = c.apartmentService.Update(ctx, tokenKeys, id, &req); err != nil {
		if errors.Is(err, service_error.ErrApartmentNotFound) {
			response.Error(ctx, http.StatusNotFound, "not_found_apartment", err)
			return
		}
		response.Error(ctx, http.StatusInternalServerError, "failed_to_update_apartment", err)
		return
	}

	response.Success(ctx, http.StatusOK, "apartment_updated_successfully", data)
}

func (c *ApartmentController) Delete(ctx *gin.Context) {
	idParam := ctx.Param("apartment_id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid_apartment_id", err)
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
	if err := c.apartmentService.Delete(ctx, tokenKeys, id); err != nil {
		if errors.Is(err, service_error.ErrApartmentNotFound) {
			response.Error(ctx, http.StatusNotFound, "not_found_apartment", err)
			return
		}
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
	apartment, err := c.apartmentService.GetByID(ctx, tokenKeys, id)
	if err != nil {
		switch err {
		case service_error.ErrApartmentNotFound:
			response.Error(ctx, http.StatusNotFound, "apartment_not_found", err)
			return

		case service_error.ErrApartmentUnauthorizedAccess:
			response.Error(ctx, http.StatusForbidden, "apartment_unauthorized_access", err)
			return

		default:
			response.Error(ctx, http.StatusInternalServerError, "internal_server_error", err)
			return
		}

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
	apartment, err := c.apartmentService.GetWithUsers(ctx, tokenKeys, id)
	if err != nil {
		switch err {

		case service_error.ErrApartmentNotFound:
			response.Error(ctx, http.StatusNotFound, "apartment_not_found", err)
			return

		case service_error.ErrApartmentUnauthorizedAccess:
			response.Error(ctx, http.StatusForbidden, "apartment_unauthorized_access", err)
			return

		default:
			response.Error(ctx, http.StatusInternalServerError, "internal_server_error", err)
			return
		}

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

	apartment, err := c.apartmentService.GetWithRules(ctx, tokenKeys, id)
	if err != nil {
		switch err {

		case service_error.ErrApartmentNotFound:
			response.Error(ctx, http.StatusNotFound, "apartment_not_found", err)
			return

		case service_error.ErrApartmentUnauthorizedAccess:
			response.Error(ctx, http.StatusForbidden, "apartment_unauthorized_access", err)
			return

		default:
			response.Error(ctx, http.StatusInternalServerError, "internal_server_error", err)
			return
		}

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

	apartment, err := c.apartmentService.GetWithAnnouncements(ctx, tokenKeys, id)
	if err != nil {
		switch err {

		case service_error.ErrApartmentNotFound:
			response.Error(ctx, http.StatusNotFound, "apartment_not_found", err)
			return

		case service_error.ErrApartmentUnauthorizedAccess:
			response.Error(ctx, http.StatusForbidden, "apartment_unauthorized_access", err)
			return

		default:
			response.Error(ctx, http.StatusInternalServerError, "internal_server_error", err)
			return
		}

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

	apartment, err := c.apartmentService.GetWithInviteCodes(ctx, tokenKeys, id)
	if err != nil {
		switch err {

		case service_error.ErrApartmentNotFound:
			response.Error(ctx, http.StatusNotFound, "apartment_not_found", err)
			return

		case service_error.ErrApartmentUnauthorizedAccess:
			response.Error(ctx, http.StatusForbidden, "apartment_unauthorized_access", err)
			return

		default:
			response.Error(ctx, http.StatusInternalServerError, "internal_server_error", err)
			return
		}

	}

	response.Success(ctx, http.StatusOK, "apartment_invite_codes_fetched_successfully", apartment)
}

func (c *ApartmentController) GetByIDWithTickets(ctx *gin.Context) {
	idParam := ctx.Param("apartment_id")

	id, err := uuid.Parse(idParam)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, "invalid_apartment_id", err)
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

	apartment, err := c.apartmentService.GetWithTickets(ctx, tokenKeys, id)
	if err != nil {
		switch err {

		case service_error.ErrApartmentNotFound:
			response.Error(ctx, http.StatusNotFound, "apartment_not_found", err)
			return

		case service_error.ErrApartmentUnauthorizedAccess:
			response.Error(ctx, http.StatusForbidden, "apartment_unauthorized_access", err)
			return

		default:
			response.Error(ctx, http.StatusInternalServerError, "internal_server_error", err)
			return
		}

	}

	response.Success(ctx, http.StatusOK, "apartment_rules_fetched_successfully", apartment)
}
