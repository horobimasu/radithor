package bot

import (
	commands_this "radithor/bot/commands/this"

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
		commands_this.NewThisCommand(),

		&commands.PingCommand{},
		&commands.NotifCommand{},
		&commands.ScreenshotCommand{},
		&commands.NotifExCommand{},
	}

	for _, cmd := range cmds {
		commandsList[cmd.Name()] = cmd
	}
}
