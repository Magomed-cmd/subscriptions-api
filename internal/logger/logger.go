package logger

import (
	"go.uber.org/zap"
)

func New(logLevel string) (*zap.Logger, error) {
	var cfg zap.Config

	if logLevel == "debug" {
		cfg = zap.NewDevelopmentConfig()
	} else {
		cfg = zap.NewProductionConfig()
	}

	return cfg.Build()
}
