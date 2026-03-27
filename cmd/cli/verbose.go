package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/buemura/acm/internal/provider"
)

// printVerbosePathInfo prints detailed information about cache paths
func printVerbosePathInfo(p *provider.Provider) {
	fmt.Fprintf(os.Stderr, "\n=== Checking cache locations for %s ===\n\n", p.Name)

	// Track unique directories we've checked
	checkedDirs := make(map[string]bool)

	for _, path := range p.CachePaths {
		// Skip glob patterns with wildcards for display purposes
		if strings.Contains(path, "*") {
			continue
		}

		// Only show each directory once
		if checkedDirs[path] {
			continue
		}
		checkedDirs[path] = true

		// Check if path exists
		info, err := os.Stat(path)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "  ✗ %s (not found)\n", path)
			} else {
				fmt.Fprintf(os.Stderr, "  ✗ %s (error: %v)\n", path, err)
			}
			continue
		}

		// Path exists - show it with details
		if info.IsDir() {
			// Count files in directory
			count := 0
			filepath.Walk(path, func(p string, fi os.FileInfo, err error) error {
				if err == nil && !fi.IsDir() {
					count++
				}
				return nil
			})
			fmt.Fprintf(os.Stderr, "  ✓ %s (directory, %d files)\n", path, count)
		} else {
			// It's a file
			fmt.Fprintf(os.Stderr, "  ✓ %s (file, %d bytes)\n", path, info.Size())
		}
	}

	fmt.Fprintf(os.Stderr, "\n")
}

// printVerboseCheckInfo prints information about path checks
func printVerboseCheckInfo(p *provider.Provider) {
	fmt.Fprintf(os.Stderr, "=== Path check results ===\n\n")

	for _, check := range p.Checks {
		found := false
		var foundPaths []string

		for _, pattern := range check.Patterns {
			matches, err := filepath.Glob(pattern)
			if err != nil {
				continue
			}
			if len(matches) > 0 {
				found = true
				foundPaths = append(foundPaths, matches...)
			}
		}

		if found {
			fmt.Fprintf(os.Stderr, "  ✓ %s: found %d match(es)\n", check.Name, len(foundPaths))
		} else {
			fmt.Fprintf(os.Stderr, "  ✗ %s: not found\n", check.Name)
		}
	}

	fmt.Fprintf(os.Stderr, "\n")
}
