package discord

import (
	"fmt"
	"strings"
)

type Flag struct {
	Exec func(ctx *CommandContext, args []string)
	Name string
	Info string
}

type FlagCommand struct {
	Flags []Flag
}

func (cmd *FlagCommand) AddFlag(
	name string,
	info string,
	exec func(ctx *CommandContext, args []string),
) {
	cmd.Flags = append(cmd.Flags, Flag{
		Name: name,
		Info: info,
		Exec: exec,
	})
}

func (cmd *FlagCommand) GetFlags() []Flag {
	return cmd.Flags
}

func (cmd *FlagCommand) RunFlag(ctx *CommandContext, firstWord string, args []string) bool {
	if !strings.HasPrefix(firstWord, "-") {
		return false
	}

	firstWord = strings.TrimPrefix(firstWord, "-")
	firstWord = strings.ToLower(firstWord)

	if firstWord == "" {
		return false
	}

	for _, flag := range cmd.Flags {
		if strings.ToLower(flag.Name) == firstWord {
			flag.Exec(ctx, args)
			return true
		}
	}

	var matches []Flag

	for _, flag := range cmd.Flags {
		if strings.HasPrefix(strings.ToLower(flag.Name), firstWord) {
			matches = append(matches, flag)
		}
	}

	switch len(matches) {
	case 0:
		return false

	case 1:
		matches[0].Exec(ctx, args)
		return true

	default:
		names := make([]string, len(matches))
		for i, flag := range matches {
			names[i] = "`-" + flag.Name + "`"
		}

		ctx.ReplyMsg(fmt.Sprintf("Flag `%s` does not exist.", firstWord))
		return false
	}
}
