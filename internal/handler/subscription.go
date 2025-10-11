package handler

import (
	"net/http"
	"strconv"
	"subscriptions-api/internal/dto"
	"subscriptions-api/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type SubscriptionHandler struct {
	svc    *service.SubscriptionService
	logger *zap.Logger
}

func NewSubscriptionHandler(svc *service.SubscriptionService, logger *zap.Logger) *SubscriptionHandler {
	return &SubscriptionHandler{svc: svc, logger: logger}
}

func (h *SubscriptionHandler) RegisterRoutes(r *gin.Engine) {
	v1 := r.Group("/api/v1/subscriptions")
	{
		v1.POST("", h.Create)
		v1.GET("/:id", h.GetByID)
		v1.PUT("/:id", h.Update)
		v1.DELETE("/:id", h.Delete)
		v1.GET("", h.List)
		v1.GET("/total", h.TotalCostForPeriod)
	}
}

func (h *SubscriptionHandler) Create(c *gin.Context) {
	req := dto.CreateSubscriptionRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid create request", zap.Error(err))
		JSONError(c, err)
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

func (h *SubscriptionHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.logger.Warn("invalid id param", zap.Error(err))
		JSONError(c, err)
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

func (h *SubscriptionHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.logger.Warn("invalid id param", zap.Error(err))
		JSONError(c, err)
		return
	}

	req := dto.UpdateSubscriptionRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid update request", zap.Error(err))
		JSONError(c, err)
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

func (h *SubscriptionHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		h.logger.Warn("invalid id param", zap.Error(err))
		JSONError(c, err)
		return
	}

	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		h.logger.Error("failed to delete subscription", zap.Error(err))
		JSONError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *SubscriptionHandler) List(c *gin.Context) {
	filter := dto.SubscriptionFilterRequest{}
	if err := c.ShouldBindQuery(&filter); err != nil {
		h.logger.Warn("invalid filter params", zap.Error(err))
		JSONError(c, err)
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

func (h *SubscriptionHandler) TotalCostForPeriod(c *gin.Context) {
	req := dto.TotalCostRequest{}
	if err := c.ShouldBindQuery(&req); err != nil {
		h.logger.Warn("invalid total cost params", zap.Error(err))
		JSONError(c, err)
		return
	}

	total, err := h.svc.TotalCostForPeriod(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("failed to calculate total cost", zap.Error(err))
		JSONError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"total": total})
}
