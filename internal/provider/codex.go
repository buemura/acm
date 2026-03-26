package provider

import (
	"os"
	"path/filepath"
)

func init() {
	home, _ := os.UserHomeDir()
	codexDir := filepath.Join(home, ".codex")
	Register(&Provider{
		Name: "codex",
		CachePaths: []string{
			filepath.Join(codexDir, "cache"),
			filepath.Join(home, ".cache", "codex"),
			filepath.Join(codexDir, "models_cache.json"),
			filepath.Join(codexDir, "log"),
			filepath.Join(codexDir, "logs_*.sqlite*"),
			filepath.Join(codexDir, "sessions"),
			filepath.Join(codexDir, "session_index.jsonl"),
			filepath.Join(codexDir, "history.jsonl"),
		},
		Checks: []PathCheck{
			{
				Name: "cache files",
				Patterns: []string{
					filepath.Join(codexDir, "cache"),
					filepath.Join(home, ".cache", "codex"),
					filepath.Join(codexDir, "models_cache.json"),
				},
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
