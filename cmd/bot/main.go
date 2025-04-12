package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
)

type Command struct {
	Command string `json:"command"`
	URL     string `json:"url"`
}

func sendCommandToScrapper(command string, url string) error {
	scrapperURL := "http://localhost:8080/command"

	cmd := Command{
		Command: command,
		URL:     url,
	}

	jsonData, err := json.Marshal(cmd)
	if err != nil {
		return fmt.Errorf("failed to marshal command: %v", err)
	}

	resp, err := http.Post(scrapperURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to send request to scrapper: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-OK response from scrapper: %s", resp.Status)
	}

	return nil
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	token := os.Getenv("TELEGRAM_TOKEN")
	if token == "" {
		log.Fatal("TELEGRAM_TOKEN is not set in .env file")
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal(err)
	}

	bot.Debug = true

	updates := bot.GetUpdatesChan(tgbotapi.UpdateConfig{
		Timeout: 10,
	})

	fmt.Println("Bot is running...")

	for update := range updates {
		if update.Message.Text == "get" {
			err := sendCommandToScrapper("get", "http://example.com")
			if err != nil {
				log.Printf("Error sending command to scrapper: %v", err)
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Failed to process command"))
				return
			}

			bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Command sent to scrapper successfully"))
		} else if update.Message != nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			bot.Send(msg)
		}
	}
}
