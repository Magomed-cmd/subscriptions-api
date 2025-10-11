package handler

import (
	"errors"
	"net/http"
	errors2 "subscriptions-api/internal/domain/errors"

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

	case errors.Is(err, errors2.ErrNotFound):
		return http.StatusNotFound, "resource not found"

	case errors.Is(err, errors2.ErrInvalidInput),
		errors.Is(err, errors2.ErrInvalidArgument),
		errors.Is(err, errors2.ErrNullViolation),
		errors.Is(err, errors2.ErrCheckViolation),
		errors.Is(err, errors2.ErrInvalidUUID),
		errors.Is(err, errors2.ErrInvalidDateFormat),
		errors.Is(err, errors2.ErrForeignKeyViolation):
		return http.StatusBadRequest, "invalid request data"

	case errors.Is(err, errors2.ErrDuplicate):
		return http.StatusConflict, "duplicate record"

	case errors.Is(err, errors2.ErrInternal):
		return http.StatusInternalServerError, "internal server error"

	default:
		return http.StatusInternalServerError, "unexpected error"
	}
}
