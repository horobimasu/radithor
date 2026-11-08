package tests

import (
	"testing"

	"radithor/funcs"
)

const (
	_FILE_PATH        = "%USERPROFILE%\\Desktop\\Test.exe"
	_FILE_DESCRIPTION = "_FILE_DESCRIPTION"

	_FILE_TO_GET_ICON_FROM_PATH = "calc.exe"
)

func TestSetExeFileFileDescription(t *testing.T) {
	err := funcs.SetExeFileFileDescription(_FILE_PATH, _FILE_DESCRIPTION)
	if err != nil {
		t.Fatalf("failed to set file description - %v", err)
	}
}

func TestSetExeFileFileIcon(t *testing.T) {
	err := funcs.SetExeFileFileIcon(_FILE_PATH, _FILE_TO_GET_ICON_FROM_PATH)
	if err != nil {
		t.Fatalf("failed to set file icon - %v", err)
	}
}
