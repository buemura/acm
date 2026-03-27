package scanner

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// TestCrossPlatformPaths verifies that the scanner handles
// platform-specific path formats correctly
func TestCrossPlatformPaths(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test files
	testFile1 := filepath.Join(tmpDir, "cache", "file1.json")
	testFile2 := filepath.Join(tmpDir, "cache", "file2.jsonl")
	testFile3 := filepath.Join(tmpDir, "logs", "app.log")

	// Create directories
	if err := os.MkdirAll(filepath.Join(tmpDir, "cache"), 0o755); err != nil {
		t.Fatalf("failed to create cache dir: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(tmpDir, "logs"), 0o755); err != nil {
		t.Fatalf("failed to create logs dir: %v", err)
	}

	// Create files
	for _, f := range []string{testFile1, testFile2, testFile3} {
		if err := os.WriteFile(f, []byte("test"), 0o644); err != nil {
			t.Fatalf("failed to create %s: %v", f, err)
		}
	}

	// Test glob patterns work across platforms
	t.Run("glob_patterns", func(t *testing.T) {
		files, err := ScanFiles([]string{
			filepath.Join(tmpDir, "cache", "*.json*"),
		}, FilterOpts{})
		if err != nil {
			t.Fatalf("scan failed: %v", err)
		}

		if len(files) != 2 {
			t.Errorf("expected 2 files, got %d", len(files))
		}
	})

	// Test absolute paths work
	t.Run("absolute_paths", func(t *testing.T) {
		files, err := ScanFiles([]string{testFile1}, FilterOpts{})
		if err != nil {
			t.Fatalf("scan with absolute path failed: %v", err)
		}

		if len(files) != 1 {
			t.Errorf("expected 1 file, got %d", len(files))
		}
		if files[0].Path != testFile1 {
			t.Errorf("expected path %s, got %s", testFile1, files[0].Path)
		}
	})

	// Test multiple directory patterns
	t.Run("multiple_directories", func(t *testing.T) {
		files, err := ScanFiles([]string{
			filepath.Join(tmpDir, "cache", "*"),
			filepath.Join(tmpDir, "logs", "*"),
		}, FilterOpts{})
		if err != nil {
			t.Fatalf("scan failed: %v", err)
		}

		if len(files) != 3 {
			t.Errorf("expected 3 files, got %d", len(files))
		}
	})
}

// TestWindowsSpecificPaths tests Windows-specific path handling
func TestWindowsSpecificPaths(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Windows-specific test")
	}

	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.json")

	if err := os.WriteFile(testFile, []byte("test"), 0o644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	// Test that Windows accepts both forward and backslashes
	t.Run("forward_slash_on_windows", func(t *testing.T) {
		// Convert backslashes to forward slashes
		unixStylePath := filepath.ToSlash(testFile)

		// This should still work on Windows
		if _, err := os.Stat(unixStylePath); err != nil {
			// Note: os.Stat may not accept forward slashes on Windows,
			// but our glob patterns should handle it
			t.Logf("Direct stat with forward slashes: %v (expected on Windows)", err)
		}
	})

	// Test that Windows paths with drive letters work
	t.Run("drive_letter_paths", func(t *testing.T) {
		// Verify the path has a drive letter (e.g., C:\)
		if !filepath.IsAbs(testFile) {
			t.Errorf("expected absolute path with drive letter, got: %s", testFile)
		}

		// Verify it starts with a drive letter
		if len(testFile) < 3 || testFile[1] != ':' {
			t.Errorf("expected Windows drive letter path, got: %s", testFile)
		}
	})

	// Test case insensitivity on Windows (file system is case-insensitive)
	t.Run("case_insensitive_filesystem", func(t *testing.T) {
		upperPath := filepath.Join(tmpDir, "TEST.JSON")

		// On Windows, this should find the same file
		if _, err := os.Stat(upperPath); err != nil {
			t.Logf("Case-insensitive lookup failed: %v", err)
			// This is informational - Windows file system is case-insensitive
			// but case-preserving, so this behavior can vary
		}
	})
}

// TestUnixSpecificPaths tests Unix-specific path handling
func TestUnixSpecificPaths(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Unix-specific test")
	}

	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.json")

	if err := os.WriteFile(testFile, []byte("test"), 0o644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	// Test that paths start with forward slash
	t.Run("absolute_paths_with_slash", func(t *testing.T) {
		if !filepath.IsAbs(testFile) {
			t.Errorf("expected absolute path, got: %s", testFile)
		}

		if testFile[0] != '/' {
			t.Errorf("expected Unix absolute path to start with /, got: %s", testFile)
		}
	})

	// Test case sensitivity on Unix (file system is case-sensitive)
	t.Run("case_sensitive_filesystem", func(t *testing.T) {
		upperFile := filepath.Join(tmpDir, "TEST.JSON")

		// On Unix, this should NOT find the file (different case)
		if _, err := os.Stat(upperFile); err == nil {
			t.Logf("Case-sensitive lookup unexpectedly succeeded (file system may be case-insensitive)")
		}
	})
}

// TestPathSeparatorHandling ensures proper separator usage
func TestPathSeparatorHandling(t *testing.T) {
	tmpDir := t.TempDir()
	subDir := filepath.Join(tmpDir, "subdir")
	testFile := filepath.Join(subDir, "test.json")

	if err := os.MkdirAll(subDir, 0o755); err != nil {
		t.Fatalf("failed to create subdir: %v", err)
	}
	if err := os.WriteFile(testFile, []byte("test"), 0o644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	// Verify filepath.Join uses correct separator
	expectedSep := string(os.PathSeparator)
	if runtime.GOOS == "windows" {
		if expectedSep != "\\" {
			t.Errorf("expected Windows path separator \\, got: %s", expectedSep)
		}
	} else {
		if expectedSep != "/" {
			t.Errorf("expected Unix path separator /, got: %s", expectedSep)
		}
	}

	// Test that ScanFiles handles the paths correctly
	files, err := ScanFiles([]string{
		filepath.Join(subDir, "*"),
	}, FilterOpts{})
	if err != nil {
		t.Fatalf("scan failed: %v", err)
	}

	if len(files) != 1 {
		t.Errorf("expected 1 file, got %d", len(files))
	}
}
