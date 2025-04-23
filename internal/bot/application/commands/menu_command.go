package commands

import (
	"context"
	"go-ItsDianthus-NotificationLink/internal/bot/application/command_handling"

	"go-ItsDianthus-NotificationLink/internal/bot/domain"
	"go-ItsDianthus-NotificationLink/internal/bot/infrastructure/telegram"
)

type MenuCommand struct {
	Bot      telegram.BotClient
	Registry *command_handling.CommandRegistry
}

func NewMenuCommand(bot telegram.BotClient, reg *command_handling.CommandRegistry) *MenuCommand {
	return &MenuCommand{Bot: bot, Registry: reg}
}

func (c *MenuCommand) Name() string { return "/menu" }

func (c *MenuCommand) Description() string {
	return "Отменить текущую операцию"
}

func (c *MenuCommand) Execute(ctx context.Context, session *domain.UserSession, args []string) error {
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

func (c *MenuCommand) IsStateful() bool {
	return false
}
