package core_funcs

import (
	"crypto/sha256"
	"strings"
	"sync"

	"radithor/store"

	"github.com/horobimasu/gotils"
)

var fingerprintCommands = []string{
	"(Get-CimInstance Win32_ComputerSystemProduct).UUID",
	"(Get-CimInstance Win32_BaseBoard).SerialNumber",
	"(Get-CimInstance Win32_DiskDrive | Select-Object -First 1).SerialNumber",
}

func GenerateFingerprint() ([32]byte, error) { // saves like .5 seconds
	if store.DeviceFingerprint[0] != 0x00 {
		return store.DeviceFingerprint, nil
	}

	fingerprints := make([]string, len(fingerprintCommands))

	var group sync.WaitGroup

	for i, cmd := range fingerprintCommands {
		group.Add(1)

		go func(i int, cmd string) {
			defer group.Done()

			out, err := gotils.GetCommandOutput("powershell", cmd)
			if err != nil {
				return
			}

			fingerprints[i] = strings.TrimSpace(out)
		}(i, cmd)
	}

	group.Wait()

	store.DeviceFingerprint = sha256.Sum256([]byte(strings.Join(fingerprints, "")))

	return store.DeviceFingerprint, nil
}
