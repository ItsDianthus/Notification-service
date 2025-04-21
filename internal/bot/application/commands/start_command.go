package commands

import (
	"context"
	"fmt"
	"go-ItsDianthus-NotificationLink/internal/bot/application/clients"
	"go-ItsDianthus-NotificationLink/internal/bot/domain"
	"go-ItsDianthus-NotificationLink/internal/bot/infrastructure/telegram"
)

type StartCommand struct {
	Bot            telegram.BotClient
	ScrapperClient clients.ScrapperClient
}

func NewStartCommand(bot telegram.BotClient, scr clients.ScrapperClient) *StartCommand {
	return &StartCommand{Bot: bot, ScrapperClient: scr}
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
		if err := c.ScrapperClient.RegisterChat(ctx, session.UserID); err != nil {
			c.Bot.SendMessage(session.UserID,
				"Ошибка при регистрации: "+err.Error(),
				nil,
			)
			return fmt.Errorf("scrapper register chat: %w", err)
		}

		session.IsRegistered = true
		c.Bot.SendMessage(session.UserID,
			"Добро пожаловать! Вы зарегистрированы. Доступные команды:",
			fullKb,
		)
	}
	return nil
}
