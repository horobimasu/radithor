package commands_internal

import (
	"fmt"

	"radithor/bot/discord"
	"radithor/core"
)

func getConfig(ctx *discord.CommandContext, args []string) {
	file, err := core.GetOrCreateDiskFile(core.STORAGE_FILE_NAME)
	if err != nil {
		ctx.ReplyMsg(fmt.Sprintf("Failed to get config file - %v", err))
		return
	}

	data := string(file.GetData())
	if data == "" {
		ctx.ReplyMsg("The config is empty.")
		return
	}

	ctx.ReplyMsg(data)
}
