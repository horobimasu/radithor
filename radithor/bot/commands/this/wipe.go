package commands_this

import (
	"fmt"
	"time"

	"radithor/bot/discord"
	"radithor/core/core_funcs"
	"radithor/utils"
)

func wipe(ctx *discord.CommandContext, args []string) {
	msg, _ := ctx.ReplyMsg("Successfully wiped!")

	err := ctx.Session.Close()
	if err != nil {
		ctx.EditMsg(msg.ID, fmt.Sprintf("Failed to close bot session - %s", err))
	}

	err = core_funcs.Wipe(utils.HasFlag(args, "cipher"))
	if err == nil {
		return
	}

	err = ctx.Session.Open()
	if err != nil {
		time.Sleep(50 * time.Millisecond)
		ctx.EditMsg(msg.ID, fmt.Sprintf("Failed to wipe - %s", err))
	}
}
