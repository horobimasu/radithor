package commands

import "radithor/bot/discord"

func (*PingCommand) Exec(ctx *discord.CommandContext, args []string) {
	ctx.ReplyMsg("pong")
}

func (*PingCommand) Name() string {
	return "ping"
}

func (*PingCommand) Info() string {
	return "Ping the bot to see if its still alive."
}

type PingCommand struct{}
