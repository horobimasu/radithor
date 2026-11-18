package tests

import (
	"testing"

	"radithor/funcs"

	"github.com/horobimasu/gotils"
	"golang.org/x/sys/windows"
)

const (
	_BOUNCE_SPEED = 10
	_REAL_TITLE   = "bouncy box title"
)

func TestBounceMessagebox(t *testing.T) {
	tempTitle := gotils.GenRandomString(9, gotils.CHARSET_ALPHABETIC)

	go func() {
		_, err := funcs.MessageBox(tempTitle, "bouncy bouncy", windows.MB_OK)
		if err != nil {
			t.Fatalf("failed to spawn msgbox - %v", err)
		}
	}()

	err := funcs.BounceMessageBox(tempTitle, _REAL_TITLE, _BOUNCE_SPEED)
	if err != nil {
		t.Fatalf("failed to bounce msgbox - %v", err)
	}
}
