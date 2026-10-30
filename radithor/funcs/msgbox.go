package funcs

import "golang.org/x/sys/windows"

func MessageBox(title string, text string, flags uint32) (int32, error) {
	return windows.MessageBox(
		windows.HWND(0),
		windows.StringToUTF16Ptr(text),
		windows.StringToUTF16Ptr(title),
		flags|windows.MB_TOPMOST|windows.MB_SETFOREGROUND,
	)
}
