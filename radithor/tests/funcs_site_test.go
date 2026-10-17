package tests

import (
	"testing"

	"radithor/funcs"
)

// remember to run as admin
// remember to remove site from %SYSTEMROOT%\System32\drivers\etc\host

const _SITE_TO_TEST = "https://www.virustotal.com/gui/home/upload"

func TestBlockWebsite(t *testing.T) {
	site, err := funcs.BlockWebsite(_SITE_TO_TEST)
	if err != nil {
		t.Fatalf("failed to block site - %v", err)
	}

	t.Logf("blocked site %q", site)
}

func TestUnblockWebsite(t *testing.T) {
	site, err := funcs.UnblockWebsite(_SITE_TO_TEST)
	if err != nil {
		t.Fatalf("failed to unblock site - %v", err)
	}

	t.Logf("unblocked site %q", site)
}
