package lockdown

import (
	"runtime"
	"syscall"

	"radithor/core"

	"github.com/horobimasu/windowsx"
	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

const _SYSTEM_POLICIES = "SOFTWARE\\Policies\\Microsoft\\Windows\\System"

const _CONTROL_CENTER_WINDOW_NAME = "ControlCenterWindow"

var (
	wasOpen bool

	systemWindowHookThread uint32

	onSystemWindowEvent = syscall.NewCallback(func(
		hook uintptr,
		event uintptr,
		handle uintptr,
		object uintptr,
		child uintptr,
		thread uintptr,
		time uintptr,
	) uintptr {
		if int32(object) != 0 {
			return 0
		}

		name, _ := windowsx.GetClassName(handle)

		open := name == _CONTROL_CENTER_WINDOW_NAME
		if open && !wasOpen {
			panicModeCallback()
		}

		wasOpen = open

		return 0
	})
)

func RemoveLockScreenNetworkUi() error {
	key, _, err := registry.CreateKey(registry.LOCAL_MACHINE, _SYSTEM_POLICIES, registry.SET_VALUE)
	if err != nil {
		return err
	}

	defer key.Close()

	return key.SetDWordValue("DontDisplayNetworkSelectionUI", 1)
}

func RestoreLockScreenNetworkUi() error {
	key, _, err := registry.CreateKey(registry.LOCAL_MACHINE, _SYSTEM_POLICIES, registry.SET_VALUE)
	if err != nil {
		return err
	}

	defer key.Close()

	return key.DeleteValue("DontDisplayNetworkSelectionUI")
}

func EnableControlCenterDetection() {
	wasOpen = false

	go func() {
		runtime.LockOSThread()

		var msg windowsx.Msg
		windowsx.PeekMessage(&msg, 0, _WM_USER, _WM_USER, windowsx.PM_NOREMOVE)

		systemWindowHookThread = windows.GetCurrentThreadId()

		hook, err := windowsx.SetEventHook(
			windowsx.EVENT_SYSTEM_FOREGROUND,
			windowsx.EVENT_SYSTEM_FOREGROUND,
			0,
			onSystemWindowEvent,
			0,
			0,
			windowsx.WINEVENT_OUTOFCONTEXT,
		)

		if err != nil {
			core.LogError("[EnableControlCenterDetection] failed to set hook")
			return
		}

		for {
			ret, _ := windowsx.GetMessage(&msg, 0, 0, 0)
			if ret <= 0 {
				break
			}
		}

		windowsx.UnhookEvent(hook)
	}()
}

func DisableControlCenterDetection() {
	windowsx.PostThreadMessage(systemWindowHookThread, _WM_QUIT, 0, 0)
}
