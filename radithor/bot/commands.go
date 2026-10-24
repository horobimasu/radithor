package bot

import (
	"radithor/bot/commands"
	"radithor/bot/discord"
)

type FlagCommand interface {
	discord.Command
	GetFlags() []discord.Flag
}

var commandsList = make(map[string]discord.Command)

func RegisterCommands() {
	cmds := []discord.Command{
		&commands.PingCommand{},
	}

	for _, cmd := range cmds {
		commandsList[cmd.Name()] = cmd
	}
}
