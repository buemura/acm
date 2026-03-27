package provider

import (
	"path/filepath"
)

func init() {
	claudeDir := buildAppDir("claude")

	// Build cache paths (includes platform-specific locations)
	cachePaths := getAllCachePaths("claude")

	// Add non-cache paths
	cachePaths = append(cachePaths,
		filepath.Join(claudeDir, "projects"),
		filepath.Join(claudeDir, "stats-cache.json"),
	)

	// Build check patterns for cache directories
	cacheCheckPatterns := getAllCachePaths("claude")

	Register(&Provider{
		Name:       "claude",
		CachePaths: cachePaths,
		Checks: []PathCheck{
			{
				Name:     "cache files",
				Patterns: cacheCheckPatterns,
			},
			{
				Name: "session transcripts",
				Patterns: []string{
					filepath.Join(claudeDir, "projects", "*", "*.jsonl"),
				},
			},
			{
				Name: "stats cache",
				Patterns: []string{
					filepath.Join(claudeDir, "stats-cache.json"),
				},
			},
		},
	})
}
