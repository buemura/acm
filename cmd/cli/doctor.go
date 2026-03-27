package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/buemura/acm/internal/provider"
	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check installation and diagnose issues",
	Long:  "Run diagnostics to verify ACM installation, check provider paths, and identify potential issues",
	Run:   runDoctor,
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}

func runDoctor(cmd *cobra.Command, args []string) {
	fmt.Println("🔍 Running ACM diagnostics...\n")

	checks := []struct {
		name string
		fn   func() (bool, string, string)
	}{
		{"Go Installation", checkGo},
		{"ACM Binary", checkBinary},
		{"PATH Configuration", checkPath},
		{"Provider Detection", checkProviders},
		{"File Permissions", checkPermissions},
	}

	passed := 0
	failed := 0
	warnings := 0

	for _, check := range checks {
		fmt.Printf("Checking %s...\n", check.name)
		ok, message, suggestion := check.fn()

		if ok {
			fmt.Printf("  ✓ %s\n", message)
			passed++
		} else if strings.Contains(message, "warning") || strings.Contains(message, "missing") {
			fmt.Printf("  ⚠ %s\n", message)
			if suggestion != "" {
				fmt.Printf("    %s\n", suggestion)
			}
			warnings++
		} else {
			fmt.Printf("  ✗ %s\n", message)
			if suggestion != "" {
				fmt.Printf("    Fix: %s\n", suggestion)
			}
			failed++
		}
		fmt.Println()
	}

	// Summary
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Printf("Summary: %d passed, %d warnings, %d failed\n", passed, warnings, failed)

	if failed == 0 && warnings == 0 {
		fmt.Println("\n✓ Everything looks good! ACM is ready to use.")
	} else if failed == 0 {
		fmt.Println("\n⚠ ACM is functional but some optional features may not work.")
	} else {
		fmt.Println("\n✗ Some critical checks failed. Please fix the issues above.")
		os.Exit(1)
	}
}

func checkGo() (bool, string, string) {
	cmd := exec.Command("go", "version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false, "Go is not installed or not in PATH", "Install Go from https://go.dev/dl/"
	}

	version := strings.TrimSpace(string(output))
	return true, fmt.Sprintf("Go found: %s", version), ""
}

func checkBinary() (bool, string, string) {
	// Check if running from go run or actual binary
	executable, err := os.Executable()
	if err != nil {
		return false, "Cannot determine executable path", ""
	}

	// Check if it's a proper installation
	if strings.Contains(executable, "go-build") {
		return true, "Running from 'go run' (development mode)", ""
	}

	// Check file size (compiled binary should be > 1MB)
	info, err := os.Stat(executable)
	if err != nil {
		return false, "Cannot stat executable", ""
	}

	sizeMB := float64(info.Size()) / (1024 * 1024)
	return true, fmt.Sprintf("Binary found at %s (%.2f MB)", executable, sizeMB), ""
}

func checkPath() (bool, string, string) {
	// Check if we're running from development or installed location
	executable, _ := os.Executable()
	cwd, _ := os.Getwd()
	execDir := filepath.Dir(executable)

	// Check if running from current directory (dev mode)
	isDevMode := strings.Contains(executable, "go-build") || execDir == cwd

	// If in dev mode, just check if Go bin exists
	if isDevMode {
		goBinPath := getGoBinPath()
		if _, err := os.Stat(goBinPath); err != nil {
			return true, "Running in dev mode (Go bin directory will be created on install)", ""
		}

		// Check if installed version exists
		expectedPath := filepath.Join(goBinPath, "acm")
		if runtime.GOOS == "windows" {
			expectedPath += ".exe"
		}
		if _, err := os.Stat(expectedPath); err == nil {
			return true, fmt.Sprintf("Running from current directory (installed version available at %s)", goBinPath), ""
		}

		return true, fmt.Sprintf("Running in dev mode (install to %s for PATH access)", goBinPath), ""
	}

	// Check if acm is in PATH
	acmPath, err := exec.LookPath("acm")
	if err != nil {
		goBinPath := getGoBinPath()
		suggestion := ""
		if runtime.GOOS == "windows" {
			suggestion = fmt.Sprintf("Add %s to PATH:\n      [Environment]::SetEnvironmentVariable(\"Path\", $env:Path + \";%s\", \"User\")", goBinPath, goBinPath)
		} else {
			suggestion = fmt.Sprintf("Add %s to PATH:\n      export PATH=\"$PATH:%s\"", goBinPath, goBinPath)
		}

		// Check if binary exists in Go bin
		expectedPath := filepath.Join(goBinPath, "acm")
		if runtime.GOOS == "windows" {
			expectedPath += ".exe"
		}
		if _, err := os.Stat(expectedPath); err == nil {
			return false, "ACM installed but not in PATH", suggestion
		}

		return false, "ACM not found in PATH", suggestion + "\n      Or run: go install github.com/buemura/acm@latest"
	}

	// Check if Go bin is in PATH
	goBinPath := getGoBinPath()
	pathEnv := os.Getenv("PATH")
	paths := filepath.SplitList(pathEnv)

	inPath := false
	for _, p := range paths {
		if strings.TrimRight(p, string(os.PathSeparator)) == strings.TrimRight(goBinPath, string(os.PathSeparator)) {
			inPath = true
			break
		}
	}

	if !inPath {
		return true, fmt.Sprintf("ACM found at %s (not from Go bin)", acmPath), ""
	}

	return true, fmt.Sprintf("ACM is in PATH (%s)", goBinPath), ""
}

func checkProviders() (bool, string, string) {
	providers := provider.List()
	if len(providers) == 0 {
		return false, "No providers registered", "This should not happen - provider registration may be broken"
	}

	foundAny := false
	details := []string{}

	for _, p := range providers {
		missing := provider.MissingPathChecks(p)
		if len(missing) == len(p.Checks) {
			details = append(details, fmt.Sprintf("    - %s: not found (is %s installed?)", p.Name, p.Name))
		} else if len(missing) > 0 {
			foundAny = true
			details = append(details, fmt.Sprintf("    - %s: found (missing: %s)", p.Name, strings.Join(missing, ", ")))
		} else {
			foundAny = true
			details = append(details, fmt.Sprintf("    - %s: found (all paths present)", p.Name))
		}
	}

	message := fmt.Sprintf("%d providers registered:\n%s", len(providers), strings.Join(details, "\n"))

	if !foundAny {
		return false,
			fmt.Sprintf("No provider cache found\n%s", message),
			"Have you used Claude Code or Codex yet? Caches are created after first use."
	}

	return true, message, ""
}

func checkPermissions() (bool, string, string) {
	home, err := os.UserHomeDir()
	if err != nil {
		return false, "Cannot determine home directory", ""
	}

	// Check if we can read/write to home directory
	testFile := filepath.Join(home, ".acm-test-permissions")
	err = os.WriteFile(testFile, []byte("test"), 0644)
	if err != nil {
		return false,
			fmt.Sprintf("Cannot write to home directory: %v", err),
			"Check directory permissions"
	}
	os.Remove(testFile)

	// Check provider directories
	providers := provider.List()
	readableCount := 0
	totalChecked := 0

	for _, p := range providers {
		for _, path := range p.CachePaths {
			// Skip glob patterns
			if strings.Contains(path, "*") {
				continue
			}

			totalChecked++
			_, err := os.Stat(path)
			if err == nil {
				// Try to read directory
				if info, err := os.Stat(path); err == nil && info.IsDir() {
					if _, err := os.ReadDir(path); err == nil {
						readableCount++
					}
				} else if err == nil {
					readableCount++
				}
			}
		}
	}

	if totalChecked == 0 {
		return true, "No cache directories to check", ""
	}

	if readableCount == totalChecked {
		return true, fmt.Sprintf("All cache directories readable (%d/%d)", readableCount, totalChecked), ""
	}

	return true,
		fmt.Sprintf("Some directories not accessible: %d/%d readable", readableCount, totalChecked),
		""
}

func getGoBinPath() string {
	goPath := os.Getenv("GOPATH")
	if goPath == "" {
		home, _ := os.UserHomeDir()
		goPath = filepath.Join(home, "go")
	}
	return filepath.Join(goPath, "bin")
}
