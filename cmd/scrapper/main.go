package main

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go-ItsDianthus-NotificationLink/internal/api/openapi/scrapper_api"
	"go-ItsDianthus-NotificationLink/internal/bot/config"
	"go-ItsDianthus-NotificationLink/internal/scrapper/application/handlers"
	"go-ItsDianthus-NotificationLink/internal/scrapper/infrastructure/repo"
	"go-ItsDianthus-NotificationLink/pkg/slogger"
	"golang.org/x/sync/errgroup"
	"log"
	"log/slog"
	"time"
)

func main() {
	// 1) Конфиг
	cfg, err := config.LoadConfig("config/scrapper.yaml")
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}

	// 2) Логгер
	lg := slogger.NewLoggerByEnvironment(cfg.Env)
	lg.Info("Starting scrapper", slog.String("env", cfg.Env))

	// 3) Echo + middleware
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:     true,
		LogMethod:  true,
		LogStatus:  true,
		LogLatency: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			lg.Info("HTTP request",
				slog.String("method", v.Method),
				slog.String("uri", v.URI),
				slog.Int("status", v.Status),
				slog.String("latency", v.Latency.String()),
			)
			return nil
		},
	}))

	// 4) Репозиторий и HTTP‑хендлер
	repo := repo.NewSubscriptionRepo()
	srv := handlers.NewServer(repo)
	scrapper_api.RegisterHandlersWithBaseURL(e, srv, "")

	// 5) Серверные таймауты
	e.Server.ReadTimeout = cfg.Server.ReadTimeout
	e.Server.WriteTimeout = cfg.Server.WriteTimeout

	// 6) Запуск API и шедулера параллельно
	ctx, cancel := context.WithCancel(context.Background())
	g, ctx := errgroup.WithContext(ctx)

	// a) HTTP‑сервер
	g.Go(func() error {
		lg.Info("HTTP listening on", slog.String("addr", cfg.Server.Address()))
		return e.Start(cfg.Server.Address())
	})

	// b) Планировщик
	g.Go(func() error {
		ticker := time.NewTicker(cfg.Scheduler.Interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-ticker.C:
				// здесь вызывать fetch‑модуль:
				//   1) repo.ListAllChatIDs()
				//   2) для каждого chatID → repo.ListLinks(chatID)
				//   3) групповой fetch GitHub/SO с таймаутом cfg.FetchGitHub.Timeout
				//   4) при новых апдейтах → POST /updates в Bot‑сервис через HTTP‑клиент с базой cfg.Bot.BaseURL() и timeout cfg.Bot.Timeout
			}
		}
	})

	// ждём Ctrl+C
	<-interruptSignal()
	cancel()
	if err := g.Wait(); err != nil && err != context.Canceled {
		lg.Error("Shutdown with error", slog.String("err", err.Error()))
	}
	lg.Info("Scrapper service stopped")
}
