package funcs

import (
	"errors"
	"math/rand"
	"time"

	"github.com/horobimasu/windowsx"
	"golang.org/x/sys/windows"
)

const _BOUNCE_FPS int64 = 60

// the original plan was to use dialog boxes like how
// endermanch did it with 2989.exe
// but that needed some retarded version thing in comctl32
// which i did NOT want to do, especially in golang
func BounceMessagebox(tempTitle string, newTitle string, speed float64) error {
	var handle uintptr

	timeout := time.Now().Add(6 * time.Second)
	for time.Now().Before(timeout) {
		theHandle, err := windowsx.FindWindow(
			windows.StringToUTF16Ptr("#32770"),
			windows.StringToUTF16Ptr(tempTitle),
		)

		if err == nil && theHandle != 0 {
			handle = theHandle
			break
		}

		time.Sleep(50 * time.Millisecond)
	}

	if handle == 0 {
		return errors.New("failed to find handle")
	}

	err := windowsx.SetWindowText(
		handle,
		windows.StringToUTF16Ptr(newTitle),
	)

	if err != nil {
		return err
	}

	return BounceWindow(handle, speed)
}

func BounceWindow(handle uintptr, speed float64) error {
	rect, err := windowsx.GetWindowRect(handle)
	if err != nil {
		return err
	}

	windowW := float64(rect.Right - rect.Left)
	windowH := float64(rect.Bottom - rect.Top)

	windowX := float64(rect.Left)
	windowY := float64(rect.Top)

	screenW := float64(windowsx.GetSystemMetrics(windowsx.SM_CXSCREEN))
	screenH := float64(windowsx.GetSystemMetrics(windowsx.SM_CYSCREEN))

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	dirX := speed * rng.Float64()
	dirY := speed * rng.Float64()

	ticker := time.NewTicker(time.Duration(1000/_BOUNCE_FPS) * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		if !windows.IsWindow(windows.HWND(handle)) {
			return nil
		}

		windowX += dirX
		windowY += dirY

		if windowX <= 0 {
			windowX = 0
			dirX = speed
		} else if windowX+windowW >= screenW {
			windowX = screenW - windowW
			dirX = -speed
		}

		if windowY <= 0 {
			windowY = 0
			dirY = speed
		} else if windowY+windowH >= screenH {
			windowY = screenH - windowH
			dirY = -speed
		}

		err := windowsx.SetWindowPos(
			handle,
			0,
			int32(windowX),
			int32(windowY),
			0,
			0,
			windowsx.SWP_NOSIZE|windowsx.SWP_NOZORDER|windowsx.SWP_NOACTIVATE,
		)

		if err != nil {
			return err
		}
	}

	return nil
}
