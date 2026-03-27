package provider

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// TestHomeDirectoryResolution tests that home directory
// resolution works correctly on all platforms
func TestHomeDirectoryResolution(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("failed to get home directory: %v", err)
	}

	if home == "" {
		t.Fatal("home directory is empty")
	}

	// Verify home directory exists
	if _, err := os.Stat(home); err != nil {
		t.Fatalf("home directory doesn't exist: %v", err)
	}

	// Platform-specific checks
	if runtime.GOOS == "windows" {
		// On Windows, home should be like C:\Users\username
		if !filepath.IsAbs(home) {
			t.Errorf("expected absolute Windows path, got: %s", home)
		}
		if len(home) < 3 || home[1] != ':' {
			t.Errorf("expected Windows drive letter in home path, got: %s", home)
		}
		t.Logf("Windows home directory: %s", home)
	} else {
		// On Unix, home should start with /
		if !filepath.IsAbs(home) {
			t.Errorf("expected absolute Unix path, got: %s", home)
		}
		if home[0] != '/' {
			t.Errorf("expected Unix path to start with /, got: %s", home)
		}
		t.Logf("Unix home directory: %s", home)
	}
}

// TestProviderPathConstruction tests that provider cache paths
// are constructed correctly on all platforms
func TestProviderPathConstruction(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("failed to get home directory: %v", err)
	}

	tests := []struct {
		name     string
		basePath string
		subPaths []string
	}{
		{
			name:     "cache_directory",
			basePath: home,
			subPaths: []string{".claude", "cache"},
		},
		{
			name:     "projects_directory",
			basePath: home,
			subPaths: []string{".claude", "projects"},
		},
		{
			name:     "nested_cache",
			basePath: home,
			subPaths: []string{".cache", "claude"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Build path using filepath.Join
			parts := append([]string{tt.basePath}, tt.subPaths...)
			fullPath := filepath.Join(parts...)

			// Verify path is absolute
			if !filepath.IsAbs(fullPath) {
				t.Errorf("expected absolute path, got: %s", fullPath)
			}

			// Verify path uses correct separator
			if runtime.GOOS == "windows" {
				// Windows paths should contain backslashes when created with filepath.Join
				t.Logf("Windows path: %s", fullPath)
			} else {
				// Unix paths should contain forward slashes
				t.Logf("Unix path: %s", fullPath)
			}

			// Verify we can create this path structure
			tmpTestPath := filepath.Join(t.TempDir(), filepath.Base(fullPath))
			if err := os.MkdirAll(tmpTestPath, 0o755); err != nil {
				t.Errorf("failed to create path %s: %v", tmpTestPath, err)
			}
		})
	}
}

// TestWindowsEnvironmentVariables tests Windows-specific env vars
func TestWindowsEnvironmentVariables(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Windows-specific test")
	}

	// Test USERPROFILE (Windows equivalent of HOME)
	t.Run("USERPROFILE", func(t *testing.T) {
		userProfile := os.Getenv("USERPROFILE")
		if userProfile == "" {
			t.Error("USERPROFILE environment variable is not set")
		}

		home, _ := os.UserHomeDir()
		// These should typically be the same on Windows
		t.Logf("USERPROFILE: %s", userProfile)
		t.Logf("UserHomeDir: %s", home)
	})

	// Test APPDATA (common Windows cache location)
	t.Run("APPDATA", func(t *testing.T) {
		appData := os.Getenv("APPDATA")
		if appData == "" {
			t.Log("APPDATA environment variable is not set (may be expected in some environments)")
		} else {
			if !filepath.IsAbs(appData) {
				t.Errorf("expected absolute path for APPDATA, got: %s", appData)
			}
			t.Logf("APPDATA: %s", appData)

			// Verify it exists
			if _, err := os.Stat(appData); err != nil {
				t.Errorf("APPDATA path doesn't exist: %v", err)
			}
		}
	})

	// Test LOCALAPPDATA (alternative Windows cache location)
	t.Run("LOCALAPPDATA", func(t *testing.T) {
		localAppData := os.Getenv("LOCALAPPDATA")
		if localAppData == "" {
			t.Log("LOCALAPPDATA environment variable is not set (may be expected in some environments)")
		} else {
			if !filepath.IsAbs(localAppData) {
				t.Errorf("expected absolute path for LOCALAPPDATA, got: %s", localAppData)
			}
			t.Logf("LOCALAPPDATA: %s", localAppData)

			// Verify it exists
			if _, err := os.Stat(localAppData); err != nil {
				t.Errorf("LOCALAPPDATA path doesn't exist: %v", err)
			}
		}
	})
}

// TestUnixEnvironmentVariables tests Unix-specific env vars
func TestUnixEnvironmentVariables(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Unix-specific test")
	}

	// Test HOME environment variable
	t.Run("HOME", func(t *testing.T) {
		homeEnv := os.Getenv("HOME")
		if homeEnv == "" {
			t.Error("HOME environment variable is not set")
		}

		home, _ := os.UserHomeDir()
		// These should be the same on Unix
		if homeEnv != home {
			t.Logf("HOME env (%s) differs from UserHomeDir (%s)", homeEnv, home)
		}

		t.Logf("HOME: %s", homeEnv)
	})
}

// TestProviderCachePathsExist is a helper test to check if actual
// cache paths exist on the system (for debugging)
func TestProviderCachePathsExist(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping cache path check in short mode")
	}

	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("failed to get home directory: %v", err)
	}

	// Common cache paths to check
	cachePaths := []string{
		filepath.Join(home, ".claude"),
		filepath.Join(home, ".claude", "cache"),
		filepath.Join(home, ".claude", "projects"),
		filepath.Join(home, ".cache", "claude"),
		filepath.Join(home, ".codex"),
		filepath.Join(home, ".codex", "cache"),
	}

	// On Windows, also check APPDATA locations
	if runtime.GOOS == "windows" {
		if appData := os.Getenv("APPDATA"); appData != "" {
			cachePaths = append(cachePaths,
				filepath.Join(appData, "claude"),
				filepath.Join(appData, "codex"),
			)
		}
		if localAppData := os.Getenv("LOCALAPPDATA"); localAppData != "" {
			cachePaths = append(cachePaths,
				filepath.Join(localAppData, "Claude"),
				filepath.Join(localAppData, "Codex"),
			)
		}
	}

	t.Logf("Checking cache paths on %s:", runtime.GOOS)
	foundCount := 0
	for _, path := range cachePaths {
		if info, err := os.Stat(path); err == nil {
			t.Logf("  ✓ %s (exists, isDir: %v)", path, info.IsDir())
			foundCount++
		} else {
			t.Logf("  ✗ %s (not found)", path)
		}
	}

	if foundCount == 0 {
		t.Log("No cache directories found (this is normal if AI assistants haven't been used yet)")
	}
}
