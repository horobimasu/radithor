package tests

import (
	"testing"

	"radithor/core/core_funcs"
)

func TestGenFingerprint(t *testing.T) {
	fingerprint, err := core_funcs.GenerateFingerprint()
	if err != nil {
		t.Fatalf("failed to generate fingerprint - %v", err)
	}

	t.Logf("fingerprint: %x", fingerprint)
}
