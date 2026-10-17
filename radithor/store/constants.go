package store

const (
	// if this is false, none of the other
	// vals in this block will matter
	DEBUG_MODE_ENABLED = true

	// overrides the value of state.DeviceFingerprint to "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"
	// and therfore overrides the results of core.GenerateFingerprint()
	DEBUG_MODE_HIDE_FINGERPRINT = true

	// the logs will also be written to this file after it is written to the main
	// logs file
	DEBUG_MODE_LOG_FILE_PROXY_TO_PATH = "~\\Desktop\\rad_debug_logs.txt"

	// the name of the mutex file when in debug mode
	// if it is empty, it will defalt to the default one in the config
	DEBUG_MODE_MUTEX_FILE_NAME = "radithor mutex"
)

const (
	DISALLOW_LOGGING = false

	SHORTEN_USERPROFILE_PATHS_IN_LOGS = true
)

const (
	VERSION_UPDATE = 1
	VERSION_PATCH  = 0
)
