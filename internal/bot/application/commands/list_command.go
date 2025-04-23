package commands

import (
	"context"
	"go-ItsDianthus-NotificationLink/internal/bot/application/clients"
	"go-ItsDianthus-NotificationLink/internal/bot/domain"
	"go-ItsDianthus-NotificationLink/internal/bot/infrastructure/telegram"
	"strings"
)

type ListCommand struct {
	Bot            telegram.BotClient
	ScrapperClient clients.ScrapperClient
}

func NewListCommand(bot telegram.BotClient, scr clients.ScrapperClient) *ListCommand {
	return &ListCommand{Bot: bot, ScrapperClient: scr}
}

func (c *ListCommand) Name() string        { return "/list" }
func (c *ListCommand) Description() string { return "Показать текущие подписки" }

func (c *ListCommand) Execute(ctx context.Context, session *domain.UserSession, args []string) error {
	subs, err := c.ScrapperClient.ListSubscriptions(ctx, session.UserID)
	if err != nil {
		c.Bot.SendMessage(session.UserID, "Ошибка получения списка: "+err.Error(), nil)
		return err
	}
	if len(subs) == 0 {
		c.Bot.SendMessage(session.UserID,
			"У вас нет подписок.",
			telegram.BuildCommandKeyboard([]string{"/track", "/list"}),
		)
		return nil
	}
	var builder strings.Builder
	builder.WriteString("Ваши подписки:\n")
	for _, link := range subs {
		builder.WriteString("• " + link + "\n")
	}
	kb := telegram.BuildCommandKeyboard([]string{"/track", "/untrack", "/list"})
	c.Bot.SendMessage(session.UserID, builder.String(), kb)
	return nil
}

func (c *ListCommand) IsStateful() bool {
	return false
}
