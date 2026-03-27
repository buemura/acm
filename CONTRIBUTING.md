# Contributing to ACM

Thank you for your interest in contributing to ACM (Agent Cache Manager)! This guide will help you get started.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Workflow](#development-workflow)
- [Testing](#testing)
- [Cross-Platform Considerations](#cross-platform-considerations)
- [Pull Request Process](#pull-request-process)
- [Code Style](#code-style)

## Code of Conduct

Be respectful and constructive in all interactions. We're here to build great software together.

## Getting Started

### Prerequisites

- **Go 1.24 or later** - [Download](https://go.dev/dl/)
- **Git** - [Download](https://git-scm.com/downloads)
- **A code editor** - [VS Code](https://code.visualstudio.com/) recommended with Go extension

### Fork and Clone

1. Fork the repository on GitHub
2. Clone your fork:
```bash
git clone https://github.com/YOUR_USERNAME/acm.git
cd acm
```

3. Add upstream remote:
```bash
git remote add upstream https://github.com/buemura/acm.git
```

### Build from Source

**Unix/macOS:**
```bash
make build
make test
```

**Windows (PowerShell):**
```powershell
.\build.ps1
.\build.ps1 -Test
```

**Any platform (direct Go commands):**
```bash
go build -o acm .
go test ./...
```

## Development Workflow

### 1. Create a Feature Branch

```bash
git checkout -b feat/your-feature-name
# or
git checkout -b fix/your-bug-fix
```

Branch naming conventions:
- `feat/` - New features
- `fix/` - Bug fixes
- `docs/` - Documentation changes
- `test/` - Test improvements
- `refactor/` - Code refactoring
- `chore/` - Maintenance tasks

### 2. Make Your Changes

Write clear, focused commits:
```bash
git add file1.go file2.go
git commit -m "feat: add support for custom cache paths"
```

Commit message format:
```
type: brief description

Optional longer explanation if needed.

Fixes #123
```

Types: `feat`, `fix`, `docs`, `test`, `refactor`, `chore`

### 3. Keep Your Branch Updated

```bash
git fetch upstream
git rebase upstream/master
```

### 4. Run Tests

Before pushing:
```bash
go test ./...
go build .
./acm doctor  # Test the binary
```

### 5. Push and Create PR

```bash
git push origin feat/your-feature-name
```

Then create a Pull Request on GitHub.

## Testing

### Running Tests

**All tests:**
```bash
go test ./...
```

**With verbose output:**
```bash
go test -v ./...
```

**Specific package:**
```bash
go test ./internal/provider
go test ./internal/scanner
```

**With coverage:**
```bash
go test -cover ./...
```

### Writing Tests

#### Test File Naming

- Test files: `*_test.go`
- Platform-specific tests: `*_windows_test.go`, `*_unix_test.go`
- Cross-platform tests: `*_crossplatform_test.go`

#### Test Function Naming

```go
func TestFunctionName(t *testing.T) {
    // Test implementation
}

func TestFunctionName_EdgeCase(t *testing.T) {
    // Test specific edge case
}
```

#### Example Test

```go
func TestScanFiles(t *testing.T) {
    tmpDir := t.TempDir()
    testFile := filepath.Join(tmpDir, "test.json")

    if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
        t.Fatalf("failed to create test file: %v", err)
    }

    files, err := scanner.ScanFiles([]string{testFile}, scanner.FilterOpts{})
    if err != nil {
        t.Fatalf("scan failed: %v", err)
    }

    if len(files) != 1 {
        t.Errorf("expected 1 file, got %d", len(files))
    }
}
```

### Platform-Specific Tests

Use build tags for platform-specific tests:

**Windows-only test:**
```go
//go:build windows

package provider

import "testing"

func TestWindowsSpecificFeature(t *testing.T) {
    // Windows-specific test
}
```

**Unix-only test:**
```go
//go:build !windows

package provider

import "testing"

func TestUnixSpecificFeature(t *testing.T) {
    // Unix/macOS-specific test
}
```

**Runtime detection:**
```go
func TestCrossPlatform(t *testing.T) {
    if runtime.GOOS == "windows" {
        // Windows-specific assertions
    } else {
        // Unix/macOS assertions
    }
}
```

### Testing on Multiple Platforms

**Local testing:**
- Test on your development platform
- Use VM or Docker for other platforms if possible

**CI testing:**
- All PRs are automatically tested on Windows, macOS, and Linux
- Check the CI results before merging

## Cross-Platform Considerations

ACM must work on Windows, macOS, and Linux. Follow these guidelines:

### Path Handling

**✅ DO:**
```go
// Use filepath.Join for all paths
path := filepath.Join(home, ".claude", "cache")

// Use os.PathSeparator for dynamic separators
sep := string(os.PathSeparator)

// Use filepath.Clean to normalize paths
clean := filepath.Clean(path)
```

**❌ DON'T:**
```go
// Don't hardcode path separators
path := home + "/.claude/cache"  // Breaks on Windows

// Don't assume forward slashes
path := strings.Split(p, "/")    // Breaks on Windows
```

### Home Directory

**✅ DO:**
```go
home, err := os.UserHomeDir()
if err != nil {
    return err
}
```

**❌ DON'T:**
```go
home := os.Getenv("HOME")  // Doesn't work on Windows
home := "/home/user"        // Hardcoded path
```

### Platform Detection

```go
import "runtime"

if runtime.GOOS == "windows" {
    // Windows-specific code
} else {
    // Unix/macOS code
}
```

### Environment Variables

**Windows:**
- `USERPROFILE` - Home directory
- `APPDATA` - Roaming app data
- `LOCALAPPDATA` - Local app data

**Unix/macOS:**
- `HOME` - Home directory
- `XDG_CACHE_HOME` - Cache directory (optional)

**Cross-platform:**
```go
home, _ := os.UserHomeDir()  // Works everywhere
```

### File Permissions

Windows doesn't use Unix permission bits. Use sensible defaults:

```go
// Create file (0644 on Unix, ignored on Windows)
os.WriteFile(path, data, 0644)

// Create directory (0755 on Unix, ignored on Windows)
os.MkdirAll(path, 0755)
```

### Line Endings

- Git should handle line ending conversion automatically
- Tests should not depend on specific line endings
- Use `strings.TrimSpace()` when comparing output

### Commands and Scripts

Provide alternatives for platform-specific tools:

- Makefile (Unix/macOS)
- PowerShell scripts (Windows)
- Direct Go commands (all platforms)

## Pull Request Process

### Before Submitting

1. ✅ All tests pass: `go test ./...`
2. ✅ Code builds: `go build .`
3. ✅ No linting errors: `go vet ./...`
4. ✅ Health check works: `./acm doctor`
5. ✅ Documentation updated (if needed)
6. ✅ Commits are clear and focused

### PR Description Template

```markdown
## Description

Brief description of what this PR does.

## Changes

- Added X
- Fixed Y
- Updated Z

## Testing

- [ ] Tested on Windows
- [ ] Tested on macOS
- [ ] Tested on Linux
- [ ] Added/updated tests
- [ ] All tests pass

## Checklist

- [ ] Code follows project conventions
- [ ] Tests added for new functionality
- [ ] Documentation updated
- [ ] No breaking changes (or documented)

Fixes #123
```

### Review Process

1. Maintainers will review your PR
2. Address any feedback
3. Once approved, your PR will be merged
4. Your contribution will be included in the next release!

## Code Style

### General Go Style

Follow standard Go conventions:
- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://go.dev/wiki/CodeReviewComments)

### ACM-Specific Conventions

**Error messages:**
```go
// Lowercase, no punctuation
return fmt.Errorf("failed to read file: %w", err)
```

**Function comments:**
```go
// ScanFiles discovers files matching the given patterns.
// It returns a deduplicated list of FileInfo sorted by path.
func ScanFiles(patterns []string, opts FilterOpts) ([]FileInfo, error) {
    // Implementation
}
```

**Package comments:**
```go
// Package scanner provides file discovery and filtering functionality
// for cache management operations.
package scanner
```

**Variable naming:**
```go
// Short names for common types
p := provider.Get("claude")  // provider
f := scanner.FileInfo{}      // file info
err := doSomething()         // error

// Descriptive names for less obvious types
cachePath := filepath.Join(home, ".claude")
fileCount := len(files)
```

### Formatting

Run `gofmt` before committing:
```bash
gofmt -w .
```

Or use your editor's auto-format feature.

## Project Structure

```
acm/
├── main.go                 # Entry point
├── cmd/cli/                # CLI commands
│   ├── root.go            # Root command
│   ├── list.go            # List command
│   ├── clean.go           # Clean command
│   ├── provider.go        # Provider command
│   ├── doctor.go          # Doctor command
│   └── verbose.go         # Verbose mode helpers
├── internal/
│   ├── provider/          # Provider implementations
│   │   ├── provider.go    # Provider interface
│   │   ├── paths.go       # Cross-platform path helpers
│   │   ├── claude.go      # Claude Code provider
│   │   └── codex.go       # Codex provider
│   ├── scanner/           # File scanning
│   │   └── scanner.go     # File discovery and filtering
│   └── ui/                # User interface
│       ├── format.go      # Output formatting
│       └── menu.go        # Interactive prompts
├── docs/                  # Documentation
│   ├── windows.md         # Windows guide
│   └── troubleshooting.md # Troubleshooting guide
├── build.ps1              # Windows build script
├── install.ps1            # Windows install script
├── Makefile               # Unix/macOS build automation
├── go.mod                 # Go module definition
├── README.md              # Project documentation
└── CONTRIBUTING.md        # This file
```

## Adding a New Provider

To add support for a new AI assistant:

1. Create `internal/provider/newprovider.go`:
```go
package provider

import (
    "os"
    "path/filepath"
)

func init() {
    home, _ := os.UserHomeDir()
    providerDir := buildAppDir("newprovider")

    Register(&Provider{
        Name: "newprovider",
        CachePaths: getAllCachePaths("newprovider"),
        Checks: []PathCheck{
            {
                Name: "cache files",
                Patterns: getAllCachePaths("newprovider"),
            },
        },
    })
}
```

2. Add tests in `internal/provider/newprovider_test.go`

3. Update README.md with the new provider

4. Test on all platforms

## Questions?

- Open an issue: https://github.com/buemura/acm/issues
- Discussion: https://github.com/buemura/acm/discussions

Thank you for contributing! 🎉
