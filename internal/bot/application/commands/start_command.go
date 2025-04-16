package commands

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go-ItsDianthus-NotificationLink/internal/bot/domain"
	"go-ItsDianthus-NotificationLink/internal/bot/infrastructure/telegram"
)

type StartCommand struct {
	Bot telegram.BotClient
}

func NewStartCommand(bot telegram.BotClient) *StartCommand {
	return &StartCommand{Bot: bot}
}

func (c *StartCommand) Name() string        { return "/start" }
func (c *StartCommand) Description() string { return "Регистрация пользователя" }

func (c *StartCommand) Execute(ctx context.Context, session *domain.UserSession, args []string) error {
	session.CurrentState = domain.StateDefault
	session.ActiveCommand = ""

	startKb := tgbotapi.NewReplyKeyboard(tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("/start"),
	))
	c.Bot.SendMessage(session.UserID,
		"Привет! Нажмите /start для начала.",
		startKb,
	)

	// РЕГА СЕССИИ в хранилище инмемори

	allCmds := []string{"/help", "/track", "/untrack", "/list"}
	fullKb := telegram.BuildCommandKeyboard(allCmds)

	c.Bot.SendMessage(session.UserID,
		"Вы зарегистрированы! Доступные команды:",
		fullKb,
	)
	return nil
}
