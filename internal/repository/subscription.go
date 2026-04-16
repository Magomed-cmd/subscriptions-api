package repository

import (
	"context"
	"errors"
	"subscriptions-api/internal/domain/entity"
	domainerrors "subscriptions-api/internal/domain/errors"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type SubscriptionRepository struct {
	db     *pgxpool.Pool
	psql   squirrel.StatementBuilderType
	logger *zap.Logger
}

func NewSubscriptionRepository(db *pgxpool.Pool, logger *zap.Logger) *SubscriptionRepository {
	return &SubscriptionRepository{
		db:     db,
		psql:   squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
		logger: logger,
	}
}

func (r *SubscriptionRepository) Create(ctx context.Context, sub *entity.Subscription) (int64, error) {
	r.logger.Info("Creating subscription")

	query := `
		INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	var id int64
	if err := r.db.QueryRow(ctx, query,
		sub.ServiceName,
		sub.Price,
		sub.UserID,
		sub.StartDate,
		sub.EndDate,
	).Scan(&id); err != nil {
		mapped := domainerrors.MapPgError(err)
		if !errors.Is(mapped, err) {
			r.logger.Warn("pg error mapped", zap.Error(mapped))
			return 0, mapped
		}
		r.logger.Error("failed to insert subscription", zap.Error(err))
		return 0, domainerrors.ErrInternal
	}

	return id, nil
}

func (r *SubscriptionRepository) GetByID(ctx context.Context, id int64) (*entity.Subscription, error) {
	r.logger.Info("Getting subscription by ID", zap.Int64("id", id))

	query := `
		SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at
		FROM subscriptions
		WHERE id = $1
	`

	var sub entity.Subscription
	err := r.db.QueryRow(ctx, query, id).Scan(
		&sub.ID, &sub.ServiceName,
		&sub.Price, &sub.UserID,
		&sub.StartDate, &sub.EndDate,
		&sub.CreatedAt, &sub.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.logger.Warn("subscription not found", zap.Int64("id", id))
			return nil, domainerrors.ErrNotFound
		}

		mapped := domainerrors.MapPgError(err)
		if !errors.Is(mapped, err) {
			r.logger.Warn("pg error mapped", zap.Error(mapped))
			return nil, mapped
		}

		r.logger.Error("failed to get subscription", zap.Error(err))
		return nil, domainerrors.ErrInternal
	}

	return &sub, nil
}

func (r *SubscriptionRepository) Update(ctx context.Context, sub *entity.Subscription) error {
	r.logger.Info("Updating subscription", zap.Int64("id", sub.ID))

	query := `
		UPDATE subscriptions
		SET service_name = $1,
		    price = $2,
		    end_date = $3,
		    updated_at = NOW()
		WHERE id = $4
	`

	cmd, err := r.db.Exec(ctx, query,
		sub.ServiceName,
		sub.Price,
		sub.EndDate,
		sub.ID,
	)
	if err != nil {
		mapped := domainerrors.MapPgError(err)
		if !errors.Is(mapped, err) {
			r.logger.Warn("pg error mapped", zap.Error(mapped))
			return mapped
		}
		r.logger.Error("failed to update subscription", zap.Error(err))
		return domainerrors.ErrInternal
	}

	if cmd.RowsAffected() == 0 {
		r.logger.Warn("subscription not found for update", zap.Int64("id", sub.ID))
		return domainerrors.ErrNotFound
	}
	r.logger.Info("subscription updated successfully", zap.Int64("id", sub.ID))

	return nil
}

func (r *SubscriptionRepository) Delete(ctx context.Context, id int64) error {
	r.logger.Info("Deleting subscription", zap.Int64("id", id))

	query := `DELETE FROM subscriptions WHERE id = $1`

	cmd, err := r.db.Exec(ctx, query, id)
	if err != nil {
		mapped := domainerrors.MapPgError(err)
		if !errors.Is(mapped, err) {
			r.logger.Warn("pg error mapped", zap.Error(mapped))
			return mapped
		}

		r.logger.Error("failed to delete subscription", zap.Error(err))
		return domainerrors.ErrInternal
	}

	if cmd.RowsAffected() == 0 {
		r.logger.Warn("subscription not found for delete", zap.Int64("id", id))
		return domainerrors.ErrNotFound
	}

	r.logger.Info("subscription deleted successfully", zap.Int64("id", id))
	return nil
}

func (r *SubscriptionRepository) List(ctx context.Context, filters *entity.SubscriptionFilter) ([]entity.Subscription, error) {
	r.logger.Info("Listing subscriptions")

	query := r.psql.
		Select("id", "service_name", "price", "user_id", "start_date", "end_date", "created_at", "updated_at").
		From("subscriptions")
	query = applySubscriptionFilters(query, filters)

	sql, args, err := query.ToSql()
	if err != nil {
		r.logger.Error("failed to build list query", zap.Error(err))
		return nil, domainerrors.ErrInternal
	}

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		mapped := domainerrors.MapPgError(err)
		if !errors.Is(mapped, err) {
			r.logger.Warn("pg error mapped", zap.Error(mapped))
			return nil, mapped
		}
		r.logger.Error("failed to query subscriptions", zap.Error(err))
		return nil, domainerrors.ErrInternal
	}
	defer rows.Close()

	subs := make([]entity.Subscription, 0)
	for rows.Next() {
		var s entity.Subscription
		if err := rows.Scan(
			&s.ID, &s.ServiceName, &s.Price, &s.UserID,
			&s.StartDate, &s.EndDate, &s.CreatedAt, &s.UpdatedAt,
		); err != nil {
			r.logger.Error("failed to scan subscription row", zap.Error(err))
			return nil, domainerrors.ErrInternal
		}
		subs = append(subs, s)
	}

	if rows.Err() != nil {
		r.logger.Error("rows iteration error", zap.Error(rows.Err()))
		return nil, domainerrors.ErrInternal
	}

	r.logger.Info("subscriptions listed successfully", zap.Int("count", len(subs)))
	return subs, nil
}

func (r *SubscriptionRepository) TotalAmount(ctx context.Context, filters *entity.SubscriptionFilter) (int64, error) {
	r.logger.Info("Calculating total subscription amount")

	query := r.psql.Select("COALESCE(SUM(price), 0)").From("subscriptions")
	query = applySubscriptionFilters(query, filters)

	sql, args, err := query.ToSql()
	if err != nil {
		r.logger.Error("failed to build total amount query", zap.Error(err))
		return 0, domainerrors.ErrInternal
	}

	r.logger.Debug("Executing SQL", zap.String("sql", sql), zap.Any("args", args))

	var total int64
	if err := r.db.QueryRow(ctx, sql, args...).Scan(&total); err != nil {
		mapped := domainerrors.MapPgError(err)
		if !errors.Is(mapped, err) {
			r.logger.Warn("pg error mapped", zap.Error(mapped))
			return 0, mapped
		}
		r.logger.Error("failed to calculate total amount", zap.Error(err))
		return 0, domainerrors.ErrInternal
	}

	r.logger.Info("total amount calculated", zap.Int64("total", total))
	return total, nil
}

func applySubscriptionFilters(query squirrel.SelectBuilder, f *entity.SubscriptionFilter) squirrel.SelectBuilder {
	if f == nil {
		return query
	}

	if f.UserID != nil {
		query = query.Where(squirrel.Eq{"user_id": *f.UserID})
	}
	if f.ServiceName != nil {
		query = query.Where(squirrel.Eq{"service_name": *f.ServiceName})
	}
	if f.StartDate != nil && f.EndDate != nil {
		query = query.Where(squirrel.And{
			squirrel.LtOrEq{"start_date": *f.EndDate},
			squirrel.GtOrEq{"end_date": *f.StartDate},
		})
	} else if f.StartDate != nil {
		query = query.Where(squirrel.GtOrEq{"start_date": *f.StartDate})
	} else if f.EndDate != nil {
		query = query.Where(squirrel.LtOrEq{"end_date": *f.EndDate})
	}

	return query
}
