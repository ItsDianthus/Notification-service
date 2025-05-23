package application

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go-ItsDianthus-NotificationLink/internal/bot/application/command_handling"
	"go-ItsDianthus-NotificationLink/internal/bot/application/erros"
	"go-ItsDianthus-NotificationLink/internal/bot/infrastructure/repo"
	"go-ItsDianthus-NotificationLink/internal/bot/infrastructure/telegram"
	"log/slog"
)

type MessageProcessor struct {
	BotClient telegram.BotClient
	Repo      *repo.InMemorySessionRepo
	Registry  *command_handling.CommandRegistry
	Updates   tgbotapi.UpdatesChannel
	Logger    *slog.Logger
}

func NewMessageProcessor(
	bot telegram.BotClient,
	repo *repo.InMemorySessionRepo,
	reg *command_handling.CommandRegistry,
	updates tgbotapi.UpdatesChannel,
	logger *slog.Logger,
) *MessageProcessor {
	return &MessageProcessor{
		BotClient: bot,
		Repo:      repo,
		Registry:  reg,
		Updates:   updates,
		Logger:    logger,
	}
}

func (p *MessageProcessor) ProcessUpdates(ctx context.Context) {
	for update := range p.Updates {
		if update.Message == nil {
			continue
		}
		chatID := update.Message.Chat.ID
		text := update.Message.Text

		p.Logger.Info("Incoming message",
			slog.Int64("chat_id", update.Message.Chat.ID),
			slog.String("text", update.Message.Text),
		)

		session := p.Repo.GetOrCreate(chatID)
		err := command_handling.HandleCmd(ctx, p.Registry, session, text)

		switch e := err.(type) {
		case nil:
			p.Logger.Info("Command executed",
				slog.Int64("chat_id", chatID),
				slog.String("input", text),
			)

		case erros.ErrNotRegistered:
			p.Logger.Warn("User not registered",
				slog.Int64("chat_id", chatID),
				slog.String("input", text),
			)
			p.BotClient.SendMessage(chatID,
				"Пожалуйста, зарегистрируйтесь кнопкой /start",
				telegram.BuildCommandKeyboard([]string{"/start"}),
			)

		case erros.ErrUnknownCommand:
			p.Logger.Warn("Unknown command",
				slog.Int64("chat_id", chatID),
				slog.String("input", text),
			)
			all := p.Registry.AllNames()
			kb := telegram.BuildCommandKeyboard(all)
			p.BotClient.SendMessage(chatID,
				"Такой команды нет. Нажмите /help, чтобы увидеть список команд",
				kb,
			)

		default:
			p.Logger.Error("Command execution failed",
				slog.Int64("chat_id", chatID),
				slog.String("input", text),
				slog.String("error", e.Error()),
			)
			p.BotClient.SendMessage(chatID, "Ошибка: "+e.Error(), nil)
		}

		p.Repo.Save(session)
	}
}
