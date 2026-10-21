package main

// import (
// 	"fmt"
// 	"runtime"
// 	"syscall"

// 	"radithor/utils"

// 	"github.com/horobimasu/windowsx"
// 	"golang.org/x/sys/windows"
// )

// const _WM_USER = 0x0400

// var callback = syscall.NewCallback(onDetect)

// func main() {
// 	go func() {
// 		runtime.LockOSThread()

// 		var msg windowsx.Msg
// 		windowsx.PeekMessage(&msg, 0, _WM_USER, _WM_USER, windowsx.PM_NOREMOVE)

// 		hook, err := windowsx.SetEventHook(
// 			windowsx.EVENT_OBJECT_CREATE,
// 			windowsx.EVENT_OBJECT_CREATE,
// 			0,
// 			callback,
// 			0,
// 			0,
// 			windowsx.WINEVENT_OUTOFCONTEXT,
// 		)

// 		if err != nil {
// 			fmt.Println("failed to set hook")
// 			return
// 		}

// 		for {
// 			ret, _ := windowsx.GetMessage(&msg, 0, 0, 0)
// 			if ret <= 0 {
// 				break
// 			}
// 		}

// 		windowsx.UnhookEvent(hook)
// 	}()

// 	fmt.Println("started")
// 	fmt.Scanln()
// }

// func onDetect(
// 	hook uintptr,
// 	event uintptr,
// 	handle uintptr,
// 	object uintptr,
// 	child uintptr,
// 	thread uintptr,
// 	time uintptr,
// ) uintptr {
// 	if int32(object) != 0 {
// 		return 0
// 	}

// 	var pid uint32
// 	windows.GetWindowThreadProcessId(windows.HWND(handle), &pid)

// 	name, err := utils.GetProcNameFromPid(pid)
// 	if err != nil {
// 		return 0
// 	}

// 	fmt.Println(name)

// 	return 0
// }
