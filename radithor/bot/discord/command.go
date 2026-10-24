package discord

type Command interface {
	Exec(ctx *CommandContext, args []string)
	Name() string
	Info() string
}
