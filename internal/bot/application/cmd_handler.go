package application

import (
	"context"
	"go-ItsDianthus-NotificationLink/internal/bot/application/errs"
	"go-ItsDianthus-NotificationLink/internal/bot/domain"
	"strings"
)

func HandleCmd(
	ctx context.Context,
	reg *CommandRegistry,
	session *domain.UserSession,
	input string,
) error {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return nil
	}
	name, args := parts[0], parts[1:]

	if !session.IsRegistered && name != "/start" {
		return errs.ErrNotRegistered{}
	}

	if session.ActiveCommand != "" {
		name = session.ActiveCommand
	}

	cmd, ok := reg.Get(name)
	if !ok {
		return errs.ErrUnknownCommand{Name: name}
	}

	session.ActiveCommand = name
	return cmd.Execute(ctx, session, args)
}
