package commands_internal

import (
	"radithor/bot/discord"
	"radithor/store"
)

func (cmd *InternalCommand) Exec(ctx *discord.CommandContext, args []string) {
	if len(args) < 1 {
		ctx.ReplyMsg("You need to provide a flag.")
		return
	}

	cmd.RunFlag(ctx, args[0], args[1:])
}

func (*InternalCommand) Name() string {
	return "internal"
}

func (*InternalCommand) Info() string {
	return "Do things with Radithor's internals."
}

// todo, test autocomplreate flags thing
func CreateCommand() *InternalCommand {
	cmd := &InternalCommand{}
	cmd.AddFlag("listDiskFiles", "List all disk files.", listDiskFiles)
	cmd.AddFlag("getConfig", "Gets the persistent config.", getConfig)
	cmd.AddFlag("isDebugMode", "Check the status of debug mode.", isDebugMode)
	cmd.AddFlag("fingerprint", "Gets the unhashed and hashed device fingerprint.", fingerprint)

	if store.DEBUG_MODE_ENABLED {
		cmd.AddFlag("getState", "Gets Radithor's state.", getState)
	}

	return cmd
}

type InternalCommand struct {
	discord.FlagCommand
}
