package commands_screen

import (
	"fmt"
	"strconv"

	"radithor/bot/discord"
	"radithor/funcs"
)

func brightness(ctx *discord.CommandContext, args []string) {
	if len(args) < 1 {
		ctx.ReplyMsg("You need to provde a number from 0 to 100 as the brightness level.")
		return
	}

	level, err := strconv.Atoi(args[0])
	if err != nil {
		ctx.ReplyMsg(fmt.Sprintf("Failed to parse number - %v", err))
		return
	}

	if level < 0 || level > 100 {
		ctx.ReplyMsg("Number needs to be from 0 to 100.")
		return
	}

	err = funcs.SetScreenBrightness(uint32(level))
	if err != nil {
		ctx.ReplyMsg(fmt.Sprintf("Failed to set screen brightness - %v", err))
		return
	}

	ctx.ReplyMsg("Successfully set brightness!")
}
