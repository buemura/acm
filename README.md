# Agent Cache Manager (ACM)

A CLI tool to manage and clean up cache files from AI coding assistants (Claude Code, Codex). Identifies and safely removes cached session data, logs, and temporary files to free up storage space.

## Features

- List and delete cache files with filtering by type and age
- Multiple provider support (Claude, Codex) with extensible architecture
- Interactive menu and direct CLI command modes
- Human-readable file size and relative time formatting
- Confirmation prompt before deletion (with `--force` override)

## Installation

**Prerequisites:** Go 1.24+

```bash
# Clone and build
make build

# Or build directly
go build -o acm .
```

## Usage

### Interactive Mode

Run without arguments to launch the interactive menu:

```bash
./acm
```

### CLI Commands

**List cache files:**

```bash
acm list                          # List all Claude cache files (default provider)
acm list -p codex                 # List Codex cache files
acm list -t json                  # Filter by file extension
acm list -a 7d                    # Filter files older than 7 days
```

**Clean cache files:**

```bash
acm clean                         # Delete Claude cache files (with confirmation)
acm clean -p codex -a 30d         # Delete Codex files older than 30 days
acm clean -f                      # Skip confirmation prompt
```

### Flags

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--provider` | `-p` | Cache provider name (`claude`, `codex`) | `claude` |
| `--type` | `-t` | Filter by file extension (e.g., `json`, `log`, `jsonl`) | — |
| `--age` | `-a` | Filter by minimum file age (e.g., `20m`, `24h`, `7d`) | — |
| `--force` | `-f` | Skip deletion confirmation (clean only) | `false` |

## Managed Cache Paths

**Claude** (`~/.claude`):
- `~/.claude/cache/` and `~/.cache/claude/` — cache directories
- `~/.claude/projects/` — project session transcripts (`.jsonl`), session index, tool result caches
- `~/.claude/stats-cache.json` — usage statistics cache

**Codex** (`~/.codex`):
- `~/.codex/cache/` and `~/.cache/codex/` — cache directories
- `~/.codex/sessions/`, `session_index.jsonl`, `history.jsonl` — session data
- `~/.codex/models_cache.json` — models metadata cache
- `~/.codex/log/`, `logs_*.sqlite*` — log and database files

## Project Structure

```
├── cmd/           # CLI commands (root, list, clean)
├── provider/      # Provider registry and implementations (Claude, Codex)
├── scanner/       # File discovery, filtering, and deduplication
├── ui/            # Interactive menu and output formatting
├── main.go        # Entry point
├── Makefile       # Build automation
└── go.mod         # Go module definition
```

## Development

```bash
make build    # Build the binary
make run      # Build and run
make clean    # Remove built binary
go test ./... # Run tests
```
