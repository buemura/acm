package provider

import (
	"os"
	"path/filepath"
)

func init() {
	home, _ := os.UserHomeDir()
	claudeDir := filepath.Join(home, ".claude")
	Register(&Provider{
		Name: "claude",
		CachePaths: []string{
			filepath.Join(claudeDir, "cache"),
			filepath.Join(home, ".cache", "claude"),
			filepath.Join(claudeDir, "projects"),
			filepath.Join(claudeDir, "stats-cache.json"),
		},
		Checks: []PathCheck{
			{
				Name: "cache files",
				Patterns: []string{
					filepath.Join(claudeDir, "cache"),
					filepath.Join(home, ".cache", "claude"),
				},
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
