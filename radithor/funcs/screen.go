package funcs

import (
	"fmt"
	"syscall"
	"unsafe"

	"github.com/horobimasu/windowsx"
	"golang.org/x/sys/windows"
)

// todo, set brightness function with LCD driver
// todo, add all the video settings maybe

var videoSubGroup = syscall.GUID{
	Data1: 0x7516B95F,
	Data2: 0xF776,
	Data3: 0x4464,
	Data4: [8]byte{0x8C, 0x53, 0x06, 0x16, 0x7F, 0x40, 0xCC, 0x99},
}

var (
	videoBrightness = syscall.GUID{
		Data1: 0xADED5E82,
		Data2: 0xB909,
		Data3: 0x4619,
		Data4: [8]byte{0x99, 0x49, 0xF5, 0xD7, 0x1D, 0xAC, 0x0B, 0xCB},
	}

	videoPowerDown = syscall.GUID{
		Data1: 0x3C0BC021,
		Data2: 0xC8A8,
		Data3: 0x4E07,
		Data4: [8]byte{0xA9, 0x73, 0x6B, 0x14, 0xCB, 0xCB, 0x2B, 0x7E},
	}
)

func SetScreenOrientation(screenName string, orientation uint32) error {
	var screen *uint16

	if screenName != "" {
		name, err := windows.UTF16PtrFromString(screenName)
		if err != nil {
			return err
		}

		screen = name
	}

	var device windowsx.DeviceMode
	device.Size = uint16(unsafe.Sizeof(device))

	err := windowsx.EnumDisplaySettings(screen, windowsx.ENUM_CURRENT_SETTINGS, &device)
	if err != nil {
		return err
	}

	if (device.DisplayOrientation^orientation)&1 == 1 {
		device.PelsWidth = device.PelsHeight
		device.PelsHeight = device.PelsWidth
	}

	device.DisplayOrientation = orientation
	device.Fields = windowsx.DM_DISPLAYORIENTATION | windowsx.DM_PELSWIDTH | windowsx.DM_PELSHEIGHT

	ret := windowsx.ChangeDisplaySettingsEx(screen, &device, 0, windowsx.CDS_TEST, 0)
	if ret != windowsx.DISP_CHANGE_SUCCESSFUL {
		return fmt.Errorf("failed to test display change: error code %d", ret)
	}

	ret = windowsx.ChangeDisplaySettingsEx(screen, &device, 0, windowsx.CDS_UPDATEREGISTRY, 0)
	if ret != windowsx.DISP_CHANGE_SUCCESSFUL {
		return fmt.Errorf("failed to change display: error code %d", ret)
	}

	return nil
}

func SetScreenBrightness(percent uint32) error {
	return internalSetDisplaySetting(videoBrightness, percent)
}

func SetScreenTimeout(seconds uint32) error {
	return internalSetDisplaySetting(videoPowerDown, seconds)
}

func internalSetDisplaySetting(setting syscall.GUID, index uint32) error {
	var scheme *syscall.GUID

	err := windowsx.PowerGetActiveScheme(0, &scheme)
	if err != nil {
		return err
	}

	defer windows.LocalFree(windows.Handle(unsafe.Pointer(scheme)))

	err = windowsx.PowerWriteACValueIndex(0, scheme, &videoSubGroup, &setting, index)
	if err != nil {
		return err
	}

	err = windowsx.PowerWriteDCValueIndex(0, scheme, &videoSubGroup, &setting, index)
	if err != nil {
		return err
	}

	return windowsx.PowerSetActiveScheme(0, scheme)
}
