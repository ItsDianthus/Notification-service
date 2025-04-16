package handler

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go-ItsDianthus-NotificationLink/internal/bot/application"
	"go-ItsDianthus-NotificationLink/internal/bot/domain"
	"go-ItsDianthus-NotificationLink/internal/bot/infrastructure/repo"
	"go-ItsDianthus-NotificationLink/internal/bot/infrastructure/telegram"
	"strings"
)

type MessageProcessor struct {
	BotClient telegram.BotClient
	Repo      *repo.InMemorySessionRepo
	Registry  *application.CommandRegistry
	Updates   tgbotapi.UpdatesChannel
}

func NewMessageProcessor(bot telegram.BotClient, repo *repo.InMemorySessionRepo,
	reg *application.CommandRegistry, updates tgbotapi.UpdatesChannel) *MessageProcessor {
	return &MessageProcessor{BotClient: bot, Repo: repo, Registry: reg, Updates: updates}
}

func (p *MessageProcessor) ProcessUpdates(ctx context.Context) {
	for update := range p.Updates {
		if update.Message == nil {
			continue
		}
		chatID := update.Message.Chat.ID
		session := p.Repo.GetOrCreateSession(chatID)

		if err := application.DispatchCommand(ctx, p.Registry, session, update.Message.Text); err != nil {
			if strings.Contains(err.Error(), "неизвестная команда") && session.CurrentState == domain.StateDefault {
				keyboard := telegram.BuildCommandKeyboard([]string{"/start"})
				p.BotClient.SendMessage(chatID,
					"Нажмите /start для начала",
					keyboard,
				)
			} else {
				p.BotClient.SendMessage(chatID, "Ошибка: "+err.Error())
			}
		}
		p.Repo.SaveSession(session)
	}
}
