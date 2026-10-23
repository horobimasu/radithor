package tests

import (
	"testing"

	"radithor/utils"
)

const (
	_SHORTCUT_PATH   = "%USERPROFILE%\\Desktop\\Test.lnk"
	_SHORTCUT_TARGET = "calc.exe"
)

func TestCreateShortcut(t *testing.T) {
	data := utils.Shortcut{
		ShortcutPath: _SHORTCUT_PATH,
		Target:       _SHORTCUT_TARGET,
		IconPath:     _SHORTCUT_TARGET,
	}

	err := utils.CreateShortcut(data)
	if err != nil {
		t.Fatalf("failed to create shortcut - %v", err)
		return
	}

	t.Logf("created shortcut")
}

func TestGetShortcutData(t *testing.T) {
	data, err := utils.GetShortcutData(_SHORTCUT_PATH)
	if err != nil {
		t.Fatalf("failed to get data - %v", err)
		return
	}

	t.Logf("data - %+v\n", data)
}

func TestUpdateShortcutData(t *testing.T) {
	data := utils.Shortcut{
		Target: _SHORTCUT_TARGET,
	}

	err := utils.UpdateShortcutData(_SHORTCUT_PATH, data)
	if err != nil {
		t.Fatalf("failed to update - %v", err)
		return
	}

	t.Logf("updated shortcut")
}
