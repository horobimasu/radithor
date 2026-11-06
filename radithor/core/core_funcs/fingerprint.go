package core_funcs

import (
	"crypto/sha256"
	"strings"
	"sync"

	"radithor/store"

	"github.com/horobimasu/gotils"
)

type fingerprintCommand struct {
	Name    string
	Command string
}

var fingerprintCommands = []fingerprintCommand{
	{"Computer product UUID", "(gcim Win32_ComputerSystemProduct).UUID"},
	{"Base board serial number", "(gcim Win32_BaseBoard).SerialNumber"},
	{"Disk drive serial number", "(gcim Win32_DiskDrive | Select-Object -First 1).SerialNumber"},
}

func GenerateFingerprint() [32]byte { // saves like .5 seconds
	if store.DeviceFingerprint[0] != 0x00 {
		return store.DeviceFingerprint
	}

	fingerprints := make([]string, len(fingerprintCommands))

	var group sync.WaitGroup

	for i, cmd := range fingerprintCommands {
		group.Add(1)

		go func(i int, cmd string) {
			defer group.Done()

			out, err := gotils.GetProgramOutput("powershell", "-nop", "-ep", "bypass", cmd)
			if err != nil {
				return
			}

			fingerprints[i] = strings.TrimSpace(out)
		}(i, cmd.Command)
	}

	group.Wait()

	store.DeviceFingerprint = sha256.Sum256([]byte(
		strings.Join(fingerprints, ""),
	))

	return store.DeviceFingerprint
}

func GetUnhashedFingerprint() map[string]string {
	results := make([]string, len(fingerprintCommands))

	var group sync.WaitGroup

	for i, cmd := range fingerprintCommands {
		group.Add(1)

		go func(i int, cmd string) {
			defer group.Done()

			out, err := gotils.GetProgramOutput("powershell", "-nop", "-ep", "bypass", cmd)
			if err != nil {
				return
			}

			results[i] = strings.TrimSpace(out)
		}(i, cmd.Command)
	}

	group.Wait()

	fingerprints := make(map[string]string, len(fingerprintCommands))
	for i, fingerprintThing := range fingerprintCommands {
		var result string

		if store.DEBUG_MODE_ENABLED && store.DEBUG_MODE_HIDE_FINGERPRINT {
			result = strings.Repeat("0", 8)
		} else {
			result = results[i]
		}

		fingerprints[fingerprintThing.Name] = result
	}

	return fingerprints
}
