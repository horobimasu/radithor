package discord

import (
	"slices"
	"strings"

	"radithor/store"

	"shared/config"

	"github.com/bwmarrin/discordgo"
)

var (
	serverOwner string

	whitelistedUsers []string

	noNeedToCheck bool
)

var FailedToInit bool

func InitWhitelistedUsers(ses *discordgo.Session, serverId string) error {
	allWhitelistedUsers := config.GetConfigString(store.WHITELISTED_OPERATOR_IDS)
	if allWhitelistedUsers == "" {
		noNeedToCheck = true
	}

	var ids []string

	for id := range strings.SplitSeq(allWhitelistedUsers, ";") {
		id = strings.TrimSpace(id)
		if id == "" {
			continue
		}

		if strings.ToLower(id) == "%owner%" {
			server, err := ses.Guild(serverId)
			if err != nil {
				return err
			}

			serverOwner = server.OwnerID
			id = server.OwnerID
		}

		ids = append(ids, id)
	}

	if len(ids) == 0 {
		noNeedToCheck = true
		return nil
	}

	whitelistedUsers = ids
	return nil
}

func CheckCanRunCommand(user *discordgo.User) (bool, string) { // if the user can run the command, the message
	if user.ID != serverOwner && store.DEBUG_MODE_ENABLED && store.DEBUG_MODE_ONLY_LET_OWNER_RUN_COMMANDS {
		return false, "Only the server owner can run commands in debug mode."
	}

	if FailedToInit && user.ID != serverOwner {
		return false, "Only the server ownder can run commands right now as whitelisted user data has failed to load."
	}

	if noNeedToCheck {
		return true, ""
	}

	if slices.Contains(whitelistedUsers, user.ID) {
		return true, ""
	}

	return false, "You are not whitelisted to use commands."
}
