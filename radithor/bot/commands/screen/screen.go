package commands_screen

import "radithor/bot/discord"

func (cmd *ScreenCommand) Exec(ctx *discord.CommandContext, args []string) {
	if len(args) < 1 {
		ctx.ReplyMsg("You need to provide a flag.")
		return
	}

	cmd.RunFlag(ctx, args[0], args[1:])
}

func (*ScreenCommand) Name() string {
	return "screen"
}

func (*ScreenCommand) Info() string {
	return "Do things with the desktop display."
}

func CreateCommand() *ScreenCommand {
	cmd := &ScreenCommand{}
	cmd.AddFlag("rotate", "Rotates the display orientation.", rotate)
	cmd.AddFlag("brightness", "Sets the desktop brightness.", brightness)
	cmd.AddFlag("timeout", "Sets the desktop timeout.", timeout)

	return cmd
}

type ScreenCommand struct {
	discord.FlagCommand
}
