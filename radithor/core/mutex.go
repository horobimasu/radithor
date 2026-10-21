package core

import (
	"errors"
	"os"
	"path/filepath"

	"radithor/store"

	"shared/config"

	"golang.org/x/sys/windows"
)

const _FILE_SHARE_0 uint32 = 0x00000000

var mutexHandle windows.Handle = windows.InvalidHandle

func CreateMutex() error {
	if mutexHandle != windows.InvalidHandle {
		return errors.New("mutex already active")
	}

	utf16Path := windows.StringToUTF16Ptr(filepath.Join(
		os.TempDir(),
		getMutexName(),
	))

	handle, err := windows.CreateFile(
		utf16Path,
		windows.GENERIC_READ|windows.GENERIC_WRITE,
		_FILE_SHARE_0,
		nil,
		windows.OPEN_ALWAYS,
		windows.FILE_ATTRIBUTE_HIDDEN|windows.FILE_ATTRIBUTE_SYSTEM|windows.FILE_FLAG_DELETE_ON_CLOSE,
		0,
	)

	if err != nil {
		return err
	}

	mutexHandle = handle

	return nil
}

func RemoveMutex() error {
	if mutexHandle == windows.InvalidHandle {
		return errors.New("mutex in use")
	}

	err := windows.CloseHandle(mutexHandle)
	if err != nil {
		return err
	}

	mutexHandle = windows.InvalidHandle

	err = os.Remove(filepath.Join(os.TempDir(), getMutexName()))
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	return nil
}

func CheckMutexActive() bool {
	err := os.Remove(filepath.Join(os.TempDir(), getMutexName()))
	if err == nil || os.IsNotExist(err) {
		return false
	}

	return true
}

func getMutexName() string {
	if store.DEBUG_MODE_ENABLED && store.DEBUG_MODE_MUTEX_FILE_NAME != "" {
		return store.DEBUG_MODE_MUTEX_FILE_NAME
	} else {
		return config.GetConfigString(store.MUTEX_NAME)
	}
}
