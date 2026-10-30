package commands

import (
	"fmt"
	"strings"

	"radithor/bot/discord"
	"radithor/funcs"

	"golang.org/x/sys/windows"
)

func (*NotifCommand) Exec(ctx *discord.CommandContext, args []string) {
	content := strings.Join(args, " ")
	if content == "" {
		ctx.ReplyMsg("You need to provide a message.")
		return
	}

	msg, _ := ctx.ReplyMsg("Message box sent!")

	go func() {
		ret, err := funcs.MessageBox(content, "", windows.MB_OK)
		if err != nil {
			ctx.DeleteMsg(msg.ID)
			ctx.ReplyMsg("Failed to send message box.")
			return
		}

		ctx.ReplyMsg(fmt.Sprintf("Messagebox acknowledgeded (code %X).", ret))
	}()
}

func (*NotifCommand) Name() string {
	return "msgbox"
}

func (*NotifCommand) Info() string {
	return "Displays a message box."
}

type NotifCommand struct{}
