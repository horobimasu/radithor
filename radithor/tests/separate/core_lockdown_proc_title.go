package main

// import (
// 	"fmt"
// 	"runtime"
// 	"syscall"

// 	"github.com/horobimasu/windowsx"
// )

// const _WM_USER = 0x0400

// var callback = syscall.NewCallback(onDetect)

// func main() {
// 	go func() {
// 		runtime.LockOSThread()

// 		var msg windowsx.Msg
// 		windowsx.PeekMessage(&msg, 0, _WM_USER, _WM_USER, windowsx.PM_NOREMOVE)

// 		hook, err := windowsx.SetEventHook(
// 			windowsx.EVENT_OBJECT_NAMECHANGE,
// 			windowsx.EVENT_OBJECT_NAMECHANGE,
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

// 	title, err := windowsx.GetWindowText(handle)
// 	if err != nil {
// 		return 0
// 	}

// 	println(title)

// 	return 0
// }
