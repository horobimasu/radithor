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
)

// the things that cannot be set with state.InitState()
var (
	// technically this can be set but only in debug mode tho
	// so it gets a pass
	DeviceFingerprint [32]byte

	// set in /main.go
	DiskFileDataEncryptionKey []byte
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
}
