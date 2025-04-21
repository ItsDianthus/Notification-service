package commands

import (
	"context"
	"fmt"
	"go-ItsDianthus-NotificationLink/internal/bot/application/command_registry"
	"go-ItsDianthus-NotificationLink/internal/bot/domain"
	"go-ItsDianthus-NotificationLink/internal/bot/infrastructure/telegram"
	"strings"
)

type HelpCommand struct {
	Bot      telegram.BotClient
	Registry *command_registry.CommandRegistry
}

func NewHelpCommand(bot telegram.BotClient, reg *command_registry.CommandRegistry) *HelpCommand {
	return &HelpCommand{Bot: bot, Registry: reg}
}

func (c *HelpCommand) Name() string        { return "/help" }
func (c *HelpCommand) Description() string { return "Показать список команд" }

func (c *HelpCommand) Execute(ctx context.Context, session *domain.UserSession, args []string) error {
	// Собираем описание команд
	names := c.Registry.AllNames()
	var lines []string
	for _, name := range names {
		cmd, ok := c.Registry.Get(name)
		if !ok {
			continue
		}
		if name == "/start" {
			continue
		}
		lines = append(lines, fmt.Sprintf("%s — %s", name, cmd.Description()))
	}
	text := "Доступные команды:\n" + strings.Join(lines, "\n")
	kb := telegram.BuildCommandKeyboard(names)
	c.Bot.SendMessage(session.UserID, text, kb)
	return nil
}
