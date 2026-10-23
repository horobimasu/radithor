package persistence

import (
	"os"
	"path/filepath"

	"radithor/utils"

	"golang.org/x/sys/windows"
)

func AddToStartupFolder(theFilePath string) error {
	path, err := windows.GetCurrentProcessToken().KnownFolderPath(windows.FOLDERID_Startup, 0)
	if err != nil {
		return err
	}

	return utils.CreateShortcut(utils.Shortcut{
		ShortcutPath: filepath.Join(path, getPersistenceName(_PERSISTENCE_NAME_STARTUP)),
		Target:       theFilePath,
		IconPath:     theFilePath,
	})
}

func RemoveFromStartupFolder() error {
	path, err := windows.GetCurrentProcessToken().KnownFolderPath(windows.FOLDERID_Startup, 0)
	if err != nil {
		return err
	}

	return os.Remove(filepath.Join(path, getPersistenceName(_PERSISTENCE_NAME_STARTUP)))
}
