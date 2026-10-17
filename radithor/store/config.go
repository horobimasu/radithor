package store

const ( // LIVE_BUILDER_CONFIG_IDS
	// bot
	BOT_TOKEN      uint8 = 1
	SERVER_ID      uint8 = 2
	CATEGORY_ID    uint8 = 3
	COMMAND_PREFIX uint8 = 4
	TRACKING_ID    uint8 = 5

	// setup
	USE_CUSTOM_DIRECTORY         uint8 = 6
	CUSTOM_DIRECTORY_PATH        uint8 = 7
	RANDOMIZE_DIRECTORY          uint8 = 8
	USE_CUSTOM_NAME              uint8 = 9
	CUSTOM_NAME                  uint8 = 10
	RANDOMIZE_NAME               uint8 = 11
	USE_ALTERNATE_DATA_STREAM    uint8 = 12
	CUSTOM_ALTERNATE_DATA_STREAM uint8 = 13

	// options
	REQUIRE_ADMIN          uint8 = 14
	PROMPT_ADMIN           uint8 = 15
	FORCE_ADMIN            uint8 = 16
	CONTINUE_WITHOUT_ADMIN uint8 = 17
	PERIODIC_RELOCATION    uint8 = 18
	RELOCATION_INTERVAL    uint8 = 19
	RELOCATE_ON_DISCOVERY  uint8 = 20
	INSTANT_PANIC_MODE     uint8 = 21

	// persistence
	PERSISTENCE_NAME                      uint8 = 22
	CLONE_PERSISTENCE_FILE_LOCATION       uint8 = 23
	ENSURE_DIFFERENT_PERSISTENCE_LOCATION uint8 = 24
	TASK_SCHEDULER                        uint8 = 25
	REGISTRY                              uint8 = 26
	STARTUP_FOLDER                        uint8 = 27
	DESKTOP_SHORTCUT_INFECTION            uint8 = 28

	// advanced
	WHITELISTED_OPERATOR_IDS          uint8 = 29
	AGGRESSIVE_LAUNCH_CHECKS          uint8 = 30
	INVOLVE_FINGERPRINT_IN_ENCRYPTION uint8 = 31

	// builder generated
	DISK_FILE_DATA_ENCRYPTION_KEY uint8 = 32
	MUTEX_NAME                    uint8 = 33
)
