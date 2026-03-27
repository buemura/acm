package provider

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestBuildAppDir(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("failed to get home directory: %v", err)
	}

	tests := []struct {
		name    string
		appName string
		want    string
	}{
		{
			name:    "claude",
			appName: "claude",
			want:    filepath.Join(home, ".claude"),
		},
		{
			name:    "codex",
			appName: "codex",
			want:    filepath.Join(home, ".codex"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildAppDir(tt.appName)
			if got != tt.want {
				t.Errorf("buildAppDir(%q) = %q, want %q", tt.appName, got, tt.want)
			}
		})
	}
}

func TestGetAllCachePaths(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("failed to get home directory: %v", err)
	}

	tests := []struct {
		name         string
		appName      string
		wantContains []string
		wantMinCount int
	}{
		{
			name:    "claude_paths",
			appName: "claude",
			wantContains: []string{
				filepath.Join(home, ".claude", "cache"),
				filepath.Join(home, ".cache", "claude"),
			},
			wantMinCount: 2, // At least Unix paths
		},
		{
			name:    "codex_paths",
			appName: "codex",
			wantContains: []string{
				filepath.Join(home, ".codex", "cache"),
				filepath.Join(home, ".cache", "codex"),
			},
			wantMinCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getAllCachePaths(tt.appName)

			// Verify minimum count
			if len(got) < tt.wantMinCount {
				t.Errorf("getAllCachePaths(%q) returned %d paths, want at least %d", tt.appName, len(got), tt.wantMinCount)
			}

			// Verify expected paths are present
			for _, wantPath := range tt.wantContains {
				found := false
				for _, gotPath := range got {
					if gotPath == wantPath {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("getAllCachePaths(%q) missing expected path: %s\nGot: %v", tt.appName, wantPath, got)
				}
			}

			// Log all paths for debugging
			t.Logf("Cache paths for %s:", tt.appName)
			for _, p := range got {
				t.Logf("  - %s", p)
			}
		})
	}
}

func TestGetAllCachePathsWindows(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Windows-specific test")
	}

	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("failed to get home directory: %v", err)
	}

	appName := "claude"
	paths := getAllCachePaths(appName)

	// On Windows, should have at least 4 paths:
	// 1. ~/.claude/cache
	// 2. ~/.cache/claude
	// 3. %APPDATA%/claude/cache (if APPDATA is set)
	// 4. %LOCALAPPDATA%/Claude/Cache (if LOCALAPPDATA is set)
	minExpected := 2
	if os.Getenv("APPDATA") != "" {
		minExpected++
	}
	if os.Getenv("LOCALAPPDATA") != "" {
		minExpected++
	}

	if len(paths) < minExpected {
		t.Errorf("expected at least %d cache paths on Windows, got %d", minExpected, len(paths))
	}

	// Verify Unix-style paths are included
	unixStylePath := filepath.Join(home, ".claude", "cache")
	found := false
	for _, p := range paths {
		if p == unixStylePath {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Unix-style path %s not found in Windows cache paths", unixStylePath)
	}

	// Verify Windows-specific paths are included if env vars are set
	if appData := os.Getenv("APPDATA"); appData != "" {
		appDataPath := filepath.Join(appData, "claude", "cache")
		found := false
		for _, p := range paths {
			if p == appDataPath {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("APPDATA path %s not found in Windows cache paths", appDataPath)
		}
	}

	if localAppData := os.Getenv("LOCALAPPDATA"); localAppData != "" {
		localAppDataPath := filepath.Join(localAppData, "Claude", "Cache")
		found := false
		for _, p := range paths {
			if p == localAppDataPath {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("LOCALAPPDATA path %s not found in Windows cache paths", localAppDataPath)
		}
	}

	t.Logf("Windows cache paths for %s:", appName)
	for _, p := range paths {
		t.Logf("  - %s", p)
	}
}

func TestCapitalize(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "lowercase",
			input: "claude",
			want:  "Claude",
		},
		{
			name:  "already_capitalized",
			input: "Codex",
			want:  "Codex",
		},
		{
			name:  "empty",
			input: "",
			want:  "",
		},
		{
			name:  "single_char",
			input: "a",
			want:  "A",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := capitalize(tt.input)
			if got != tt.want {
				t.Errorf("capitalize(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestBuildCachePaths(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("failed to get home directory: %v", err)
	}

	tests := []struct {
		name         string
		appName      string
		subPaths     []string
		wantContains string
	}{
		{
			name:         "cache_subpath",
			appName:      "claude",
			subPaths:     []string{"cache"},
			wantContains: filepath.Join(home, ".claude", "cache"),
		},
		{
			name:         "nested_subpath",
			appName:      "claude",
			subPaths:     []string{"projects", "test"},
			wantContains: filepath.Join(home, ".claude", "projects", "test"),
		},
		{
			name:         "no_subpath",
			appName:      "claude",
			subPaths:     []string{},
			wantContains: filepath.Join(home, ".claude"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := buildCachePaths(tt.appName, tt.subPaths...)
			if len(got) == 0 {
				t.Errorf("buildCachePaths(%q, %v) returned empty slice", tt.appName, tt.subPaths)
				return
			}

			// Check if expected path is in results
			found := false
			for _, p := range got {
				if p == tt.wantContains {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("buildCachePaths(%q, %v) missing expected path: %s\nGot: %v", tt.appName, tt.subPaths, tt.wantContains, got)
			}

			// On Windows, should have more paths
			if runtime.GOOS == "windows" && len(got) < 2 {
				t.Errorf("on Windows, expected multiple paths, got: %v", got)
			}
		})
	}
}

func TestProviderPathsNotEmpty(t *testing.T) {
	// This test ensures that the provider initialization works
	// and produces non-empty paths
	providers := []string{"claude", "codex"}

	for _, name := range providers {
		t.Run(name, func(t *testing.T) {
			p, err := Get(name)
			if err != nil {
				t.Fatalf("provider %q not registered: %v", name, err)
			}
			if p == nil {
				t.Fatalf("provider %q is nil", name)
			}

			if len(p.CachePaths) == 0 {
				t.Errorf("provider %q has no cache paths", name)
			}

			if len(p.Checks) == 0 {
				t.Errorf("provider %q has no checks", name)
			}

			// Verify paths are absolute
			for _, path := range p.CachePaths {
				// Skip glob patterns
				if strings.Contains(path, "*") {
					continue
				}
				if !filepath.IsAbs(path) {
					t.Errorf("provider %q has relative cache path: %s", name, path)
				}
			}

			t.Logf("Provider %q has %d cache paths", name, len(p.CachePaths))
		})
	}
}
