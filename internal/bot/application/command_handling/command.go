package command_handling

import (
	"context"
	"go-ItsDianthus-NotificationLink/internal/bot/domain"
)

type Command interface {
	Name() string

	Description() string

	Execute(ctx context.Context, session *domain.UserSession, args []string) error

	IsStateful() bool
}
