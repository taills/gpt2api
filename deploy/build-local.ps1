# gpt2api Windows 预构建脚本 — 生成单文件可执行 (Linux/amd64)
#
# 用法:
#   powershell -NoProfile -File deploy/build-local.ps1            # 常规构建
#   powershell -NoProfile -File deploy/build-local.ps1 -NoWeb    # 跳过前端构建

param(
    [switch]$NoWeb
)

$ErrorActionPreference = 'Stop'
if ($PSVersionTable.PSVersion.Major -ge 7) {
    $PSNativeCommandUseErrorActionPreference = $false
}

$root     = Resolve-Path "$PSScriptRoot/.."
$embedDir = Join-Path $root "internal/server/web"
Set-Location $root

Write-Host "[build-local] repo = $root"

# ---- step1: 前端 ----
if (-not $NoWeb) {
    Write-Host "[build-local] step1 = npm run build (web)"
    Push-Location (Join-Path $root "web")
    try {
        if (-not (Test-Path node_modules)) {
            npm install --no-audit --no-fund --loglevel=error
            if ($LASTEXITCODE -ne 0) { throw "npm install failed" }
        }
        npm run build
        if ($LASTEXITCODE -ne 0) { throw "npm run build failed" }
    } finally {
        Pop-Location
    }

    Write-Host "[build-local] step1 = copy dist → internal/server/web/"
    if (Test-Path $embedDir) { Remove-Item -Recurse -Force $embedDir }
    New-Item -ItemType Directory -Force $embedDir | Out-Null
    Copy-Item (Join-Path $root "web/dist/*") $embedDir -Recurse
} else {
    Write-Host "[build-local] step1 = skipped (-NoWeb)"
}

# ---- step2: 后端 ----
Write-Host "[build-local] step2 = cross-build gpt2api (linux/amd64)"
$env:GOOS      = "linux"
$env:GOARCH    = "amd64"
$env:CGO_ENABLED = "0"
New-Item -ItemType Directory -Force deploy/bin | Out-Null
go build -ldflags "-s -w" -o deploy/bin/gpt2api ./cmd/server
if ($LASTEXITCODE -ne 0) { throw "go build failed" }

Write-Host "[build-local] done."
Get-Item deploy/bin/gpt2api | Format-Table Name, Length -AutoSize
