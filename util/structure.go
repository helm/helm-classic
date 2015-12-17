package util

import "path/filepath"

// cachePath is the suffix for the cache.
const cachePath = "cache"

// workspacePath is the user's workspace directory.
const workspacePath = "workspace"

// workspaceChartPath is the directory that contains a user's workspace charts.
const workspaceChartPath = "workspace/charts"

// CacheDirectory - File path to cache directory based on home
func CacheDirectory(home string, paths ...string) string {
	fragments := append([]string{home, cachePath}, paths...)
	return filepath.Join(fragments...)
}

// WorkspaceChartDirectory - File path to workspace chart directory based on home
func WorkspaceChartDirectory(home string, paths ...string) string {
	fragments := append([]string{home, workspaceChartPath}, paths...)
	return filepath.Join(fragments...)
}
