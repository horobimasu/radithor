package main

// import (
// 	"fmt"
// 	"runtime"
// 	"syscall"

// 	"github.com/horobimasu/windowsx"
// )

// const (
// 	_CONTROL_CENTER_WINDOW_NAME = "ControlCenterWindow"

// 	_WM_USER = 0x0400
// )

// var (
// 	wasOpen bool

// 	callback = syscall.NewCallback(onDetect)
// )

// func main() {
// 	wasOpen = false

// 	go func() {
// 		runtime.LockOSThread()

// 		var msg windowsx.Msg
// 		windowsx.PeekMessage(&msg, 0, _WM_USER, _WM_USER, windowsx.PM_NOREMOVE)

// 		hook, err := windowsx.SetEventHook(
// 			windowsx.EVENT_SYSTEM_FOREGROUND,
// 			windowsx.EVENT_SYSTEM_FOREGROUND,
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

// 	className, err := windowsx.GetClassName(handle)
// 	if err != nil {
// 		fmt.Printf("failed to get class name - %v\n", err)
// 		return 0
// 	}

// 	open := className == _CONTROL_CENTER_WINDOW_NAME
// 	if open && !wasOpen {
// 		fmt.Println("control center opened")
// 	}

// 	wasOpen = open

// 	return 0
// }
