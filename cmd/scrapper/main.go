package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go-ItsDianthus-NotificationLink/internal/api/openapi/scrapper_api"
	"go-ItsDianthus-NotificationLink/internal/scrapper/application/handlers"
	"go-ItsDianthus-NotificationLink/internal/scrapper/config"
	"go-ItsDianthus-NotificationLink/internal/scrapper/infrastructure/repo"
	"go-ItsDianthus-NotificationLink/pkg/slogger"
	"log"
)

func main() {
	cfg, err := config.LoadConfig("config/scrapper.yaml")
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	lg := slogger.NewLoggerByEnvironment(cfg.Env)
	lg.Info("Starting Scrapper")

	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	r := repo.NewSubscriptionRepo()
	srv := handlers.NewServer(r)
	scrapper_api.RegisterHandlersWithBaseURL(e, srv, "")

	e.Server.ReadTimeout = cfg.Server.ReadTimeout
	e.Server.WriteTimeout = cfg.Server.WriteTimeout

	addr := cfg.Server.Address()
	lg.Info("Listening on " + addr)
	log.Fatal(e.Start(addr))
}
