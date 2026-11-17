package shared

const (
	ENCRYPTION_KEY_LEN = 32
	FINGERPRINT_LEN    = 8
)

const (
	RADITHOR_CONFIG_PE_SECTION_NAME = ".rodata"

	RADITHOR_ENCRYPTION_KEYS_LEN = (ENCRYPTION_KEY_LEN * 2) // the len of the EncryptionKeys struct without the fingerprint
)

// generated with scripts/build-live
var (
	RADITHOR_ENCRYPTION_KEYS_FINGERPRINT = []byte{ // LIVE_BUILDER_MATCH:ENCRYPTION_KEYS_FINGERPRINT
		0x3C, 0x41, 0x22, 0xB5, // Fingerprint
		0x77, 0xE3, 0x45, 0xC1, // Fingerprint
	}
)
