package commands

import (
	"context"
	"go-ItsDianthus-NotificationLink/internal/bot/application/command_registry"

	"go-ItsDianthus-NotificationLink/internal/bot/domain"
	"go-ItsDianthus-NotificationLink/internal/bot/infrastructure/telegram"
)

type CancelCommand struct {
	Bot      telegram.BotClient
	Registry *command_registry.CommandRegistry
}

func NewMenuCommand(bot telegram.BotClient, reg *command_registry.CommandRegistry) *CancelCommand {
	return &CancelCommand{Bot: bot, Registry: reg}
}

func (c *CancelCommand) Name() string { return "/menu" }

func (c *CancelCommand) Description() string {
	return "Отменить текущую операцию"
}

func (c *CancelCommand) Execute(ctx context.Context, session *domain.UserSession, args []string) error {
	session.CurrentState = domain.StateDefault
	session.ActiveCommand = ""
	session.TempData = nil

	cmds := c.Registry.AllExceptStart()
	kb := telegram.BuildCommandKeyboard(cmds)

	c.Bot.SendMessage(session.UserID,
		"Вы вернулись в главное меню. Выберите любую команду",
		kb,
	)
	return nil
}
