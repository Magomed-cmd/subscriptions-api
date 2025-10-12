package dto

import (
	"time"

	"subscriptions-api/internal/domain/entity"
	errors2 "subscriptions-api/internal/domain/errors"

	"github.com/google/uuid"
)

type CreateSubscriptionRequest struct {
	ServiceName string     `json:"service_name" binding:"required" example:"Netflix"`
	Price       int64      `json:"price" binding:"required,gt=0" example:"499"`
	UserID      string     `json:"user_id" binding:"required,uuid" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	StartDate   time.Time  `json:"start_date" binding:"required" example:"2025-07-01T00:00:00Z"`
	EndDate     *time.Time `json:"end_date,omitempty" example:"2025-12-01T00:00:00Z"`
}

type UpdateSubscriptionRequest struct {
	ServiceName *string    `json:"service_name,omitempty" example:"Yandex Plus"`
	Price       *int64     `json:"price,omitempty" binding:"omitempty,gt=0" example:"599"`
	EndDate     *time.Time `json:"end_date,omitempty" example:"2025-12-31T00:00:00Z"`
}

type TotalCostRequest struct {
	UserID      *string    `form:"user_id" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	ServiceName *string    `form:"service_name" example:"Netflix"`
	PeriodStart *time.Time `form:"period_start" binding:"required" time_format:"2006-01-02T15:04:05Z" example:"2025-10-12T00:00:00Z"`
	PeriodEnd   *time.Time `form:"period_end" binding:"required" time_format:"2006-01-02T15:04:05Z" example:"2025-12-12T00:00:00Z"`
}

type SubscriptionResponse struct {
	ID          int64      `json:"id" example:"1"`
	ServiceName string     `json:"service_name" example:"Netflix"`
	Price       int64      `json:"price" example:"499"`
	UserID      string     `json:"user_id" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	StartDate   time.Time  `json:"start_date" example:"2025-07-01T00:00:00Z"`
	EndDate     *time.Time `json:"end_date,omitempty" example:"2025-12-01T00:00:00Z"`
	CreatedAt   time.Time  `json:"created_at" example:"2025-10-12T19:33:28Z"`
	UpdatedAt   time.Time  `json:"updated_at" example:"2025-10-12T19:33:28Z"`
}

type SubscriptionFilterRequest struct {
	UserID      *string `form:"user_id"`
	ServiceName *string `form:"service_name"`
	StartDate   *string `form:"start_date"`
	EndDate     *string `form:"end_date"`
}

func ToEntitySubscription(req *CreateSubscriptionRequest) (*entity.Subscription, error) {
	if req == nil {
		return nil, errors2.ErrInvalidInput
	}
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, errors2.ErrInvalidInput
	}

	return &entity.Subscription{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      userID,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}, nil
}

func ToResponseSubscription(sub *entity.Subscription) *SubscriptionResponse {
	if sub == nil {
		return nil
	}
	return &SubscriptionResponse{
		ID:          sub.ID,
		ServiceName: sub.ServiceName,
		Price:       sub.Price,
		UserID:      sub.UserID.String(),
		StartDate:   sub.StartDate,
		EndDate:     sub.EndDate,
		CreatedAt:   sub.CreatedAt,
		UpdatedAt:   sub.UpdatedAt,
	}
}

func ToEntityFilter(req *SubscriptionFilterRequest) (*entity.SubscriptionFilter, error) {
	filter := &entity.SubscriptionFilter{}
	if req == nil {
		return filter, nil
	}

	if req.UserID != nil {
		uid, err := uuid.Parse(*req.UserID)
		if err != nil {
			return nil, errors2.ErrInvalidInput
		}
		filter.UserID = &uid
	}

	filter.ServiceName = req.ServiceName

	parseDate := func(s *string) *time.Time {
		if s == nil || *s == "" {
			return nil
		}
		t, err := time.Parse("2006-01-02", *s)
		if err != nil {
			return nil
		}
		return &t
	}

	filter.StartDate = parseDate(req.StartDate)
	filter.EndDate = parseDate(req.EndDate)

	return filter, nil
}

func ToEntityFilterTotalCost(req *TotalCostRequest) (*entity.SubscriptionFilter, error) {
	filter := &entity.SubscriptionFilter{}
	if req == nil {
		return filter, nil
	}
	
	if req.UserID != nil {
		uid, err := uuid.Parse(*req.UserID)
		if err != nil {
			return nil, errors2.ErrInvalidInput
		}
		filter.UserID = &uid
	}

	filter.ServiceName = req.ServiceName

	filter.StartDate = req.PeriodStart
	filter.EndDate = req.PeriodEnd

	return filter, nil
}
