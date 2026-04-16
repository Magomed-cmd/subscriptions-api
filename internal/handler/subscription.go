package handler

import (
	"context"
	"net/http"
	"strconv"
	domainerrors "subscriptions-api/internal/domain/errors"
	"subscriptions-api/internal/dto"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SubscriptionService interface {
	Create(ctx context.Context, req *dto.CreateSubscriptionRequest) (*dto.SubscriptionResponse, error)
	GetByID(ctx context.Context, id int64) (*dto.SubscriptionResponse, error)
	Update(ctx context.Context, id int64, req *dto.UpdateSubscriptionRequest) (*dto.SubscriptionResponse, error)
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, filterReq *dto.SubscriptionFilterRequest) ([]dto.SubscriptionResponse, error)
	TotalCostForPeriod(ctx context.Context, req *dto.TotalCostRequest) (int64, error)
}

type SubscriptionHandler struct {
	svc    SubscriptionService
	logger *zap.Logger
}

func NewSubscriptionHandler(svc SubscriptionService, logger *zap.Logger) *SubscriptionHandler {
	return &SubscriptionHandler{svc: svc, logger: logger}
}

// @Summary Создать подписку
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param body body dto.CreateSubscriptionRequest true "Данные подписки"
// @Success 201 {object} dto.SubscriptionResponse
// @Failure 400 {object} ErrorResponse
// @Router /subscriptions [post]
func (h *SubscriptionHandler) Create(c *gin.Context) {
	req := dto.CreateSubscriptionRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid create request", zap.Error(err))
		JSONError(c, domainerrors.ErrInvalidInput)
		return
	}

	resp, err := h.svc.Create(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("failed to create subscription", zap.Error(err))
		JSONError(c, err)
		return
	}

	c.JSON(http.StatusCreated, resp)
}

// @Summary Получить подписку по ID
// @Tags subscriptions
// @Produce json
// @Param id path int true "ID подписки"
// @Success 200 {object} dto.SubscriptionResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /subscriptions/{id} [get]
func (h *SubscriptionHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.logger.Warn("invalid id param", zap.Error(err))
		JSONError(c, domainerrors.ErrInvalidInput)
		return
	}

	resp, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		h.logger.Error("failed to get subscription", zap.Error(err))
		JSONError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Summary Обновить подписку
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path int true "ID подписки"
// @Param body body dto.UpdateSubscriptionRequest true "Данные для обновления"
// @Success 200 {object} dto.SubscriptionResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /subscriptions/{id} [put]
func (h *SubscriptionHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.logger.Warn("invalid id param", zap.Error(err))
		JSONError(c, domainerrors.ErrInvalidInput)
		return
	}

	req := dto.UpdateSubscriptionRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid update request", zap.Error(err))
		JSONError(c, domainerrors.ErrInvalidInput)
		return
	}

	resp, err := h.svc.Update(c.Request.Context(), id, &req)
	if err != nil {
		h.logger.Error("failed to update subscription", zap.Error(err))
		JSONError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Summary Удалить подписку
// @Tags subscriptions
// @Param id path int true "ID подписки"
// @Success 204
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /subscriptions/{id} [delete]
func (h *SubscriptionHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.logger.Warn("invalid id param", zap.Error(err))
		JSONError(c, domainerrors.ErrInvalidInput)
		return
	}

	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		h.logger.Error("failed to delete subscription", zap.Error(err))
		JSONError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// @Summary Получить список подписок
// @Tags subscriptions
// @Produce json
// @Param user_id query string false "ID пользователя"
// @Param service_name query string false "Название сервиса"
// @Success 200 {array} dto.SubscriptionResponse
// @Failure 400 {object} ErrorResponse
// @Router /subscriptions [get]
func (h *SubscriptionHandler) List(c *gin.Context) {
	filter := dto.SubscriptionFilterRequest{}
	if err := c.ShouldBindQuery(&filter); err != nil {
		h.logger.Warn("invalid filter params", zap.Error(err))
		JSONError(c, domainerrors.ErrInvalidInput)
		return
	}

	resp, err := h.svc.List(c.Request.Context(), &filter)
	if err != nil {
		h.logger.Error("failed to list subscriptions", zap.Error(err))
		JSONError(c, err)
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Summary Получить суммарную стоимость подписок за период
// @Tags subscriptions
// @Produce json
// @Param period_start query string true "Начало периода (2006-01-02T15:04:05Z)"
// @Param period_end query string true "Конец периода (2006-01-02T15:04:05Z)"
// @Param user_id query string false "ID пользователя"
// @Param service_name query string false "Название сервиса"
// @Success 200 {object} map[string]int
// @Failure 400 {object} ErrorResponse
// @Router /subscriptions/total [get]
func (h *SubscriptionHandler) TotalCostForPeriod(c *gin.Context) {
	req := dto.TotalCostRequest{}
	if err := c.ShouldBindQuery(&req); err != nil {
		h.logger.Warn("invalid total cost params", zap.Error(err))
		JSONError(c, domainerrors.ErrInvalidInput)
		return
	}

	if req.PeriodStart == nil || req.PeriodEnd == nil {
		h.logger.Warn("missing required period dates")
		JSONError(c, domainerrors.ErrInvalidInput)
		return
	}

	if req.PeriodEnd.Before(*req.PeriodStart) {
		h.logger.Warn("period end before period start")
		JSONError(c, domainerrors.ErrInvalidInput)
		return
	}

	h.logger.Info("calculating total cost", zap.Any("request", req))

	total, err := h.svc.TotalCostForPeriod(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("failed to calculate total cost", zap.Error(err))
		JSONError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"total": total})
}
