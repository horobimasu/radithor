package funcs

import (
	"errors"
	"runtime"
	"time"
	"unsafe"

	"github.com/horobimasu/windowsx"
	"golang.org/x/sys/windows"
)

func GetClipboardText() (string, error) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	var opened bool

	for range 50 {
		err := windowsx.OpenClipboard(0)
		if err == nil {
			opened = true
			break
		}

		time.Sleep(50 * time.Millisecond)
	}

	if !opened {
		return "", errors.New("failed to open clipboard")
	}

	defer windowsx.CloseClipboard()

	handle, err := windowsx.GetClipboardData(windowsx.CF_UNICODETEXT)
	if err != nil {
		return "", err
	}

	ptr, err := windowsx.GlobalLock(handle)
	if err != nil {
		return "", err
	}

	defer windowsx.GlobalUnlock(handle)

	size, err := windowsx.GlobalSize(handle)
	if err != nil {
		return "", err
	}

	return windows.UTF16ToString(unsafe.Slice(
		(*uint16)(unsafe.Pointer(ptr)),
		int(size)/2,
	)), nil
}

func CopyToClipboard(text string) error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	var opened bool

	for range 50 {
		err := windowsx.OpenClipboard(0)
		if err == nil {
			opened = true
			break
		}

		time.Sleep(50 * time.Millisecond)
	}

	if !opened {
		return errors.New("failed to open clipboard")
	}

	defer windowsx.CloseClipboard()

	err := windowsx.EmptyClipboard()
	if err != nil {
		return err
	}

	utf16Text, err := windows.UTF16FromString(text)
	if err != nil {
		return err
	}

	handle, err := windowsx.GlobalAlloc(
		windowsx.GMEM_MOVEABLE,
		uintptr(len(utf16Text)*2),
	)

	if err != nil {
		return err
	}

	ptr, err := windowsx.GlobalLock(handle)
	if err != nil {
		windowsx.GlobalFree(handle)
		return err
	}

	copy(unsafe.Slice((*uint16)(unsafe.Pointer(ptr)), len(utf16Text)), utf16Text)

	windowsx.GlobalUnlock(handle)

	_, err = windowsx.SetClipboardData(windowsx.CF_UNICODETEXT, handle)
	if err != nil {
		windowsx.GlobalFree(handle)
		return err
	}

	return nil
}
