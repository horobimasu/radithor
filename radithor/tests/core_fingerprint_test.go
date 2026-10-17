package tests

import (
	"testing"

	"radithor/core"
)

func TestGenFingerprint(t *testing.T) {
	fingerprint, err := core.GenerateFingerprint()
	if err != nil {
		t.Fatalf("failed to generate fingerprint - %v", err)
	}

	t.Logf("fingerprint: %x", fingerprint)
}
