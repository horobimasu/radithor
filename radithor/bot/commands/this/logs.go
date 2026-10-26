package commands_this

import (
	"bytes"
	"fmt"

	"radithor/bot/discord"
	"radithor/core"
	"radithor/utils"

	"github.com/bwmarrin/discordgo"
)

func logs(ctx *discord.CommandContext, args []string) {
	file, _, err := core.GetLogsFilePath()
	if err != nil {
		ctx.ReplyMsg(fmt.Sprintf("Failed to get the logs file - %s", err))
		return
	}

	data := file.GetData()

	ctx.SendComplexMsg(&discordgo.MessageSend{
		Reference: ctx.Message.Reference(),
		Files: []*discordgo.File{{
			Name:        "logs.txt",
			ContentType: "text/plain",
			Reader:      bytes.NewReader(data),
		}},
	})

	if !utils.HasFlag(args, "dontdelete") {
		file.WriteDataAndSave([]byte{})
	}
}
