package commands

import (
	"bytes"
	"fmt"
	"strings"

	"radithor/bot/discord"
	"radithor/utils"

	"github.com/horobimasu/gotils"
)

func (*RunCommand) Exec(ctx *discord.CommandContext, args []string) {
	if len(args) < 1 {
		ctx.ReplyMsg("You need to provide a program to run.")
		return
	}

	program := args[0]
	arguments := utils.ParseArgs(args[1:])

	switch strings.ToLower(program) {
	case "cmd":
		if len(arguments) == 0 {
			ctx.ReplyMsg("You need to provide 1 or more arguments to execute this program.")
			return
		}

		arguments = append([]string{"/c"}, arguments...)

	case "powershell", "ps", "ps1":
		if len(arguments) == 0 {
			ctx.ReplyMsg("You need to provide 1 or more arguments to execute this program.")
			return
		}

		program = "powershell"
		arguments = append([]string{"-c"}, arguments...)
	}

	cmd := gotils.StartProgram(program, arguments...)

	var output bytes.Buffer
	cmd.Stdout = &output

	var error bytes.Buffer
	cmd.Stderr = &error

	err := cmd.Run()
	if err != nil {
		ctx.ReplyMsg(fmt.Sprintf("Failed to run - %v", err))
		return
	}

	var response strings.Builder

	if output.Len() > 0 {
		fmt.Fprintf(&response, "**Output:**\n```\n%s```\n", output.String())
	}

	if error.Len() > 0 {
		fmt.Fprintf(&response, "**Error:**\n```\n%s```", error.String())
	}

	if response.Len() == 0 {
		response.WriteString("No output.")
	}

	ctx.ReplyMsg(response.String())
}

func (*RunCommand) Name() string {
	return "run"
}

func (*RunCommand) Info() string {
	return "Run programs with arguments."
}

type RunCommand struct{}
