package tests

import (
	"testing"

	"radithor/funcs"

	"shared"

	"golang.org/x/sys/windows"
)

const (
	_BOUNCE_SPEED = 10
	_REAL_TITLE   = "bouncy box title"
)

func TestBounceMessagebox(t *testing.T) {
	tempTitle, _ := shared.GenerateRandomString(9)

	go func() {
		_, err := funcs.MessageBox(tempTitle, "bouncy bouncy", windows.MB_OK)
		if err != nil {
			t.Fatalf("failed to spawn msgbox - %v", err)
		}
	}()

	err := funcs.BounceMessagebox(tempTitle, _REAL_TITLE, _BOUNCE_SPEED)
	if err != nil {
		t.Fatalf("failed to bounce msgbox - %v", err)
	}
}
