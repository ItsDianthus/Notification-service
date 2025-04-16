package main

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"go-ItsDianthus-NotificationLink/internal/bot/application"
	"go-ItsDianthus-NotificationLink/internal/bot/application/commands"
	"go-ItsDianthus-NotificationLink/internal/bot/domain"
	"go-ItsDianthus-NotificationLink/internal/bot/handler"
	"go-ItsDianthus-NotificationLink/internal/bot/infrastructure/repo"
	"go-ItsDianthus-NotificationLink/internal/bot/infrastructure/telegram"
	"go-ItsDianthus-NotificationLink/pkg/log_package"
	"log/slog"
	"os"
)

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Ошибка загрузки .env файла: %v\n", err)
	}

	env := getEnv("APP_ENV", log_package.EnvLocal)

	logger := log_package.NewLoggerByEnvironment(env)
	logger.Info("Запуск сервиса бота", slog.String("env", env))

	config := domain.BotConfig{
		Token:           getEnv("TELEGRAM_BOT_TOKEN", ""),
		ScrapperBaseURL: getEnv("SCRAPPER_BASE_URL", ""),
	}
	logger.Info("Конфигурация бота загружена",
		slog.String("token", config.Token),
		slog.String("scrapper_url", config.ScrapperBaseURL),
	)

	// Создаём Telegram-клиент
	tg, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		logger.Info("Не удалось создать Telegram API: %v", err)
	}
	updates := tg.GetUpdatesChan(tgbotapi.NewUpdate(0))

	// Инициализируем репозиторий сессий, регистр и команды
	repo := repo.NewInMemorySessionRepo()
	registry := application.NewCommandRegistry()
	botClient := telegram.NewTgBotClient(tg)

	// Регистрируем команды
	registry.Register(commands.NewStartCommand(botClient))

	// Запускаем процессор
	processor := handler.NewMessageProcessor(botClient, repo, registry, updates)
	ctx := context.Background()
	processor.ProcessUpdates(ctx)
}
