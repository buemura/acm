package provider

import (
	"os"
	"path/filepath"
	"runtime"
)

// buildCachePaths returns cache paths for a given app, including
// platform-specific locations (APPDATA on Windows, ~/.cache on Unix)
func buildCachePaths(appName string, subPaths ...string) []string {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}

	paths := make([]string, 0)

	// Unix-style paths (work on all platforms including Windows)
	appDir := filepath.Join(home, "."+appName)
	paths = append(paths, filepath.Join(append([]string{appDir}, subPaths...)...))

	// Alternative cache location (~/.cache/appName)
	if len(subPaths) > 0 && subPaths[0] == "cache" {
		paths = append(paths, filepath.Join(home, ".cache", appName))
	}

	// Windows-specific paths
	if runtime.GOOS == "windows" {
		// APPDATA: C:\Users\username\AppData\Roaming
		if appData := os.Getenv("APPDATA"); appData != "" {
			paths = append(paths, filepath.Join(append([]string{appData, appName}, subPaths...)...))
		}

		// LOCALAPPDATA: C:\Users\username\AppData\Local
		if localAppData := os.Getenv("LOCALAPPDATA"); localAppData != "" {
			// Use capitalized app name for Windows convention
			capitalizedName := capitalize(appName)
			paths = append(paths, filepath.Join(append([]string{localAppData, capitalizedName}, subPaths...)...))
		}
	}

	return paths
}

// buildAppDir returns the base directory for an app (~/.appname)
func buildAppDir(appName string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, "."+appName)
}

// capitalize returns the string with the first letter capitalized
func capitalize(s string) string {
	if s == "" {
		return ""
	}
	// Convert first character to uppercase if it's lowercase
	first := s[0]
	if first >= 'a' && first <= 'z' {
		first = first - 32
	}
	if len(s) == 1 {
		return string(first)
	}
	return string(first) + s[1:]
}

// getAllCachePaths returns all possible cache directory paths for an app
func getAllCachePaths(appName string) []string {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}

	paths := make([]string, 0)

	// Primary Unix-style cache
	appDir := filepath.Join(home, "."+appName)
	paths = append(paths, filepath.Join(appDir, "cache"))

	// Alternative Unix cache location
	paths = append(paths, filepath.Join(home, ".cache", appName))

	// Windows-specific paths
	if runtime.GOOS == "windows" {
		if appData := os.Getenv("APPDATA"); appData != "" {
			paths = append(paths, filepath.Join(appData, appName, "cache"))
		}
		if localAppData := os.Getenv("LOCALAPPDATA"); localAppData != "" {
			capitalizedName := capitalize(appName)
			paths = append(paths, filepath.Join(localAppData, capitalizedName, "Cache"))
		}
	}

	return paths
}
