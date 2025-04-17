package main

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"go-ItsDianthus-NotificationLink/internal/bot/application"
	"go-ItsDianthus-NotificationLink/internal/bot/application/commands"
	"go-ItsDianthus-NotificationLink/internal/bot/config"
	"go-ItsDianthus-NotificationLink/internal/bot/infrastructure/repo"
	"go-ItsDianthus-NotificationLink/internal/bot/infrastructure/telegram"
	"go-ItsDianthus-NotificationLink/pkg/slogger"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func main() {
	// Загрузка .env
	if err := godotenv.Load(); err != nil {
		fmt.Printf(".env not found or couldn't load: %v", err)
	}

	// Чтение конфигурации из YAML + env-переменные
	cfg, err := config.LoadConfig("config/bot.yaml")
	if err != nil {
		log.Fatalf("cannot load config: %v", err)
	}

	// Настройка логгера по env (local/prod)
	lg := slogger.NewLoggerByEnvironment(cfg.Env)
	lg.Info("Starting bot service", slog.String("env", cfg.Env))

	// Создание Telegram API клиента
	tgAPI, err := tgbotapi.NewBotAPI(cfg.Telegram.Token)
	if err != nil {
		slog.Error("Failed to create Telegram API client", slogger.ErrAttr(err))
		os.Exit(1)
	}
	botClient := telegram.NewTgBotClient(tgAPI)

	// HTTP-клиент для Scrapper-сервиса
	// scrapperClient := application.NewScrapperHTTPClient(cfg.Scrapper.Address(), cfg.Scrapper.Timeout)

	// Инициализация in-memory репозитория сессий
	sessRepo := repo.NewInMemorySessionRepo()

	// Регистрация команд
	r := application.NewCommandRegistry()
	r.Register(commands.NewStartCommand(botClient))
	//r.Register(commands.NewHelpCommand(botClient, r))
	//r.Register(commands.NewTrackCommand(botClient, scrapperClient))
	//r.Register(commands.NewUntrackCommand(botClient, scrapperClient))
	//r.Register(commands.NewListCommand(botClient, scrapperClient))

	// Настройка горутин на получение и обработку сообщений
	updates := tgAPI.GetUpdatesChan(tgbotapi.NewUpdate(0))
	proc := application.NewMessageProcessor(botClient, sessRepo, r, updates, lg)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go proc.ProcessUpdates(ctx)

	<-ctx.Done()
	lg.Info("Shutdown signal received, exiting...")
}
