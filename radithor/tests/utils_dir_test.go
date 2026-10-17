package tests

import (
	"testing"

	"radithor/utils"
)

func TestGetRandDir(t *testing.T) {
	for i := range 30 {
		path, err := utils.GetRandomDir("~", 4, 20, []string{})
		if err != nil {
			t.Fatal(err)
		}

		t.Logf("%d -> %s", i+1, path)
	}
}
