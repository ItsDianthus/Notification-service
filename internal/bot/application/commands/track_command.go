package commands

import (
	"context"
	"fmt"
	appClients "go-ItsDianthus-NotificationLink/internal/bot/application/clients"
	"go-ItsDianthus-NotificationLink/internal/bot/application/command_registry"
	"net/url"
	"strings"

	"go-ItsDianthus-NotificationLink/internal/bot/domain"
	"go-ItsDianthus-NotificationLink/internal/bot/infrastructure/telegram"
)

type TrackCommand struct {
	Bot            telegram.BotClient
	ScrapperClient appClients.ScrapperClient
	Registry       *command_registry.CommandRegistry
}

func NewTrackCommand(bot telegram.BotClient, scr appClients.ScrapperClient, reg *command_registry.CommandRegistry) *TrackCommand {
	return &TrackCommand{Bot: bot, ScrapperClient: scr, Registry: reg}
}

func (c *TrackCommand) Name() string { return "/track" }
func (c *TrackCommand) Description() string {
	return "Добавить новую ссылку для отслеживания"
}

func (c *TrackCommand) Execute(ctx context.Context, session *domain.UserSession, args []string) error {
	if session.TempData == nil {
		session.TempData = make(map[string]interface{})
	}

	switch session.CurrentState {
	case domain.StateDefault:
		session.CurrentState = domain.StateAwaitingLink
		session.ActiveCommand = c.Name()
		c.Bot.SendMessage(session.UserID,
			"Пожалуйста, отправьте ссылку для отслеживания (URL):",
			telegram.BuildCommandKeyboard([]string{"/menu"}),
		)
		return nil

	case domain.StateAwaitingLink:
		if len(args) == 0 {
			c.Bot.SendMessage(session.UserID,
				"Некорректный формат. Отправьте ссылку:",
				telegram.BuildCommandKeyboard([]string{"/menu"}),
			)
			return nil
		}
		link := args[0]
		parsed, err := url.Parse(link)
		if err != nil || (parsed.Scheme != "https" || !(strings.HasPrefix(parsed.Host, "github.com") || strings.HasPrefix(parsed.Host, "stackoverflow.com"))) {
			c.Bot.SendMessage(session.UserID,
				"Неподдерживаемый URL. Отправьте ссылку вида https://github.com/... или https://stackoverflow.com/....",
				telegram.BuildCommandKeyboard([]string{"/menu"}),
			)
			return nil
		}
		session.TempData["link"] = link

		session.CurrentState = domain.StateAwaitingTags
		c.Bot.SendMessage(session.UserID,
			"Введите тэги (через пробел) или нажмите /skip, чтобы пропустить:",
			telegram.BuildCommandKeyboard([]string{"/menu", "/skip"}),
		)
		return nil

	case domain.StateAwaitingTags:
		var tags []string
		if len(args) == 1 && args[0] == "/skip" {
			tags = nil
		} else if len(args) > 0 {
			tags = args
		}
		session.TempData["tags"] = tags

		session.CurrentState = domain.StateAwaitingFilters
		c.Bot.SendMessage(session.UserID,
			"Настройте фильтры в формате key:value через пробел или нажмите /skip, чтобы пропустить:",
			telegram.BuildCommandKeyboard([]string{"/menu", "/skip"}),
		)
		return nil

	case domain.StateAwaitingFilters:
		var filters map[string]string
		if len(args) == 1 && args[0] == "/skip" {
			filters = nil
		} else {
			filters = make(map[string]string)
			for _, tok := range args {
				parts := strings.SplitN(tok, ":", 2)
				if len(parts) == 2 {
					filters[parts[0]] = parts[1]
				}
			}
		}
		session.TempData["filters"] = filters

		link, _ := session.TempData["link"].(string)
		tags, _ := session.TempData["tags"].([]string)

		if err := c.ScrapperClient.AddSubscription(ctx, session.UserID, link, tags, filters); err != nil {
			return fmt.Errorf("Не удалось добавить подписку: %w", err)
		}

		session.CurrentState = domain.StateDefault
		session.ActiveCommand = ""
		session.TempData = nil

		cmds := c.Registry.AllExceptStart()
		kb := telegram.BuildCommandKeyboard(cmds)
		c.Bot.SendMessage(session.UserID,
			fmt.Sprintf("Подписка на %s успешно добавлена!", link),
			kb,
		)
		return nil

	default:
		session.CurrentState = domain.StateDefault
		session.ActiveCommand = ""
		session.TempData = nil
		c.Bot.SendMessage(session.UserID,
			"Произошла ошибка. Попробуйте ещё раз командой /track",
			telegram.BuildCommandKeyboard([]string{"/track"}),
		)
		return nil
	}
}
