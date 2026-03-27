# Agent Cache Manager (ACM)

A CLI tool to manage and clean up cache files from AI coding assistants (Claude Code, Codex). Identifies and safely removes cached session data, logs, and temporary files to free up storage space.

## Features

- List and delete cache files with filtering by type and age
- Multiple provider support (Claude, Codex) with extensible architecture
- Human-readable file size and relative time formatting
- Confirmation prompt before deletion (with `--force` override)

## Installation

### Via `go install`

**Prerequisites:** Go 1.24+

```bash
go install github.com/buemura/acm@latest
```

#### Windows Setup

After installation on Windows, you may need to add Go's bin directory to your PATH:

**PowerShell (Run as Administrator):**
```powershell
# Check if Go bin is in PATH
$env:Path -split ';' | Select-String "go\\bin"

# If not found, add it permanently
[Environment]::SetEnvironmentVariable(
    "Path",
    [Environment]::GetEnvironmentVariable("Path", "User") + ";$env:USERPROFILE\go\bin",
    "User"
)
```

**Verify installation:**
```powershell
# Restart your terminal, then run:
acm --version
acm provider
```

**Common Windows installation locations:**
- Binary: `%USERPROFILE%\go\bin\acm.exe`
- Go path: `C:\Users\<username>\go\bin\acm.exe`

### From source

**Unix/macOS:**
```bash
git clone https://github.com/buemura/acm.git
cd acm
make build
```

**Windows (PowerShell):**
```powershell
git clone https://github.com/buemura/acm.git
cd acm
go build -o acm.exe .
# Optional: Move to Go bin directory
Move-Item acm.exe "$env:USERPROFILE\go\bin\acm.exe"
```

## Usage

Run without arguments to see available commands and flags:

```bash
acm
```

### Commands

**List cache files:**

```bash
acm list                          # List all Claude cache files (default provider)
acm list -p codex                 # List Codex cache files
acm list -t json                  # Filter by file extension
acm list -a 7d                    # Filter files older than 7 days
acm list -v                       # Show verbose output (checked paths, what was found)
```

**Clean cache files:**

```bash
acm clean                         # Delete Claude cache files (with confirmation)
acm clean -p codex -a 30d         # Delete Codex files older than 30 days
acm clean -f                      # Skip confirmation prompt
acm clean -v                      # Show verbose output before cleaning
```

**Show supported providers:**

```bash
acm provider                      # List all available providers
```

### Flags

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--provider` | `-p` | Cache provider name (`claude`, `codex`) | `claude` |
| `--type` | `-t` | Filter by file extension (e.g., `json`, `log`, `jsonl`) | — |
| `--age` | `-a` | Filter by minimum file age (e.g., `20m`, `24h`, `7d`) | — |
| `--verbose` | `-v` | Show detailed path information and what was found | `false` |
| `--force` | `-f` | Skip deletion confirmation (clean only) | `false` |

## Managed Cache Paths

### Claude

**Unix/macOS:**
- `~/.claude/cache/` and `~/.cache/claude/` — cache directories
- `~/.claude/projects/` — project session transcripts (`.jsonl`), session index, tool result caches
- `~/.claude/stats-cache.json` — usage statistics cache

**Windows:**
- `%USERPROFILE%\.claude\cache\` — main cache directory (e.g., `C:\Users\username\.claude\cache\`)
- `%USERPROFILE%\.cache\claude\` — alternate cache location
- `%USERPROFILE%\.claude\projects\` — project session transcripts
- `%USERPROFILE%\.claude\stats-cache.json` — usage statistics cache

### Codex

**Unix/macOS:**
- `~/.codex/cache/` and `~/.cache/codex/` — cache directories
- `~/.codex/sessions/`, `session_index.jsonl`, `history.jsonl` — session data
- `~/.codex/models_cache.json` — models metadata cache
- `~/.codex/log/`, `logs_*.sqlite*` — log and database files

**Windows:**
- `%USERPROFILE%\.codex\cache\` — main cache directory
- `%USERPROFILE%\.cache\codex\` — alternate cache location
- `%USERPROFILE%\.codex\sessions\` — session data
- `%USERPROFILE%\.codex\models_cache.json` — models metadata cache
- `%USERPROFILE%\.codex\log\` — log files

## Project Structure

```
├── main.go                      # Application entry point
├── cmd/
│   └── cli/                     # CLI commands (root, list, clean, provider)
├── internal/
│   ├── provider/                # Provider registry and implementations (Claude, Codex)
│   ├── scanner/                 # File discovery, filtering, and deduplication
│   └── ui/                      # Interactive menu and output formatting
├── Makefile                     # Build automation
└── go.mod                       # Go module definition
```

## Development

**Unix/macOS:**
```bash
make build    # Build the binary
make run      # Build and run
make clean    # Remove built binary
go test ./... # Run tests
```

**Windows (PowerShell scripts):**
```powershell
.\build.ps1              # Build debug version
.\build.ps1 -Release     # Build optimized release version
.\build.ps1 -Clean       # Clean build artifacts
.\build.ps1 -Test        # Run all tests
.\build.ps1 -All         # Clean, test, and build release

.\install.ps1            # Install to Go bin directory
.\install.ps1 -CheckOnly # Check installation status
.\install.ps1 -Force     # Force reinstall
```

**Windows (direct Go commands):**
```powershell
go build -o acm.exe .        # Build the binary
go run .                     # Build and run
Remove-Item acm.exe          # Remove built binary
go test ./...                # Run tests
```

## Troubleshooting

### Windows: "acm is not recognized as an internal or external command"

This means the Go bin directory is not in your PATH. Solutions:

1. **Add to PATH permanently** (recommended):
   ```powershell
   [Environment]::SetEnvironmentVariable(
       "Path",
       [Environment]::GetEnvironmentVariable("Path", "User") + ";$env:USERPROFILE\go\bin",
       "User"
   )
   ```
   Then restart your terminal.

2. **Run with full path**:
   ```powershell
   & "$env:USERPROFILE\go\bin\acm.exe" list
   ```

3. **Use from current directory** (if built from source):
   ```powershell
   .\acm.exe list
   ```

### Windows: "Access denied" or permission errors

Run PowerShell as Administrator when:
- Modifying PATH environment variables
- Installing to system-wide locations

For regular usage, Administrator rights are **not** required.

### No cache files found

If `acm list` returns no results:

1. **Use verbose mode to see what paths are being checked**:
   ```bash
   acm list -v
   ```
   This shows which cache locations exist and which ones were not found.

2. **Verify cache locations exist**:
   ```bash
   # Unix/macOS
   ls -la ~/.claude
   ls -la ~/.codex

   # Windows (PowerShell)
   Get-ChildItem -Path $env:USERPROFILE\.claude -Force
   Get-ChildItem -Path $env:USERPROFILE\.codex -Force
   ```

3. **Check if you've used the AI assistant** - cache files are only created after using Claude Code or Codex

4. **Try a different provider**:
   ```bash
   acm list -p codex
   acm list -p claude
   ```

### Cross-platform path issues

ACM automatically handles path differences between platforms:
- Unix/macOS: `/home/user/.claude`
- Windows: `C:\Users\user\.claude`

Both forward slashes (`/`) and backslashes (`\`) work on Windows.
