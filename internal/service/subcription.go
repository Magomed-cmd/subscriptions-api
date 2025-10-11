package service

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"subscriptions-api/internal/domain/entity"
	errors2 "subscriptions-api/internal/domain/errors"
	"subscriptions-api/internal/dto"
	"time"

	"go.uber.org/zap"
)

type SubscriptionRepository interface {
	Create(ctx context.Context, sub *entity.Subscription) (int64, error)
	GetByID(ctx context.Context, id int64) (*entity.Subscription, error)
	Update(ctx context.Context, sub *entity.Subscription) error
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context, filters *entity.SubscriptionFilter) ([]entity.Subscription, error)
	TotalAmount(ctx context.Context, filters *entity.SubscriptionFilter) (int64, error)
}

type SubscriptionService struct {
	repo   SubscriptionRepository
	logger *zap.Logger
}

func NewSubscriptionService(repo SubscriptionRepository, logger *zap.Logger) *SubscriptionService {
	return &SubscriptionService{
		repo:   repo,
		logger: logger,
	}
}

func (s *SubscriptionService) Create(ctx context.Context, req *dto.CreateSubscriptionRequest) (*dto.SubscriptionResponse, error) {
	s.logger.Info("Creating subscription", zap.String("service_name", req.ServiceName))

	sub, err := dto.ToEntitySubscription(req)
	if err != nil {
		s.logger.Warn("Invalid input DTO", zap.Error(err))
		return nil, errors2.ErrInvalidInput
	}

	id, err := s.repo.Create(ctx, sub)
	if err != nil {
		s.logger.Error("Failed to create subscription", zap.Error(err))
		return nil, err
	}

	sub.ID = id
	s.logger.Info("Subscription created", zap.Int64("id", id))
	return dto.ToResponseSubscription(sub), nil
}

func (s *SubscriptionService) GetByID(ctx context.Context, id int64) (*dto.SubscriptionResponse, error) {
	s.logger.Info("Getting subscription", zap.Int64("id", id))

	if id <= 0 {
		return nil, errors2.ErrInvalidInput
	}

	sub, err := s.repo.GetByID(ctx, id)
	if err != nil {
		s.logger.Error("Failed to get subscription", zap.Error(err))
		return nil, err
	}

	return dto.ToResponseSubscription(sub), nil
}

func (s *SubscriptionService) Update(ctx context.Context, id int64, req *dto.UpdateSubscriptionRequest) (*dto.SubscriptionResponse, error) {
	s.logger.Info("Updating subscription", zap.Int64("id", id))

	if id <= 0 {
		return nil, errors2.ErrInvalidInput
	}

	if req.ServiceName == nil && req.Price == nil && req.EndDate == nil {
		return nil, errors2.ErrNothingToUpdate
	}

	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.ServiceName != nil {
		existing.ServiceName = *req.ServiceName
	}
	if req.Price != nil {
		existing.Price = *req.Price
	}
	if req.EndDate != nil {
		existing.EndDate = req.EndDate
	}

	existing.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, existing); err != nil {
		s.logger.Error("Failed to update subscription", zap.Error(err))
		return nil, err
	}

	return dto.ToResponseSubscription(existing), nil
}

func (s *SubscriptionService) Delete(ctx context.Context, id int64) error {
	s.logger.Info("Deleting subscription", zap.Int64("id", id))

	if id <= 0 {
		return errors2.ErrInvalidInput
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("Failed to delete subscription", zap.Error(err))
		return err
	}

	s.logger.Info("Subscription deleted", zap.Int64("id", id))
	return nil
}

func (s *SubscriptionService) List(ctx context.Context, filterReq *dto.SubscriptionFilterRequest) ([]dto.SubscriptionResponse, error) {
	s.logger.Info("Listing subscriptions")

	filter, err := dto.ToEntityFilter(filterReq)
	if err != nil {
		s.logger.Warn("Invalid filter DTO", zap.Error(err))
		return nil, errors2.ErrInvalidInput
	}

	subs, err := s.repo.List(ctx, filter)
	if err != nil {
		s.logger.Error("Failed to list subscriptions", zap.Error(err))
		return nil, err
	}

	resp := make([]dto.SubscriptionResponse, len(subs))
	for i := range subs {
		resp[i] = *dto.ToResponseSubscription(&subs[i])
	}

	s.logger.Info("Subscriptions listed", zap.Int("count", len(resp)))
	return resp, nil
}

func (s *SubscriptionService) TotalCostForPeriod(ctx context.Context, req *dto.TotalCostRequest) (int64, error) {
	s.logger.Info("Calculating total cost for period",
		zap.String("period_start", req.PeriodStart),
		zap.String("period_end", req.PeriodEnd))

	start, err := parseMonthYear(req.PeriodStart)
	if err != nil {
		return 0, errors2.ErrInvalidInput
	}

	end, err := parseMonthYear(req.PeriodEnd)
	if err != nil {
		return 0, errors2.ErrInvalidInput
	}

	if end.Before(start) {
		return 0, errors2.ErrInvalidInput
	}

	filter, err := dto.ToEntityFilter(&dto.SubscriptionFilterRequest{
		UserID:      req.UserID,
		ServiceName: req.ServiceName,
	})
	if err != nil {
		return 0, errors2.ErrInvalidInput
	}

	filter.StartDate = &start
	filter.EndDate = &end

	total, err := s.repo.TotalAmount(ctx, filter)
	if err != nil {
		s.logger.Error("Failed to calculate total", zap.Error(err))
		return 0, err
	}

	s.logger.Info("Total cost calculated", zap.Int64("total", total))
	return total, nil
}

func parseMonthYear(str string) (time.Time, error) {
	if str == "" {
		return time.Time{}, errors.New("empty date string")
	}

	parts := strings.Split(str, "-")
	if len(parts) != 2 {
		return time.Time{}, errors.New("invalid format, expected MM-YYYY")
	}

	month, err := strconv.Atoi(parts[0])
	if err != nil || month < 1 || month > 12 {
		return time.Time{}, errors.New("invalid month")
	}

	year, err := strconv.Atoi(parts[1])
	if err != nil || year < 1900 || year > 2100 {
		return time.Time{}, errors.New("invalid year")
	}

	return time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC), nil
}
