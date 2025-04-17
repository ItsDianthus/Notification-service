package commands

import (
	"context"
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

	allCmds := []string{"/help", "/track", "/untrack", "/list"}
	fullKb := telegram.BuildCommandKeyboard(allCmds)

	if session.IsRegistered {
		c.Bot.SendMessage(session.UserID,
			"Вы уже зарегистрированы. Доступные команды:",
			fullKb,
		)
	} else {
		session.IsRegistered = true
		c.Bot.SendMessage(session.UserID,
			"Добро пожаловать! Вы зарегистрированы. Доступные команды:",
			fullKb,
		)
	}
	return nil
}
