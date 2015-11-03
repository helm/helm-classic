package release

import (
	"testing"
)

func TestLatest(t *testing.T) {
	rr, err := Latest()
	if err != nil {
		t.Errorf("Failed to get latest: %s", err)
	}

	if *rr.ID <= 0 {
		t.Errorf("ID below zero.")
	}
}

func TestLatestVersion(t *testing.T) {
	v, err := LatestVersion()
	if err != nil {
		t.Error(err)
	}

	if v == "" {
		t.Error("Expected tag, not empty string")
	}
}

func TestLatestDownloadURL(t *testing.T) {
	v, err := LatestDownloadURL()
	if err != nil {
		t.Error(err)
	}

	if v == "" {
		t.Error("Expected URL, not empty string")
	}

}
