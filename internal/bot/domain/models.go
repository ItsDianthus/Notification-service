package domain

import (
	"time"
)

type TelegramMessage struct {
	MessageID  int64     `json:"message_id"`
	ChatID     int64     `json:"chat_id"`
	FromUserID int64     `json:"from_user_id"`
	Text       string    `json:"text"`
	Timestamp  time.Time `json:"timestamp"`
}

type BotConfig struct {
	Token           string `json:"token"`
	ScrapperBaseURL string `json:"scrapper_base_url"`
}

type TelegramUpdate struct {
	UpdateID int              `json:"update_id"`
	Message  *TelegramMessage `json:"message"`
}

type UpdatesResponse struct {
	Ok     bool             `json:"ok"`
	Result []TelegramUpdate `json:"result"`
}
