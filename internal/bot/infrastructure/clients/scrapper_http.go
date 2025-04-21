// internal/bot/infrastructure/clients/scrapper_http.go
package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go-ItsDianthus-NotificationLink/internal/api/openapi/scrapper_api"
	"io"
	"net/http"
	"strconv"
	"time"
)

type ScrapperHTTPClient struct {
	baseURL string
	client  *http.Client
}

func NewScrapperHTTPClient(baseURL string, timeout time.Duration) *ScrapperHTTPClient {
	return &ScrapperHTTPClient{
		baseURL: baseURL,
		client:  &http.Client{Timeout: timeout},
	}
}

func (c *ScrapperHTTPClient) RegisterChat(ctx context.Context, chatID int64) error {
	url := fmt.Sprintf("%s/tg-chat/%d", c.baseURL, chatID)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("unexpected status %d", resp.StatusCode)
	}
	return nil
}

func (c *ScrapperHTTPClient) UnregisterChat(ctx context.Context, chatID int64) error {
	url := fmt.Sprintf("%s/tg-chat/%d", c.baseURL, chatID)
	req, _ := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected status %d", resp.StatusCode)
	}
	return nil
}

func (c *ScrapperHTTPClient) AddSubscription(
	ctx context.Context,
	chatID int64,
	link string,
	tags []string,
	filters map[string]string,
) error {
	reqBody := scrapper_api.AddLinkRequest{
		Link: &link,
		Tags: &tags,
	}
	if len(filters) > 0 {
		fs := make([]string, 0, len(filters))
		for k, v := range filters {
			fs = append(fs, fmt.Sprintf("%s:%s", k, v))
		}
		reqBody.Filters = &fs
	}

	b, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("marshal AddLinkRequest: %w", err)
	}

	url := c.baseURL + "/links"
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Tg-Chat-Id", strconv.FormatInt(chatID, 10))

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	body, _ := io.ReadAll(resp.Body)
	var apiErr scrapper_api.ApiErrorResponse
	if err := json.Unmarshal(body, &apiErr); err == nil && apiErr.Description != nil {
		return fmt.Errorf("scrapper error: %s", *apiErr.Description)
	}

	return fmt.Errorf("unexpected status %d from scrapper", resp.StatusCode)
}

func (c *ScrapperHTTPClient) RemoveSubscription(
	ctx context.Context,
	chatID int64,
	link string,
) error {
	url := c.baseURL + "/links"
	reqBody := scrapper_api.RemoveLinkRequest{Link: &link}
	b, _ := json.Marshal(reqBody)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Tg-Chat-Id", strconv.FormatInt(chatID, 10))

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}
	return fmt.Errorf("unexpected status %d from scrapper", resp.StatusCode)
}

func (c *ScrapperHTTPClient) ListSubscriptions(
	ctx context.Context,
	chatID int64,
) ([]string, error) {
	url := c.baseURL + "/links"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Tg-Chat-Id", strconv.FormatInt(chatID, 10))

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status %d from scrapper", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	var listResp scrapper_api.ListLinksResponse
	if err := json.Unmarshal(body, &listResp); err != nil {
		return nil, err
	}

	urls := make([]string, 0)
	if listResp.Links != nil {
		for _, item := range *listResp.Links {
			if item.Url != nil {
				urls = append(urls, *item.Url)
			}
		}
	}
	return urls, nil
}
