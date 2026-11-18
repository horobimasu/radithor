package commands_screen

import (
	"fmt"
	"strconv"

	"radithor/bot/discord"
	"radithor/funcs"
)

func timeout(ctx *discord.CommandContext, args []string) {
	if len(args) < 1 {
		ctx.ReplyMsg("You need to provde a number 0 or more as the seconds until the timeout.")
		return
	}

	timeout, err := strconv.Atoi(args[0])
	if err != nil {
		ctx.ReplyMsg(fmt.Sprintf("Failed to parse number - %v", err))
		return
	}

	if timeout < 0 {
		ctx.ReplyMsg("Number needs to be 0 or more.")
		return
	}

	err = funcs.SetScreenTimeout(uint32(timeout))
	if err != nil {
		ctx.ReplyMsg(fmt.Sprintf("Failed to set screen timeout - %v", err))
		return
	}

	ctx.ReplyMsg("Successfully set screen timeout!")
}
