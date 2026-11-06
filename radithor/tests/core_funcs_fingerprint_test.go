package tests

import (
	"testing"

	"radithor/core/core_funcs"
)

func TestGenFingerprint(t *testing.T) {
	t.Logf("fingerprint: %x", core_funcs.GenerateFingerprint())
}

func TestGenUnhashedFingerprint(t *testing.T) {
	for name, value := range core_funcs.GetUnhashedFingerprint() {
		t.Logf("%s: %s", name, value)
	}
}
