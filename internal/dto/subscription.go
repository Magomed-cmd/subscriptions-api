package dto

import (
	"time"

	"subscriptions-api/internal/domain/entity"

	"github.com/google/uuid"
)

type CreateSubscriptionRequest struct {
	ServiceName string     `json:"service_name" binding:"required"`
	Price       int64      `json:"price" binding:"required,gte=0"`
	UserID      string     `json:"user_id" binding:"required,uuid"`
	StartDate   time.Time  `json:"start_date" binding:"required"`
	EndDate     *time.Time `json:"end_date,omitempty"`
}

type UpdateSubscriptionRequest struct {
	ServiceName *string    `json:"service_name,omitempty"`
	Price       *int64     `json:"price,omitempty" binding:"omitempty,gte=0"`
	EndDate     *time.Time `json:"end_date,omitempty"`
}

type SubscriptionFilterRequest struct {
	UserID      *string `form:"user_id"`
	ServiceName *string `form:"service_name"`
	StartDate   *string `form:"start_date"`
	EndDate     *string `form:"end_date"`
}

type TotalCostRequest struct {
	UserID      *string `form:"user_id"`
	ServiceName *string `form:"service_name"`
	PeriodStart string  `form:"period_start" binding:"required"`
	PeriodEnd   string  `form:"period_end" binding:"required"`
}

type SubscriptionResponse struct {
	ID          int64      `json:"id"`
	ServiceName string     `json:"service_name"`
	Price       int64      `json:"price"`
	UserID      string     `json:"user_id"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     *time.Time `json:"end_date,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func ToEntitySubscription(req *CreateSubscriptionRequest) (*entity.Subscription, error) {
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, err
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

	if req.UserID != nil {
		uid, err := uuid.Parse(*req.UserID)
		if err != nil {
			return nil, err
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
