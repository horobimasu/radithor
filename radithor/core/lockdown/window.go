package lockdown

import (
	"runtime"
	"syscall"

	"radithor/core"

	"github.com/horobimasu/windowsx"
	"golang.org/x/sys/windows"
)

const _CONTROL_CENTER_WINDOW_NAME = "ControlCenterWindow"

var (
	wasOpen bool

	controlCenterThreadId uint32

	controlCenterCallback = syscall.NewCallback(onControlCenterWindowOpen)
)

func EnableControlCenterWindowDetection() {
	wasOpen = false

	go func() {
		runtime.LockOSThread()

		var msg windowsx.Msg
		windowsx.PeekMessage(&msg, 0, _WM_USER, _WM_USER, windowsx.PM_NOREMOVE)

		controlCenterThreadId = windows.GetCurrentThreadId()

		hook, err := windowsx.SetEventHook(
			windowsx.EVENT_SYSTEM_FOREGROUND,
			windowsx.EVENT_SYSTEM_FOREGROUND,
			0,
			controlCenterCallback,
			0,
			0,
			windowsx.WINEVENT_OUTOFCONTEXT,
		)

		if err != nil {
			core.LogError("[EnableControlCenterWindowDetection] failed to set hook")
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

func DisableControlCenterWindowDetection() {
	windowsx.PostThreadMessage(controlCenterThreadId, _WM_QUIT, 0, 0)
}

func onControlCenterWindowOpen(
	hook uintptr,
	event uintptr,
	handle uintptr,
	object uintptr,
	child uintptr,
	thread uintptr,
	time uintptr,
) uintptr {
	className, _ := windowsx.GetClassName(handle)

	open := className == _CONTROL_CENTER_WINDOW_NAME
	if open && !wasOpen {
		panicModeCallback()
	}

	wasOpen = open

	return 0
}
