package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v3"
	"github.com/robertt3kuk/xiaoma-test-task/config"
	"github.com/robertt3kuk/xiaoma-test-task/init/httpserver"
	"github.com/robertt3kuk/xiaoma-test-task/init/logger"
	"github.com/robertt3kuk/xiaoma-test-task/init/postgres"
	v1 "github.com/robertt3kuk/xiaoma-test-task/internal/delivery/http/v1"
	"github.com/robertt3kuk/xiaoma-test-task/internal/service"
)

func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)
	pg, err := postgres.New(cfg.PG.URL, postgres.MaxPoolSize(cfg.PoolMax))
	if err != nil {
		panic(err)
	}

	repo := service.NewRepo(pg)
	service := service.New(repo)

	handler := fiber.New()
	v1.NewRouter(handler, l, service)

	httpServer := httpserver.New(handler.Handler(), cfg.HTTP.Port)
	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
