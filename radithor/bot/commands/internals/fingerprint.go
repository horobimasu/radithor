package commands_internal

import (
	"fmt"
	"strings"

	"radithor/bot/discord"
	"radithor/core/core_funcs"
)

func fingerprint(ctx *discord.CommandContext, args []string) {
	var response strings.Builder

	response.WriteString("```\n")
	for name, value := range core_funcs.GetUnhashedFingerprint() {
		fmt.Fprintf(&response, "%s: %s\n", name, value)
	}

	fmt.Fprintf(&response, "Hashed: %x\n```", core_funcs.GenerateFingerprint())

	ctx.ReplyMsg(response.String())
}
