package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type BotClient interface {
	SendMessage(chatID int64, text string, replyMarkup ...interface{})
}

type TgBotClient struct {
	api *tgbotapi.BotAPI
}

func NewTgBotClient(api *tgbotapi.BotAPI) *TgBotClient {
	return &TgBotClient{api: api}
}

func (c *TgBotClient) SendMessage(chatID int64, text string, replyMarkup ...interface{}) {
	msg := tgbotapi.NewMessage(chatID, text)

	if len(replyMarkup) > 0 {
		switch rm := replyMarkup[0].(type) {
		case tgbotapi.ReplyKeyboardMarkup:
			msg.ReplyMarkup = rm
		case tgbotapi.InlineKeyboardMarkup:
			msg.ReplyMarkup = rm
		case tgbotapi.ReplyKeyboardRemove:
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		}
	}

	c.api.Send(msg)
}
