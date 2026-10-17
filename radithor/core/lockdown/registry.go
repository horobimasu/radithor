package lockdown

import (
	"errors"
	"radithor/funcs"
	"radithor/store"

	"golang.org/x/sys/windows/registry"
)

const _SYSTEM_POLICIES = "SOFTWARE\\Policies\\Microsoft\\Windows\\System"

func LockdownRegistry() error {
	if !store.IsElevated {
		return errors.New("elevation is required to run this function")
	}

	for _, proc := range ProcsToBlock {
		funcs.BlockProcess(proc)
	}

	key, _, err := registry.CreateKey(registry.LOCAL_MACHINE, _SYSTEM_POLICIES, registry.SET_VALUE)
	if err != nil {
		return err
	}

	defer key.Close()

	return key.SetDWordValue("DontDisplayNetworkSelectionUI", 1)
}

func UnlockdownRegistry() error {
	if !store.IsElevated {
		return errors.New("elevation is required to run this function")
	}

	for _, proc := range ProcsToBlock {
		funcs.UnblockProcess(proc)
	}

	key, _, err := registry.CreateKey(registry.LOCAL_MACHINE, _SYSTEM_POLICIES, registry.SET_VALUE)
	if err != nil {
		return err
	}

	defer key.Close()

	return key.DeleteValue("DontDisplayNetworkSelectionUI")
}
