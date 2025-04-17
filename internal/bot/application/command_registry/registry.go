package command_registry

type CommandRegistry struct {
	commands map[string]Command
}

func NewCommandRegistry() *CommandRegistry {
	return &CommandRegistry{commands: make(map[string]Command)}
}

func (r *CommandRegistry) Register(cmd Command) {
	r.commands[cmd.Name()] = cmd
}

func (r *CommandRegistry) Get(name string) (Command, bool) {
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

func (r *CommandRegistry) AllExceptStart() []string {
	all := r.AllNames()
	filtered := make([]string, 0, len(all)-1)
	for _, name := range all {
		if name == "/start" {
			continue
		}
		filtered = append(filtered, name)
	}
	return filtered
}
