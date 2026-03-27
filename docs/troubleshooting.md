# Troubleshooting Guide

Common issues and solutions for ACM on all platforms.

## Table of Contents

- [Quick Diagnostics](#quick-diagnostics)
- [Installation Issues](#installation-issues)
- [PATH Issues](#path-issues)
- [Provider Issues](#provider-issues)
- [Permission Issues](#permission-issues)
- [Platform-Specific Issues](#platform-specific-issues)

## Quick Diagnostics

### Step 1: Run Health Check

```bash
acm doctor
```

This command checks:
- Go installation
- ACM binary location
- PATH configuration
- Provider detection
- File permissions

### Step 2: Use Verbose Mode

```bash
acm list -v        # See which paths are checked
acm clean -v       # See what will be cleaned
```

Verbose mode shows:
- Which cache locations exist
- Which are not found
- File counts for each location
- Path check results

## Installation Issues

### "Command not found: acm" or "acm is not recognized"

**Symptoms:**
```bash
$ acm --version
bash: acm: command not found
```

**Cause:** ACM is not installed or not in PATH.

**Solutions:**

**Option 1: Install via go install**
```bash
go install github.com/buemura/acm@latest
```

**Option 2: Add Go bin to PATH**

Unix/macOS:
```bash
export PATH="$PATH:$HOME/go/bin"
# Add to ~/.bashrc or ~/.zshrc for persistence
```

Windows (PowerShell):
```powershell
[Environment]::SetEnvironmentVariable(
    "Path",
    $env:Path + ";$env:USERPROFILE\go\bin",
    "User"
)
# Restart terminal
```

**Option 3: Use full path**
```bash
# Unix/macOS
~/go/bin/acm --version

# Windows
%USERPROFILE%\go\bin\acm.exe --version
```

### "go: cannot find main module"

**Symptoms:**
```bash
$ go install .
go: go.mod file not found in current directory or any parent directory
```

**Cause:** Running `go install .` outside the ACM repository.

**Solution:**

If using `go install` to install from GitHub:
```bash
# Don't be in any directory
go install github.com/buemura/acm@latest
```

If building from source:
```bash
# Clone first
git clone https://github.com/buemura/acm.git
cd acm
go install .
```

### "package github.com/buemura/acm is not in GOROOT"

**Cause:** Old Go version or misconfigured GOPATH.

**Solution:**

1. Check Go version (need 1.24+):
```bash
go version
```

2. Update Go if needed: https://go.dev/dl/

3. Clear Go cache:
```bash
go clean -modcache
go install github.com/buemura/acm@latest
```

## PATH Issues

### ACM installed but not accessible

**Symptoms:**
- `acm doctor` shows "ACM installed but not in PATH"
- Need to use full path to run ACM

**Solution:**

**Unix/macOS:**

Add to `~/.bashrc`, `~/.zshrc`, or `~/.bash_profile`:
```bash
export PATH="$PATH:$HOME/go/bin"
```

Then reload:
```bash
source ~/.bashrc  # or ~/.zshrc
```

**Windows:**

PowerShell (Administrator):
```powershell
[Environment]::SetEnvironmentVariable(
    "Path",
    [Environment]::GetEnvironmentVariable("Path", "User") + ";$env:USERPROFILE\go\bin",
    "User"
)
```

Restart terminal after adding to PATH.

### Verify PATH configuration

```bash
# Unix/macOS
echo $PATH | tr ':' '\n' | grep go

# Windows (PowerShell)
$env:Path -split ';' | Select-String "go"

# Check if acm is found
which acm          # Unix/macOS
Get-Command acm    # Windows PowerShell
```

## Provider Issues

### "No claude cache found. Is claude installed?"

**Cause:** Claude Code hasn't been used yet, or cache is in unexpected location.

**Solutions:**

**Step 1: Verify Claude Code is installed**

Check if cache directory exists:

Unix/macOS:
```bash
ls -la ~/.claude
```

Windows (PowerShell):
```powershell
Get-ChildItem -Path $env:USERPROFILE\.claude -Force
```

**Step 2: Use Claude Code at least once**

Cache files are only created after using Claude Code. Run a simple command or open a project with Claude Code.

**Step 3: Check with verbose mode**

```bash
acm list -v -p claude
```

This shows all locations ACM checks.

**Step 4: Check alternative locations**

ACM checks these locations (in order):

Unix/macOS:
- `~/.claude/cache`
- `~/.cache/claude`
- `~/.claude/projects`

Windows:
- `%USERPROFILE%\.claude\cache`
- `%USERPROFILE%\.cache\claude`
- `%APPDATA%\claude\cache`
- `%LOCALAPPDATA%\Claude\Cache`
- `%USERPROFILE%\.claude\projects`

### "Provider missing expected cache files"

**Symptoms:**
```
Note: claude provider is missing expected stats cache in standard locations.
```

**Cause:** Some cache files haven't been created yet. This is normal if you've only used Claude Code briefly.

**Solution:**

This is usually just a warning, not an error. ACM will still work with the cache files that exist.

To see what's found:
```bash
acm list -v
```

### Wrong provider selected

**Symptoms:** No files found, but you know cache exists.

**Solution:**

Specify the provider explicitly:

```bash
acm list -p claude      # For Claude Code
acm list -p codex       # For Codex

acm provider            # List all available providers
```

## Permission Issues

### "Permission denied" when listing files

**Cause:** Insufficient permissions to read cache directories.

**Solutions:**

**On Unix/macOS:**

Check permissions:
```bash
ls -la ~/.claude
```

Fix if needed:
```bash
chmod 755 ~/.claude
chmod 644 ~/.claude/cache/*
```

**On Windows:**

Usually not an issue since you own the files in your user directory. If you see permission errors:

1. Run as your regular user (not Administrator)
2. Check if antivirus is blocking access
3. Verify file ownership:
```powershell
Get-Acl $env:USERPROFILE\.claude
```

### "Permission denied" when deleting files

**Cause:** Files are in use or read-only.

**Solutions:**

1. Close Claude Code/Codex before cleaning
2. Check if files are read-only:

Unix/macOS:
```bash
ls -la ~/.claude/cache
```

Windows:
```powershell
Get-ChildItem -Path $env:USERPROFILE\.claude\cache | Select-Object Name, IsReadOnly
```

3. Use force flag (use with caution):
```bash
acm clean -f
```

## Platform-Specific Issues

### macOS: "acm cannot be opened because the developer cannot be verified"

**Cause:** macOS Gatekeeper blocking unsigned binary.

**Solution:**

If you built from source or downloaded from releases:
```bash
xattr -d com.apple.quarantine ~/go/bin/acm
```

Or: System Preferences > Security & Privacy > Allow

**Better solution:** Install via `go install`:
```bash
go install github.com/buemura/acm@latest
```

### Windows: "Windows protected your PC" SmartScreen warning

**Cause:** Binary is not signed with a certificate.

**Solution:**

Click "More info" > "Run anyway"

Or install via `go install`:
```powershell
go install github.com/buemura/acm@latest
```

### Linux: Snap/Flatpak home directory isolation

**Cause:** Snap/Flatpak applications may use different home directories.

**Solution:**

Check if Claude Code is installed via Snap:
```bash
snap list | grep claude
```

If yes, cache might be in:
```bash
~/snap/claude-code/current/.claude
```

You may need to add custom provider paths (feature request - not yet implemented).

### WSL (Windows Subsystem for Linux)

**Cache location in WSL:**

If you run Claude Code on Windows but ACM in WSL, specify Windows paths:

```bash
# Access Windows home from WSL
cd /mnt/c/Users/YourName/.claude
```

Or run ACM from Windows PowerShell instead of WSL.

## Performance Issues

### "acm list" is slow

**Cause:** Large number of cache files or slow filesystem.

**Solutions:**

1. Filter by file type:
```bash
acm list -t jsonl      # Only .jsonl files
```

2. Filter by age:
```bash
acm list -a 7d         # Only files older than 7 days
```

3. Clean old files:
```bash
acm clean -a 30d       # Delete files older than 30 days
```

### "acm clean" is slow

**Cause:** Many files to delete.

**Solution:**

This is expected. Deleting thousands of files takes time. Progress is shown during deletion.

## Debug Mode

### Get more information

1. Use verbose flag:
```bash
acm list -v
acm clean -v
```

2. Run diagnostics:
```bash
acm doctor
```

3. Check Go environment:
```bash
go env
```

### Report a bug

When reporting issues, include:

1. Output of `acm doctor`
2. Output of `acm list -v`
3. Your OS and version
4. Go version (`go version`)
5. How you installed ACM

## Getting Help

- **Run diagnostics**: `acm doctor`
- **Use verbose mode**: `acm list -v`
- **Check documentation**: https://github.com/buemura/acm
- **Report issues**: https://github.com/buemura/acm/issues
- **Windows-specific help**: See [docs/windows.md](windows.md)
