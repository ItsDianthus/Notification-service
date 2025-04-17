package erros

import "fmt"

type ErrNotRegistered struct{}

func (e ErrNotRegistered) Error() string {
	return "registration required"
}

type ErrUnknownCommand struct{ Name string }

func (e ErrUnknownCommand) Error() string {
	return fmt.Sprintf("unknown command: %s", e.Name)
}
