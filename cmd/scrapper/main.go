package main

import (
	"github.com/labstack/echo/v4"
	"go-ItsDianthus-NotificationLink/internal/api/openapi/scrapper_api"
	"go-ItsDianthus-NotificationLink/internal/scrapper/application/handlers"
	"go-ItsDianthus-NotificationLink/internal/scrapper/infrastructure/repo"
	"time"
)

func main() {
	e := echo.New()

	repo := repo.NewSubscriptionRepo()
	srv := handlers.NewServer(repo)

	scrapper_api.RegisterHandlersWithBaseURL(e, srv, "")

	e.Server.ReadTimeout = 5 * time.Second
	e.Server.WriteTimeout = 5 * time.Second

	// Сюда логгер
}
