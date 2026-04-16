package handler

import (
	"errors"
	"net/http"
	domainerrors "subscriptions-api/internal/domain/errors"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func JSONError(c *gin.Context, err error) {
	status, msg := mapErrorToResponse(err)
	c.JSON(status, ErrorResponse{Error: msg})
}

func mapErrorToResponse(err error) (int, string) {
	switch {
	case err == nil:
		return http.StatusOK, ""

	case errors.Is(err, domainerrors.ErrNotFound):
		return http.StatusNotFound, "resource not found"

	case errors.Is(err, domainerrors.ErrInvalidInput),
		errors.Is(err, domainerrors.ErrInvalidArgument),
		errors.Is(err, domainerrors.ErrNullViolation),
		errors.Is(err, domainerrors.ErrCheckViolation),
		errors.Is(err, domainerrors.ErrInvalidUUID),
		errors.Is(err, domainerrors.ErrInvalidDateFormat),
		errors.Is(err, domainerrors.ErrForeignKeyViolation),
		errors.Is(err, domainerrors.ErrNothingToUpdate):
		return http.StatusBadRequest, "invalid request data"

	case errors.Is(err, domainerrors.ErrDuplicate):
		return http.StatusConflict, "duplicate record"

	case errors.Is(err, domainerrors.ErrInternal):
		return http.StatusInternalServerError, "internal server error"

	default:
		return http.StatusInternalServerError, "unexpected error"
	}
}
