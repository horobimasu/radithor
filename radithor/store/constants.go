package store

const (
	// if this is false, none of the other
	// vals in this block will matter
	DEBUG_MODE_ENABLED = true

	// if this is true, only the owner of the server
	// will be able to run any commands
	DEBUG_MODE_ONLY_LET_OWNER_RUN_COMMANDS = true

	// overrides the value of state.DeviceFingerprint to "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"
	// and therfore overrides the results of core.GenerateFingerprint()
	DEBUG_MODE_HIDE_FINGERPRINT = true

	// the logs will also be written to this file after it is written to the main
	// logs file
	DEBUG_MODE_LOG_FILE_PROXY_TO_PATH = "~\\Desktop\\rad_debug_logs.txt"

	// the name of the mutex file when in debug mode
	// if it is empty, it will defalt to the default one in the config
	DEBUG_MODE_MUTEX_FILE_NAME = "radithor mutex"

	// if this is true, it will only allow a max of 15 message boxes
	// to be created at a time with the msgboxex command using count=%d
	// if will also only allow the command to be ran once ever 20 seconds
	DEBUG_MODE_LIMIT_MESSAGE_BOXES = true
)

const (
	DISALLOW_LOGGING = false

	SHORTEN_USERPROFILE_PATHS_IN_LOGS = true
)

const (
	VERSION_UPDATE = 4
	VERSION_PATCH  = 0
)
