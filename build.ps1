# Build script for ACM (Agent Cache Manager)
# PowerShell alternative to Makefile for Windows users

param(
    [switch]$Release,
    [switch]$Clean,
    [switch]$Run,
    [switch]$Test,
    [switch]$All
)

$ErrorActionPreference = "Stop"

$BinaryName = "acm.exe"
$BuildFlags = @()

if ($Release) {
    $BuildFlags += "-ldflags", "-s -w"
    Write-Host "Building release version..." -ForegroundColor Green
} else {
    Write-Host "Building debug version..." -ForegroundColor Green
}

function Build-Binary {
    Write-Host "`nBuilding $BinaryName..." -ForegroundColor Cyan

    if ($Release) {
        go build -ldflags "-s -w" -o $BinaryName .
    } else {
        go build -o $BinaryName .
    }

    if ($LASTEXITCODE -eq 0) {
        Write-Host "[OK] Build successful: $BinaryName" -ForegroundColor Green

        # Show binary size
        $size = (Get-Item $BinaryName).Length
        $sizeKB = [math]::Round($size / 1KB, 2)
        $sizeMB = [math]::Round($size / 1MB, 2)
        Write-Host "  Size: $sizeMB MB" -ForegroundColor Gray
    } else {
        Write-Host "[ERROR] Build failed" -ForegroundColor Red
        exit 1
    }
}

function Clean-Build {
    Write-Host "`nCleaning build artifacts..." -ForegroundColor Cyan

    if (Test-Path $BinaryName) {
        Remove-Item $BinaryName -Force
        Write-Host "[OK] Removed $BinaryName" -ForegroundColor Green
    } else {
        Write-Host "  Nothing to clean" -ForegroundColor Gray
    }
}

function Run-Binary {
    if (-not (Test-Path $BinaryName)) {
        Write-Host "`nBinary not found, building first..." -ForegroundColor Yellow
        Build-Binary
    }

    Write-Host "`nRunning $BinaryName..." -ForegroundColor Cyan
    & ".\$BinaryName" $args
}

function Run-Tests {
    Write-Host "`nRunning tests..." -ForegroundColor Cyan
    go test ./... -v

    if ($LASTEXITCODE -eq 0) {
        Write-Host "[OK] All tests passed" -ForegroundColor Green
    } else {
        Write-Host "[ERROR] Tests failed" -ForegroundColor Red
        exit 1
    }
}

function Show-Help {
    Write-Host @"

ACM Build Script
================

Usage: .\build.ps1 [options]

Options:
  -Release     Build optimized release version (smaller binary)
  -Clean       Remove built binary
  -Run         Build and run the binary
  -Test        Run all tests
  -All         Clean, test, and build release version
  -Help        Show this help message

Examples:
  .\build.ps1              # Build debug version
  .\build.ps1 -Release     # Build release version
  .\build.ps1 -Clean       # Clean build artifacts
  .\build.ps1 -Test        # Run tests
  .\build.ps1 -All         # Full build pipeline

"@ -ForegroundColor White
}

# Main execution
if ($All) {
    Clean-Build
    Run-Tests
    $Release = $true
    Build-Binary
    Write-Host "`n[OK] All tasks completed successfully!" -ForegroundColor Green
}
elseif ($Clean) {
    Clean-Build
}
elseif ($Test) {
    Run-Tests
}
elseif ($Run) {
    Run-Binary
}
elseif ($args.Count -gt 0 -and ($args[0] -eq "-h" -or $args[0] -eq "--help")) {
    Show-Help
}
else {
    Build-Binary
}
