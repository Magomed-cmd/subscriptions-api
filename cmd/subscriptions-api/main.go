// @title Subscriptions API
// @version 1.0
// @description REST API для управления подписками пользователей
// @host localhost:8080
// @BasePath /api/v1
package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	_ "subscriptions-api/docs"
	"subscriptions-api/internal/config"
	"subscriptions-api/internal/db/postgres"
	"subscriptions-api/internal/handler"
	"subscriptions-api/internal/logger"
	"subscriptions-api/internal/repository"
	"subscriptions-api/internal/router"
	"subscriptions-api/internal/service"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.LoadConfig("./configs")
	if err != nil {
		panic("failed to load config")
	}

	log, err := logger.New(cfg.App.LogLevel)
	if err != nil {
		panic("failed to initialize logger")
	}
	defer log.Sync()

	db, err := postgres.NewConnection(cfg.GetDSN(), log)
	if err != nil {
		panic("failed to initialize database connection")
	}
	defer db.Close()

	repo := repository.NewSubscriptionRepository(db, log)
	svc := service.NewSubscriptionService(repo, log)
	subHandler := handler.NewSubscriptionHandler(svc, log)

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	router.InitRoutes(r, subHandler)

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{Addr: addr, Handler: r}

	go func() {
		log.Info("server starting", zap.String("host", cfg.Server.Host), zap.Int("port", cfg.Server.Port))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("server error", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Error("server forced to shutdown", zap.Error(err))
	}
	log.Info("server exited properly")
}
