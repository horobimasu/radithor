package commands

import (
	"fmt"
	"strconv"

	"radithor/bot/discord"
	"radithor/funcs"

	"github.com/bwmarrin/discordgo"
)

func (*ScreenshotCommand) Exec(ctx *discord.CommandContext, args []string) {
	var screen int

	if len(args) == 1 {
		num, err := strconv.Atoi(args[0])
		if err == nil {
			screen = num
		}
	}

	pic, err := funcs.CaptureScreen(screen)
	if err != nil {
		ctx.ReplyMsg(fmt.Sprintf("Failed to take screenshot - %v", err))
		return
	}

	ctx.SendComplexMsg(&discordgo.MessageSend{
		Reference: ctx.Message.Reference(),
		Files: []*discordgo.File{{
			Name:   "ss.jpg",
			Reader: pic,
		}},
	})
}

func (*ScreenshotCommand) Name() string {
	return "ss"
}

func (*ScreenshotCommand) Info() string {
	return "Takes a screenshot of the screen."
}

type ScreenshotCommand struct{}
