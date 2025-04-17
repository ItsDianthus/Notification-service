package clients

import (
	"context"
	"go-ItsDianthus-NotificationLink/internal/api/openapi/bot_api"
)

// Ещё будет
type ScrapperClient interface {
	AddLink(ctx context.Context, chatID int64, link string, tags, filters []string) (*bot_api.LinkUpdate, error)
	RemoveLink(ctx context.Context, chatID int64, link string) error
	ListLinks(ctx context.Context, chatID int64) ([]bot_api.LinkUpdate, error)
}
