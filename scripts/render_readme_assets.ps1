$ErrorActionPreference = "Stop"

$root = Resolve-Path (Join-Path $PSScriptRoot "..")
$screenshots = Join-Path $root "screenshots"
$port = 5514
$process = $null
$stdout = Join-Path $env:TEMP ("gov-comment-ingestor-" + [guid]::NewGuid().ToString() + "-stdout.log")
$stderr = Join-Path $env:TEMP ("gov-comment-ingestor-" + [guid]::NewGuid().ToString() + "-stderr.log")
$edgeCandidates = @(
    "C:\Program Files (x86)\Microsoft\Edge\Application\msedge.exe",
    "C:\Program Files\Microsoft\Edge\Application\msedge.exe"
)

New-Item -ItemType Directory -Force -Path $screenshots | Out-Null

function Get-EdgePath {
    foreach ($candidate in $edgeCandidates) {
        if (Test-Path $candidate) {
            return $candidate
        }
    }

    throw "Microsoft Edge was not found."
}

function Wait-ForUrl {
    param([string]$Url)
    for ($i = 0; $i -lt 40; $i++) {
        try {
            Invoke-WebRequest -Uri $Url -UseBasicParsing | Out-Null
            return
        } catch {
            Start-Sleep -Milliseconds 750
        }
    }

    throw "Timed out waiting for $Url"
}

try {
    $process = Start-Process -FilePath "go.exe" `
        -ArgumentList "run", "./cmd/server" `
        -WorkingDirectory $root `
        -RedirectStandardOutput $stdout `
        -RedirectStandardError $stderr `
        -WindowStyle Hidden `
        -PassThru

    Wait-ForUrl "http://127.0.0.1:$port/"

    $edge = Get-EdgePath
    $targets = @(
        @{ Url = "http://127.0.0.1:$port/"; File = "01-overview-proof.png"; Size = "1600,1450" },
        @{ Url = "http://127.0.0.1:$port/ingest-lane"; File = "02-ingest-lane-proof.png"; Size = "1600,1360" },
        @{ Url = "http://127.0.0.1:$port/source-adapters"; File = "03-source-adapters-proof.png"; Size = "1600,1320" },
        @{ Url = "http://127.0.0.1:$port/verification"; File = "04-verification-proof.png"; Size = "1600,1180" }
    )

    foreach ($target in $targets) {
        & $edge `
            --headless `
            --disable-gpu `
            --hide-scrollbars `
            "--window-size=$($target.Size)" `
            "--screenshot=$(Join-Path $screenshots $target.File)" `
            $target.Url | Out-Null
    }
} finally {
    if ($process -and -not $process.HasExited) {
        Stop-Process -Id $process.Id -Force
        $process.WaitForExit()
    }

    if (Test-Path $stdout) {
        try {
            Remove-Item $stdout -Force
        } catch {
            Start-Sleep -Milliseconds 250
            Remove-Item $stdout -Force -ErrorAction SilentlyContinue
        }
    }

    if (Test-Path $stderr) {
        try {
            Remove-Item $stderr -Force
        } catch {
            Start-Sleep -Milliseconds 250
            Remove-Item $stderr -Force -ErrorAction SilentlyContinue
        }
    }
}
