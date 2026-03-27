# Windows Installation and Usage Guide

Complete guide for installing and using ACM on Windows.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Installation Methods](#installation-methods)
- [Post-Installation Setup](#post-installation-setup)
- [Building from Source](#building-from-source)
- [Common Issues](#common-issues)
- [Windows-Specific Features](#windows-specific-features)

## Prerequisites

### Required

- **Windows 10 or later** (Windows 11 recommended)
- **Go 1.24 or later** - Download from [go.dev/dl](https://go.dev/dl/)

### Optional

- **PowerShell 5.1 or later** (for build scripts)
- **Git for Windows** (for building from source)

## Installation Methods

### Method 1: Via `go install` (Recommended)

This is the simplest method and works on all platforms.

**Step 1: Install ACM**

```powershell
go install github.com/buemura/acm@latest
```

This installs ACM to `%USERPROFILE%\go\bin\acm.exe` (typically `C:\Users\YourName\go\bin\acm.exe`).

**Step 2: Verify Installation**

```powershell
# Check if ACM is accessible
& "$env:USERPROFILE\go\bin\acm.exe" --version

# Or run the health check
& "$env:USERPROFILE\go\bin\acm.exe" doctor
```

**Step 3: Add to PATH (if needed)**

If the command `acm` doesn't work, you need to add Go's bin directory to your PATH:

```powershell
# Run this in PowerShell (Administrator)
[Environment]::SetEnvironmentVariable(
    "Path",
    [Environment]::GetEnvironmentVariable("Path", "User") + ";$env:USERPROFILE\go\bin",
    "User"
)

# Restart your terminal, then verify:
acm --version
```

### Method 2: Download from Releases

**Step 1: Download Binary**

Go to [Releases](https://github.com/buemura/acm/releases) and download `acm-windows-amd64.exe`.

**Step 2: Rename and Move**

```powershell
# Rename the file
Rename-Item acm-windows-amd64.exe acm.exe

# Move to Go bin directory (or any directory in PATH)
Move-Item acm.exe "$env:USERPROFILE\go\bin\acm.exe"
```

**Step 3: Verify**

```powershell
acm --version
```

### Method 3: Build from Source

See [Building from Source](#building-from-source) section below.

## Post-Installation Setup

### Verify Installation

Run the built-in health check:

```powershell
acm doctor
```

This will check:
- ✓ Go installation
- ✓ ACM binary location
- ✓ PATH configuration
- ✓ Provider detection (Claude, Codex)
- ✓ File permissions

Expected output:
```
🔍 Running ACM diagnostics...

Checking Go Installation...
  ✓ Go found: go version go1.24.0 windows/amd64

Checking ACM Binary...
  ✓ Binary found at C:\Users\YourName\go\bin\acm.exe (3.9 MB)

Checking PATH Configuration...
  ✓ ACM is in PATH (C:\Users\YourName\go\bin)

Checking Provider Detection...
  ✓ 2 providers registered:
    - claude: found (missing: stats cache)
    - codex: not found (is codex installed?)

Checking File Permissions...
  ✓ All cache directories readable (4/4)

━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
Summary: 5 passed, 0 warnings, 0 failed

✓ Everything looks good! ACM is ready to use.
```

### Check Cache Locations

Use verbose mode to see where ACM looks for cache files:

```powershell
acm list -v
```

Output shows:
```
=== Checking cache locations for claude ===

  ✓ C:\Users\YourName\.claude\cache (directory, 5 files)
  ✓ C:\Users\YourName\.cache\claude (directory, 0 files)
  ✗ C:\Users\YourName\AppData\Roaming\claude\cache (not found)
  ✓ C:\Users\YourName\AppData\Local\Claude\Cache (directory, 10 files)

=== Path check results ===

  ✓ cache files: found 2 match(es)
  ✓ session transcripts: found 15 match(es)
```

## Building from Source

### Using PowerShell Scripts (Recommended)

ACM includes PowerShell scripts that make building on Windows easy.

**Step 1: Clone Repository**

```powershell
git clone https://github.com/buemura/acm.git
cd acm
```

**Step 2: Build**

```powershell
# Build debug version
.\build.ps1

# Build optimized release version
.\build.ps1 -Release

# Run tests
.\build.ps1 -Test

# Full pipeline (clean, test, build)
.\build.ps1 -All
```

**Step 3: Install**

```powershell
# Install to Go bin directory
.\install.ps1

# Check installation status
.\install.ps1 -CheckOnly

# Force reinstall
.\install.ps1 -Force
```

The install script will:
- Copy binary to `%USERPROFILE%\go\bin`
- Check if Go bin is in PATH
- Provide instructions if PATH needs updating
- Verify the installation

### Using Go Commands Directly

If you prefer not to use PowerShell scripts:

```powershell
# Build
go build -o acm.exe .

# Build optimized
go build -ldflags "-s -w" -o acm.exe .

# Install
go install .

# Test
go test ./...
```

## Common Issues

### Issue: "acm is not recognized as an internal or external command"

**Cause:** ACM is not in your PATH.

**Solution:**

```powershell
# Option 1: Add Go bin to PATH permanently
[Environment]::SetEnvironmentVariable(
    "Path",
    [Environment]::GetEnvironmentVariable("Path", "User") + ";$env:USERPROFILE\go\bin",
    "User"
)
# Restart terminal

# Option 2: Use full path
& "$env:USERPROFILE\go\bin\acm.exe" list

# Option 3: Run from current directory (if built from source)
.\acm.exe list
```

### Issue: "Access Denied" or Permission Errors

**Cause:** Running without necessary permissions.

**When Administrator is needed:**
- Modifying PATH environment variables
- Installing to system-wide locations (C:\Program Files)

**When Administrator is NOT needed:**
- Normal ACM usage (list, clean)
- Installing to user directories (%USERPROFILE%)
- Building from source

**Solution:**
- For PATH modification: Run PowerShell as Administrator
- For regular use: Administrator rights are not required

### Issue: "No cache files found"

**Cause:** AI assistant hasn't been used yet, or ACM is looking in wrong locations.

**Solution:**

```powershell
# 1. Run diagnostics
acm doctor

# 2. Check with verbose mode
acm list -v

# 3. Verify AI assistant is installed
# For Claude Code: Check if C:\Users\YourName\.claude exists
Get-ChildItem -Path $env:USERPROFILE\.claude -Force

# 4. Use the AI assistant at least once to create cache files
```

### Issue: "GOPATH not set" or "Go bin not found"

**Cause:** GOPATH environment variable is not set.

**Solution:**

```powershell
# Check GOPATH
echo $env:GOPATH

# If empty, set it
[Environment]::SetEnvironmentVariable("GOPATH", "$env:USERPROFILE\go", "User")

# Restart terminal
```

By default, Go uses `%USERPROFILE%\go` if GOPATH is not set.

## Windows-Specific Features

### Cache Locations

ACM checks multiple Windows-specific locations:

1. **Unix-style** (backward compatible):
   - `%USERPROFILE%\.claude\cache`
   - `%USERPROFILE%\.cache\claude`

2. **Windows Standard**:
   - `%APPDATA%\claude\cache` (Roaming data)
   - `%LOCALAPPDATA%\Claude\Cache` (Local data)

### PowerShell Scripts

- `build.ps1` - Build automation
- `install.ps1` - Installation helper with PATH management

### Environment Variables

ACM recognizes these Windows environment variables:
- `USERPROFILE` - User home directory
- `APPDATA` - Roaming application data
- `LOCALAPPDATA` - Local application data
- `GOPATH` - Go workspace (optional)

### Path Handling

ACM automatically handles Windows path differences:
- Accepts both forward slashes (`/`) and backslashes (`\`)
- Uses correct path separator for the platform
- Handles drive letters (e.g., `C:\`)
- Case-insensitive file system on Windows

## Tips and Best Practices

### 1. Use PowerShell Scripts

The PowerShell scripts (`build.ps1`, `install.ps1`) are designed specifically for Windows and handle many edge cases automatically.

### 2. Run `acm doctor` After Installation

Always run `acm doctor` after installing or updating to verify everything is configured correctly.

### 3. Use Verbose Mode for Troubleshooting

When something doesn't work, use `-v` flag:

```powershell
acm list -v
acm clean -v
```

### 4. Keep Go Updated

ACM requires Go 1.24+. Check your version:

```powershell
go version
```

Update Go from [go.dev/dl](https://go.dev/dl/) if needed.

### 5. Use Windows Terminal (Optional)

[Windows Terminal](https://aka.ms/terminal) provides a better experience than Command Prompt:
- Better Unicode support
- Better colors
- Multiple tabs
- Better PowerShell integration

## Uninstalling

To completely remove ACM:

```powershell
# Remove binary
Remove-Item "$env:USERPROFILE\go\bin\acm.exe"

# Optionally remove Go bin from PATH
# (Manual: System Properties > Environment Variables)

# ACM doesn't create any config files or registry entries
```

## Getting Help

- **GitHub Issues**: https://github.com/buemura/acm/issues
- **Documentation**: https://github.com/buemura/acm
- **Health Check**: `acm doctor`
- **Verbose Mode**: `acm list -v` or `acm clean -v`
