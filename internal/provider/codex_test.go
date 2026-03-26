package provider

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func TestCodexProviderIncludesCacheLogSessionTargets(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("failed to resolve user home: %v", err)
	}

	p, err := Get("codex")
	if err != nil {
		t.Fatalf("failed to load codex provider: %v", err)
	}

	expectedPaths := []string{
		filepath.Join(home, ".codex", "models_cache.json"),
		filepath.Join(home, ".codex", "log"),
		filepath.Join(home, ".codex", "logs_*.sqlite*"),
		filepath.Join(home, ".codex", "sessions"),
		filepath.Join(home, ".codex", "session_index.jsonl"),
		filepath.Join(home, ".codex", "history.jsonl"),
	}

	for _, want := range expectedPaths {
		if !slices.Contains(p.CachePaths, want) {
			t.Fatalf("codex cache paths missing %q", want)
		}
	}

	expectedChecks := []string{"cache files", "log files", "session files"}
	for _, want := range expectedChecks {
		found := false
		for _, check := range p.Checks {
			if check.Name == want {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("codex checks missing %q", want)
		}
	}
}
