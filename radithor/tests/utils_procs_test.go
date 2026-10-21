package tests

import (
	"strings"
	"testing"

	"radithor/utils"
)

func TestGetAllProcsTitles(t *testing.T) {
	all, err := utils.GetAllRunningProcs()
	if err != nil {
		t.Fatal(err)
	}

	for _, proc := range all {
		titles, err := utils.GetWindowTitlesFromPid(proc.Pid)
		if err != nil {
			t.Logf("failed to get titles - %v", err)
		}

		if len(titles) == 0 {
			continue
		}

		t.Logf("%s - [%s]", proc.Name, strings.Join(titles, ", "))
	}
}
