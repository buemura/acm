package provider

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func TestMissingPathChecks(t *testing.T) {
	tmpDir := t.TempDir()
	existing := filepath.Join(tmpDir, "present.log")

	touchFile(t, existing)

	p := &Provider{
		Name: "test",
		Checks: []PathCheck{
			{
				Name:     "present files",
				Patterns: []string{existing},
			},
			{
				Name:     "missing files",
				Patterns: []string{filepath.Join(tmpDir, "does-not-exist*")},
			},
		},
	}

	missing := MissingPathChecks(p)
	if len(missing) != 1 {
		t.Fatalf("expected 1 missing check, got %d: %v", len(missing), missing)
	}
	if !slices.Contains(missing, "missing files") {
		t.Fatalf("expected missing checks to include %q: %v", "missing files", missing)
	}
}

func touchFile(t *testing.T, path string) {
	t.Helper()
	if err := os.WriteFile(path, []byte("x"), 0o644); err != nil {
		t.Fatalf("failed to create test file %s: %v", path, err)
	}
}
