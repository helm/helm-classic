package chart

import (
	"strings"
	"testing"
)

func TestRepoName(t *testing.T) {
	name := RepoName(".")
	if name == "" {
		t.Fatalf("Expected a git URL.")
	}
	if !strings.HasSuffix(name, ".git") {
		t.Errorf("Expected %s to end with '.git'", name)
	}
}
