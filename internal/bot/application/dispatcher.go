package application

import (
	"context"
	"fmt"
	"go-ItsDianthus-NotificationLink/internal/bot/domain"
	"strings"
)

func DispatchCommand(ctx context.Context, reg *CommandRegistry, session *domain.UserSession, input string) error {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return nil
	}
	name, args := parts[0], parts[1:]

	if session.ActiveCommand != "" {
		name = session.ActiveCommand
	}

	cmd, ok := reg.Get(name)
	if !ok {
		return fmt.Errorf("неизвестная команда: %s", name)
	}
	session.ActiveCommand = name
	return cmd.Execute(ctx, session, args)
}
