package command_handling

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

	originalName := parts[0]
	var name string
	var args []string

	// Если пользователь в stateful команде, но не хочет сбрасывать
	if session.ActiveCommand != "" && originalName != "/menu" {
		name = session.ActiveCommand
		args = parts
	} else {
		name = originalName
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

	if cmd.IsStateful() {
		session.ActiveCommand = name
	} else {
		session.ActiveCommand = ""
	}

	return cmd.Execute(ctx, session, args)
}
