package application

import "go-ItsDianthus-NotificationLink/internal/bot/application/commands"

type CommandRegistry struct {
	commands map[string]commands.Command
}

func NewCommandRegistry() *CommandRegistry {
	return &CommandRegistry{commands: make(map[string]commands.Command)}
}

func (r *CommandRegistry) Register(cmd commands.Command) {
	r.commands[cmd.Name()] = cmd
}

func (r *CommandRegistry) Get(name string) (commands.Command, bool) {
	cmd, ok := r.commands[name]
	return cmd, ok
}

func (r *CommandRegistry) AllNames() []string {
	names := make([]string, 0, len(r.commands))
	for name := range r.commands {
		names = append(names, name)
	}
	return names
}
