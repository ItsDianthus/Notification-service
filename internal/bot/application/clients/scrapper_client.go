package clients

import (
	"context"
)

type ScrapperClient interface {
	AddSubscription(ctx context.Context, chatID int64,
		link string, tags []string, filters map[string]string) error
	RemoveSubscription(ctx context.Context, chatID int64, link string) error
	ListSubscriptions(ctx context.Context, chatID int64) ([]string, error)
}
