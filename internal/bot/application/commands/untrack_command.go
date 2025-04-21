package commands

import (
	"context"
	"fmt"
	"go-ItsDianthus-NotificationLink/internal/bot/application/clients"
	"go-ItsDianthus-NotificationLink/internal/bot/domain"
	"go-ItsDianthus-NotificationLink/internal/bot/infrastructure/telegram"
)

type UntrackCommand struct {
	Bot            telegram.BotClient
	ScrapperClient clients.ScrapperClient
}

func NewUntrackCommand(bot telegram.BotClient, scr clients.ScrapperClient) *UntrackCommand {
	return &UntrackCommand{Bot: bot, ScrapperClient: scr}
}

func (c *UntrackCommand) Name() string        { return "/untrack" }
func (c *UntrackCommand) Description() string { return "Отписаться от ссылки" }

func (c *UntrackCommand) Execute(ctx context.Context, session *domain.UserSession, args []string) error {
	// Ожидаем URL сразу
	if len(args) == 0 {
		c.Bot.SendMessage(session.UserID,
			"Укажите ссылку: /untrack <url>",
			telegram.BuildCommandKeyboard([]string{"/track", "/untrack", "/list"}),
		)
		return nil
	}
	link := args[0]
	err := c.ScrapperClient.RemoveSubscription(ctx, session.UserID, link)
	fullKb := telegram.BuildCommandKeyboard([]string{"/track", "/untrack", "/list"})
	if err != nil {
		c.Bot.SendMessage(session.UserID, "Ошибка при удалении: "+err.Error(), nil)
		return err
	}
	c.Bot.SendMessage(session.UserID,
		fmt.Sprintf("Вы отписались от %s", link),
		fullKb,
	)
	return nil
}
