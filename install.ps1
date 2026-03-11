Param(
    [string]$RepoOwner = "thejamesnick",
    [string]$RepoName = "keysync",
    [string]$BinaryName = "keysync.exe"
)

Write-Host "🔐 Installing KeySync for Windows..." -ForegroundColor Cyan

# 1. Detect Arch
$arch = $env:PROCESSOR_ARCHITECTURE.ToLower()
switch ($arch) {
    "amd64" { $target = "windows-amd64" }
    "arm64" { $target = "windows-arm64" }
    default {
        Write-Error "❌ Unsupported architecture: $arch"
        exit 1
    }
}

Write-Host "  🔍 Detected: windows/$arch (Target: $target)"

# 2. Find Latest Release Asset URL from GitHub
$releasesUrl = "https://api.github.com/repos/$RepoOwner/$RepoName/releases/latest"

try {
    $response = Invoke-RestMethod -Uri $releasesUrl -Headers @{ "User-Agent" = "keysync-installer" }
} catch {
    Write-Error "❌ Error: Failed to query GitHub releases API."
    exit 1
}

if (-not $response.assets) {
    Write-Error "❌ Error: No assets found on latest release. Check https://github.com/$RepoOwner/$RepoName/releases"
    exit 1
}

$asset = $response.assets | Where-Object { $_.browser_download_url -like "*$target*" } | Select-Object -First 1

if (-not $asset) {
    Write-Error "❌ Error: Could not find release asset for $target."
    Write-Host "   Ensure a Windows release exists at https://github.com/$RepoOwner/$RepoName/releases"
    exit 1
}

$assetUrl = $asset.browser_download_url
Write-Host "  📥 Downloading from GitHub..."

$tempPath = Join-Path $env:TEMP $BinaryName

try {
    Invoke-WebRequest -Uri $assetUrl -OutFile $tempPath -UseBasicParsing
} catch {
    Write-Error "❌ Error: Failed to download asset."
    exit 1
}

# 3. Install to a user-local bin directory (no admin needed)
$installDir = Join-Path $env:USERPROFILE "bin"
if (-not (Test-Path $installDir)) {
    New-Item -ItemType Directory -Path $installDir | Out-Null
}

$destPath = Join-Path $installDir $BinaryName
Move-Item -Force -Path $tempPath -Destination $destPath

Write-Host "  ✅ Installed to $destPath" -ForegroundColor Green

# 4. PATH hint
if ($env:PATH -notlike "*$installDir*") {
    Write-Host ""
    Write-Host "ℹ️ Add this directory to your PATH to run 'keysync' from anywhere:" -ForegroundColor Yellow
    Write-Host "   $installDir"
    Write-Host ""
    Write-Host "Windows 10/11:" -ForegroundColor Yellow
    Write-Host "  1. Press Win + R, type 'systempropertiesadvanced' and press Enter."
    Write-Host "  2. Click 'Environment Variables...'."
    Write-Host "  3. Under 'User variables', select 'Path' → Edit → New → add:" 
    Write-Host "     $installDir"
    Write-Host "  4. Open a new PowerShell window and run: keysync --help"
}

Write-Host ""
Write-Host "  ✅ Done! Open a new PowerShell window and run:" -ForegroundColor Green
Write-Host "     keysync --help"

