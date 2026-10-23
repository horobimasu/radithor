package funcs

import "golang.org/x/sys/windows/registry"

const _EXEC_OPTIONS = "SOFTWARE\\Microsoft\\Windows NT\\CurrentVersion\\Image File Execution Options\\"

func BlockProcess(name string) error {
	key, _, err := registry.CreateKey(registry.LOCAL_MACHINE, _EXEC_OPTIONS+name, registry.SET_VALUE)
	if err != nil {
		return err
	}

	defer key.Close()

	return key.SetStringValue("Debugger", "nul")
}

func UnblockProcess(name string) error {
	err := registry.DeleteKey(registry.LOCAL_MACHINE, _EXEC_OPTIONS+name)
	if err == nil {
		return nil
	}

	key, err := registry.OpenKey(registry.LOCAL_MACHINE, _EXEC_OPTIONS+name, registry.SET_VALUE)
	if err != nil {
		return err
	}

	defer key.Close()

	return key.DeleteValue("Debugger")
}
