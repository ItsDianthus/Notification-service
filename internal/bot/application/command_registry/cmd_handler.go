package command_registry

import (
	"context"
	"go-ItsDianthus-NotificationLink/internal/bot/application/erros"
	"strings"

	"go-ItsDianthus-NotificationLink/internal/bot/domain"
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

	name := parts[0]
	var args []string
	if session.ActiveCommand != "" && name != "/menu" {
		name = session.ActiveCommand
		args = parts
	} else {
		args = parts[1:]
	}

	if name == "/menu" {
		session.CurrentState = domain.StateDefault
		session.ActiveCommand = ""
		session.TempData = nil

		menuCmd, _ := reg.Get("/menu") // /menu точно зарегистрирована
		return menuCmd.Execute(ctx, session, nil)
	}

	if !session.IsRegistered && name != "/start" {
		return erros.ErrNotRegistered{}
	}

	cmd, ok := reg.Get(name)
	if !ok {
		return erros.ErrUnknownCommand{Name: name}
	}

	session.ActiveCommand = name
	return cmd.Execute(ctx, session, args)
}
