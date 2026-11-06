package commands_internal

import (
	"fmt"

	"radithor/bot/discord"
	"radithor/store"
)

func isDebugMode(ctx *discord.CommandContext, args []string) {
	ctx.ReplyMsg(fmt.Sprintf("In debug mode: `%t`", store.DEBUG_MODE_ENABLED))
}
