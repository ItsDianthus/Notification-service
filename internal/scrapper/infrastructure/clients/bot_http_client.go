package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go-ItsDianthus-NotificationLink/internal/api/openapi/bot_api"
	"net/http"
	"time"
)

type BotClient struct {
	baseURL string
	timeout time.Duration
}

func NewBotClient(baseURL string, timeout time.Duration) *BotClient {
	return &BotClient{baseURL: baseURL, timeout: timeout}
}

func (b *BotClient) NotifyUpdate(ctx context.Context, update bot_api.LinkUpdate) error {
	client := &http.Client{Timeout: b.timeout}
	bts, err := json.Marshal(update)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, b.baseURL+"/updates", bytes.NewReader(bts))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("bot update failed: %d", resp.StatusCode)
	}
	return nil
}
