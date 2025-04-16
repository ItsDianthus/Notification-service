package telegram

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func BuildCommandKeyboard(cmdNames []string) tgbotapi.ReplyKeyboardMarkup {
	var rows [][]tgbotapi.KeyboardButton
	const perRow = 2
	for i := 0; i < len(cmdNames); i += perRow {
		end := i + perRow
		if end > len(cmdNames) {
			end = len(cmdNames)
		}
		var row []tgbotapi.KeyboardButton
		for _, name := range cmdNames[i:end] {
			row = append(row, tgbotapi.NewKeyboardButton(name))
		}
		rows = append(rows, row)
	}
	return tgbotapi.NewReplyKeyboard(rows...)
}
