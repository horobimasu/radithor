package commands

import (
	"fmt"
	"os"

	"radithor/bot/discord"
	"radithor/utils"
	"radithor/utils/discord_utils"

	"github.com/bwmarrin/discordgo"
)

func (*DownloadCommand) Exec(ctx *discord.CommandContext, args []string) {
	urls := discord_utils.GetUrls(ctx)
	if len(urls) == 0 {
		ctx.ReplyMsg("Failed to find any URLs.")
		return
	}

	var downloadToMsg *discordgo.Message

	dstDir := utils.UnwrapQuotes(args)
	if dstDir == "" {
		dstDir = os.TempDir()
		downloadToMsg, _ = ctx.ReplyMsg(fmt.Sprintf("No directory provided, downloading to: `%s`", dstDir))
	} else {
		info, err := os.Stat(dstDir)
		if err != nil {
			ctx.DeleteMsg(downloadToMsg.ID)
			ctx.ReplyMsg(fmt.Sprintf("Failed to get info on `%s` - %s", dstDir, err))
			return
		}

		if !info.IsDir() {
			ctx.DeleteMsg(downloadToMsg.ID)
			ctx.ReplyMsg(fmt.Sprintf("Path `%s` is not a directory", dstDir))
			return
		}
	}

	downloadMsg, _ := ctx.ReplyMsg(fmt.Sprintf("Downloading %d file(s)...", len(urls)))

	for _, url := range urls {
		progressStream := utils.DownloadFileWithProgress(url, dstDir)

		var err error
		var path string

		fileName := discord_utils.GetFileNameFromUrl(url)

		for progress := range progressStream {
			if progress.Err != nil {
				err = progress.Err
				continue
			}

			if progress.Path != "" {
				path = progress.Path
				continue
			}

			status := fmt.Sprintf("Downloading `%s`...", fileName)

			if progress.Percent >= 0 {
				status = fmt.Sprintf("Downloading `%s` (%d%%)...", fileName, progress.Percent)
			}

			ctx.EditMsg(status, downloadMsg.ID)
		}

		if err != nil {
			ctx.ReplyMsg(fmt.Sprintf("Failed to download `%s` - %s", fileName, err))
			continue
		}

		ctx.ReplyMsg(fmt.Sprintf("Downloaded `%s` to: `%s`", fileName, path))
	}

	ctx.DeleteMsg(downloadToMsg.ID)
	ctx.DeleteMsg(downloadMsg.ID)
}

func (*DownloadCommand) Name() string {
	return "download"
}

func (*DownloadCommand) Info() string {
	return "Download files to the device."
}

type DownloadCommand struct{}
