package provider

import (
	"path/filepath"
)

func init() {
	codexDir := buildAppDir("codex")

	// Build cache paths (includes platform-specific locations)
	cachePaths := getAllCachePaths("codex")

	// Add non-cache paths
	cachePaths = append(cachePaths,
		filepath.Join(codexDir, "models_cache.json"),
		filepath.Join(codexDir, "log"),
		filepath.Join(codexDir, "logs_*.sqlite*"),
		filepath.Join(codexDir, "sessions"),
		filepath.Join(codexDir, "session_index.jsonl"),
		filepath.Join(codexDir, "history.jsonl"),
	)

	// Build check patterns for cache directories
	cacheCheckPatterns := getAllCachePaths("codex")
	cacheCheckPatterns = append(cacheCheckPatterns, filepath.Join(codexDir, "models_cache.json"))

	Register(&Provider{
		Name:       "codex",
		CachePaths: cachePaths,
		Checks: []PathCheck{
			{
				Name:     "cache files",
				Patterns: cacheCheckPatterns,
			},
			{
				Name: "log files",
				Patterns: []string{
					filepath.Join(codexDir, "log"),
					filepath.Join(codexDir, "logs_*.sqlite*"),
				},
			},
			{
				Name: "session files",
				Patterns: []string{
					filepath.Join(codexDir, "sessions"),
					filepath.Join(codexDir, "session_index.jsonl"),
					filepath.Join(codexDir, "history.jsonl"),
				},
			},
		},
	})
}
