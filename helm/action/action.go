package action

import (
	"os"
	"path"
)

// CachePath is path to repository checkout
var CachePath string

// CacheChartPath is path to charts inside the checkout
var CacheChartPath string

// WorkspacePath is path to the local workspace
var WorkspacePath string

// WorkspaceChartPath is path to charts inside the workspace
var WorkspaceChartPath string

// DefaultNS to use during kubectl command execution
const DefaultNS = "default"

var helmpaths = []string{CachePath, WorkspacePath}

func init() {

	envCachePath := os.ExpandEnv("$HELM_CACHE")
	if envCachePath == "" {
		CachePath = "cache"
	} else {
		CachePath = envCachePath
	}
	CacheChartPath = path.Join(CachePath, "charts")

	envWorkspacePath := os.ExpandEnv("$HELM_WORKSPACE")
	if envWorkspacePath == "" {
		WorkspacePath = "workspace"
	} else {
		WorkspacePath = envWorkspacePath
	}
	WorkspaceChartPath = path.Join(WorkspacePath, "charts")

}
