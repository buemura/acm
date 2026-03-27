# Install script for ACM (Agent Cache Manager)
# Installs the binary to Go bin directory and helps with PATH configuration

param(
    [switch]$Force,
    [switch]$CheckOnly,
    [switch]$Help
)

$ErrorActionPreference = "Stop"

$BinaryName = "acm.exe"
$SourceBinary = ".\$BinaryName"

function Show-Help {
    Write-Host @"

ACM Install Script
==================

Usage: .\install.ps1 [options]

Options:
  -Force       Overwrite existing binary without prompting
  -CheckOnly   Check installation status without installing
  -Help        Show this help message

Examples:
  .\install.ps1            # Install ACM to Go bin directory
  .\install.ps1 -Force     # Force reinstall
  .\install.ps1 -CheckOnly # Check if ACM is installed

"@ -ForegroundColor White
}

function Get-GoBinPath {
    $goPath = $env:GOPATH
    if (-not $goPath) {
        $goPath = "$env:USERPROFILE\go"
    }
    return "$goPath\bin"
}

function Test-InPath {
    param([string]$Directory)

    $paths = $env:Path -split ';'
    foreach ($path in $paths) {
        if ($path.TrimEnd('\') -eq $Directory.TrimEnd('\')) {
            return $true
        }
    }
    return $false
}

function Show-PathInstructions {
    param([string]$Directory)

    Write-Host "`n[WARNING] $Directory is not in your PATH" -ForegroundColor Yellow
    Write-Host "`nTo add it permanently, run this command in an Administrator PowerShell:" -ForegroundColor Cyan
    Write-Host @"

[Environment]::SetEnvironmentVariable(
    "Path",
    [Environment]::GetEnvironmentVariable("Path", "User") + ";$Directory",
    "User"
)

"@ -ForegroundColor White
    Write-Host "After adding to PATH, restart your terminal and run: acm --version" -ForegroundColor Gray
}

function Check-Installation {
    $goBinPath = Get-GoBinPath
    $targetPath = "$goBinPath\$BinaryName"

    Write-Host "`n=== ACM Installation Status ===" -ForegroundColor Cyan

    # Check if binary exists in Go bin
    if (Test-Path $targetPath) {
        Write-Host "[OK] ACM is installed at: $targetPath" -ForegroundColor Green

        # Check version
        try {
            $version = & $targetPath --version 2>&1 | Select-Object -First 1
            Write-Host "  Version: $version" -ForegroundColor Gray
        } catch {
            Write-Host "  (Unable to determine version)" -ForegroundColor Gray
        }
    } else {
        Write-Host "[NOT FOUND] ACM is not installed" -ForegroundColor Red
        Write-Host "  Expected location: $targetPath" -ForegroundColor Gray
    }

    # Check if Go bin is in PATH
    if (Test-InPath -Directory $goBinPath) {
        Write-Host "[OK] Go bin directory is in PATH" -ForegroundColor Green
    } else {
        Show-PathInstructions -Directory $goBinPath
    }

    # Check if acm command is available
    Write-Host "`nTesting command availability..." -ForegroundColor Gray
    try {
        $null = Get-Command acm -ErrorAction Stop
        Write-Host "[OK] 'acm' command is available in current session" -ForegroundColor Green
    } catch {
        Write-Host "[NOT AVAILABLE] 'acm' command is not available" -ForegroundColor Red
        Write-Host "  You may need to restart your terminal or update PATH" -ForegroundColor Gray
    }

    Write-Host ""
}

function Install-Binary {
    # Check if source binary exists
    if (-not (Test-Path $SourceBinary)) {
        Write-Host "[ERROR] Binary not found: $SourceBinary" -ForegroundColor Red
        Write-Host "  Build it first by running: .\build.ps1" -ForegroundColor Yellow
        exit 1
    }

    # Get Go bin path
    $goBinPath = Get-GoBinPath

    # Create Go bin directory if it doesn't exist
    if (-not (Test-Path $goBinPath)) {
        Write-Host "Creating Go bin directory: $goBinPath" -ForegroundColor Cyan
        New-Item -ItemType Directory -Path $goBinPath -Force | Out-Null
    }

    $targetPath = "$goBinPath\$BinaryName"

    # Check if binary already exists
    if ((Test-Path $targetPath) -and -not $Force) {
        $response = Read-Host "$BinaryName already exists. Overwrite? (y/N)"
        if ($response -ne 'y' -and $response -ne 'Y') {
            Write-Host "Installation cancelled." -ForegroundColor Yellow
            exit 0
        }
    }

    # Copy binary
    Write-Host "`nInstalling ACM..." -ForegroundColor Cyan
    Copy-Item -Path $SourceBinary -Destination $targetPath -Force

    Write-Host "[OK] Installed to: $targetPath" -ForegroundColor Green

    # Get file size
    $size = (Get-Item $targetPath).Length
    $sizeMB = [math]::Round($size / 1MB, 2)
    Write-Host "  Size: $sizeMB MB" -ForegroundColor Gray

    # Check PATH
    if (-not (Test-InPath -Directory $goBinPath)) {
        Show-PathInstructions -Directory $goBinPath
    } else {
        Write-Host "`n[OK] Installation complete!" -ForegroundColor Green
        Write-Host "  Run 'acm --version' to verify" -ForegroundColor Gray
    }

    Write-Host ""
}

# Main execution
if ($Help) {
    Show-Help
}
elseif ($CheckOnly) {
    Check-Installation
}
else {
    Install-Binary
    Write-Host "`nVerifying installation..." -ForegroundColor Cyan
    Check-Installation
}
