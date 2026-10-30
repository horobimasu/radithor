package core

import (
	"os"

	"radithor/store"

	"golang.org/x/sys/windows"
)

func ExitFatalOnDebugMode(reason string) {
	if !store.DEBUG_MODE_ENABLED {
		return
	}

	windows.MessageBox(
		windows.HWND(0),
		windows.StringToUTF16Ptr(reason),
		windows.StringToUTF16Ptr("radithor debug mode [fatal error]"),
		windows.MB_OK|windows.MB_ICONERROR,
	)

	os.Exit(1)
}
