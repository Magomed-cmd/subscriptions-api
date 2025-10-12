package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func NewConnection(dsn string, logger *zap.Logger) (*pgxpool.Pool, error) {
	logger.Debug("Connecting to PostgreSQL", zap.String("dsn", dsn))

	dbPool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		logger.Error("Failed to connect to Postgres", zap.Error(err))
		return nil, err
	}

	if err := dbPool.Ping(context.Background()); err != nil {
		dbPool.Close()
		logger.Error("Failed to ping Postgres", zap.Error(err))
		return nil, err
	}

	logger.Info("PostgresSQL connection established")
	return dbPool, nil
}
