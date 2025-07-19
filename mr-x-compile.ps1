# This script compiles and installs the mr-x CLI for Windows.
# It makes the 'mr-x' command available in both PowerShell and CMD.
# To run: Right-click this file and select "Run with PowerShell".

# --- FIX: Change the current location to the script's directory ---
# This ensures that 'go build' runs in the correct project folder.
Set-Location -Path $PSScriptRoot

Write-Host "▶️ Compiling Mr. Albert X for Windows..." -ForegroundColor Yellow

# Set environment variables for cross-compilation to Windows
$env:GOOS = "windows"
$env:GOARCH = "amd64"

# Compile the Go program. The output will be 'mr-x.exe'.
go build -o mr-x.exe .

if (-not $?) {
    Write-Host "❌ Compilation failed. Please ensure Go is installed correctly." -ForegroundColor Red
    pause
    exit 1
}

Write-Host "✅ Compilation successful." -ForegroundColor Green
Write-Host "▶️ Installing mr-x command..." -ForegroundColor Yellow

# Define a standard installation directory in the user's profile.
$installPath = "$env:USERPROFILE\.local\bin"

# Create the directory if it doesn't exist.
if (-not (Test-Path $installPath)) {
    New-Item -ItemType Directory -Force -Path $installPath | Out-Null
    Write-Host "   -> Created installation directory at $installPath"
}

# Move the compiled executable to the installation directory.
Move-Item -Path ".\mr-x.exe" -Destination $installPath -Force

# Check if the installation directory is already in the user's PATH.
$currentUserPath = [System.Environment]::GetEnvironmentVariable("Path", "User")

if (-not ($currentUserPath -like "*$installPath*")) {
    Write-Host "   -> Adding installation directory to your PATH..."
    $newPath = "$currentUserPath;$installPath"
    [System.Environment]::SetEnvironmentVariable("Path", $newPath, "User")
    Write-Host "   -> PATH updated successfully."
} else {
    Write-Host "   -> Installation directory is already in your PATH."
}

# Clean up the environment variables used for compilation
Remove-Item Env:\GOOS
Remove-Item Env:\GOARCH

Write-Host ""
Write-Host "🚀 Mr. Albert X installed successfully!" -ForegroundColor Green
Write-Host "You can now run 'mr-x' from any new PowerShell or CMD window."
Write-Host "IMPORTANT: Please close and reopen your terminal for the changes to take effect."
pause
