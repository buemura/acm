package scanner

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func TestScanFilesSupportsGlobPatternsAndDedupes(t *testing.T) {
	tmpDir := t.TempDir()
	logMain := filepath.Join(tmpDir, "logs_1.sqlite")
	logWal := filepath.Join(tmpDir, "logs_1.sqlite-wal")

	if err := os.WriteFile(logMain, []byte("main"), 0o644); err != nil {
		t.Fatalf("failed to create %s: %v", logMain, err)
	}
	if err := os.WriteFile(logWal, []byte("wal"), 0o644); err != nil {
		t.Fatalf("failed to create %s: %v", logWal, err)
	}

	files, err := ScanFiles([]string{
		filepath.Join(tmpDir, "logs_*.sqlite*"),
		logMain,
	}, FilterOpts{})
	if err != nil {
		t.Fatalf("scan failed: %v", err)
	}

	if len(files) != 2 {
		t.Fatalf("expected 2 files after dedupe, got %d", len(files))
	}

	paths := make([]string, 0, len(files))
	for _, f := range files {
		paths = append(paths, f.Path)
	}
	if !slices.Contains(paths, logMain) {
		t.Fatalf("expected scan results to include %s", logMain)
	}
	if !slices.Contains(paths, logWal) {
		t.Fatalf("expected scan results to include %s", logWal)
	}
}
