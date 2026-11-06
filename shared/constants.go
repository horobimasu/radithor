package shared

const ENCRYPTION_KEY_LEN = 32

const (
	RADITHOR_CONFIG_PE_SECTION_NAME = ".rodata"

	RADITHOR_ENCRYPTION_KEYS_LEN = ENCRYPTION_KEY_LEN + ENCRYPTION_KEY_LEN
)

var (
	// generated with scripts/build-live
	RADITHOR_RAW_ENCRYPTION_KEYS_FINGERPRINT = []byte{ // LIVE_BUILDER_MATCH:ENCRYPTION_KEYS_FINGERPRINT
		0x0E, 0x1B, 0x22, 0xB6, // Fingerprint
		0x20, 0xF4, 0x2F, 0xC0, // Fingerprint
	}
)
