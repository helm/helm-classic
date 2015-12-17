package util

import "path/filepath"

// CachePath is the suffix for the cache.
const CachePath = "cache"

// WorkspacePath is the user's workspace directory.
const WorkspacePath = "workspace"

// WorkspaceChartPath is the directory that contains a user's workspace charts.
const WorkspaceChartPath = "workspace/charts"

// CacheDirectory - File path to cache directory based on home
func CacheDirectory(home string, paths ...string) string {
	fragments := append([]string{home, CachePath}, paths...)
	return filepath.Join(fragments...)
}
