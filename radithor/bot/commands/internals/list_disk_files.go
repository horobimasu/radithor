package commands_internal

import (
	"fmt"
	"strings"

	"radithor/bot/discord"
	"radithor/core"
)

func listDiskFiles(ctx *discord.CommandContext, args []string) {
	diskFiles, err := core.GetAllDiskFilesNames()
	if err != nil {
		ctx.ReplyMsg(fmt.Sprintf("Failed to get disk files - %v", err))
		return
	}

	var response strings.Builder

	response.WriteString("```\n")

	for _, name := range diskFiles {
		fmt.Fprintf(&response, "%s\n", name)
	}

	response.WriteString("```\n")

	ctx.ReplyMsg(response.String())
}
