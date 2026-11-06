package store

import (
	"os"
	"path/filepath"

	"golang.org/x/sys/windows"
)

var (
	IsElevated bool

	ExecPath string

	ExecDir string

	InPanicMode bool

	DiskFileDataEncryptionKeySet bool
)

// the things that cannot be set with state.InitState()
var (
	// technically this can be set but only in debug mode tho
	// so it gets a pass
	//
	// set in /core/core_funcs.GenerateFingerprint
	DeviceFingerprint [32]byte

	// set in /main.initConfig
	DiskFileDataEncryptionKey [32]byte
)

func InitState() {
	IsElevated = windows.GetCurrentProcessToken().IsElevated()

	ExecPath, _ = os.Executable()

	ExecDir = filepath.Dir(ExecPath)

	if DEBUG_MODE_ENABLED && DEBUG_MODE_HIDE_FINGERPRINT {
		for i := range DeviceFingerprint {
			DeviceFingerprint[i] = 0x41
		}
	}

	InPanicMode = false

	DiskFileDataEncryptionKeySet = false
}
