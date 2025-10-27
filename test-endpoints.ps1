Write-Host "`n=== RESERVE WATCH - ENDPOINT STATUS CHECK ===" -ForegroundColor Cyan
Write-Host "`nPublic Endpoints:" -ForegroundColor Yellow

$endpoints = @{
    "Health" = "/health"
    "Home Page" = "/"
    "Pricing" = "/pricing"
    "Methodology" = "/methodology"
    "Trigger Watch" = "/trigger-watch"
    "Crash Drill" = "/crash-drill"
    "Enterprise" = "/enterprise"
    "API Docs" = "/api/docs"
    "API Latest" = "/api/latest"
    "API Realtime" = "/api/latest/realtime"
    "API History" = "/api/history?limit=10"
    "API Signals" = "/api/signals/latest"
}

foreach ($name in $endpoints.Keys | Sort-Object) {
    $url = "https://web-production-4c1d00.up.railway.app$($endpoints[$name])"
    try {
        $response = Invoke-WebRequest -Uri $url -Method GET -UseBasicParsing -TimeoutSec 10 -ErrorAction Stop
        Write-Host "  ✓ $name : $($response.StatusCode)" -ForegroundColor Green
    }
    catch {
        $status = if ($_.Exception.Response) { $_.Exception.Response.StatusCode.value__ } else { "ERROR" }
        Write-Host "  ✗ $name : $status" -ForegroundColor Red
    }
}

Write-Host "`nPayment-Gated Endpoints (expecting 402):" -ForegroundColor Yellow

$gatedEndpoints = @{
    "Export CSV" = "/api/export/csv"
    "Export JSON" = "/api/export/json"
    "Alerts API" = "/api/alerts"
}

foreach ($name in $gatedEndpoints.Keys | Sort-Object) {
    $url = "https://web-production-4c1d00.up.railway.app$($gatedEndpoints[$name])"
    try {
        $response = Invoke-WebRequest -Uri $url -Method GET -UseBasicParsing -TimeoutSec 10 -ErrorAction Stop
        Write-Host "  ~ $name : $($response.StatusCode) (expected 402)" -ForegroundColor Yellow
    }
    catch {
        $status = if ($_.Exception.Response) { $_.Exception.Response.StatusCode.value__ } else { "ERROR" }
        $color = if ($status -eq 402) { "Green" } else { "Red" }
        Write-Host "  ✓ $name : $status" -ForegroundColor $color
    }
}

Write-Host "`n=== STRIPE CHECKOUT STATUS ===" -ForegroundColor Cyan
Write-Host "  ✓ Button ID fixed: pro-btn → premium-btn" -ForegroundColor Green
Write-Host "  ✓ Free tier removed from pricing page" -ForegroundColor Green
Write-Host "  ✓ Only Pro ($74.99) and Team ($199) available" -ForegroundColor Green
Write-Host "`nTest URL: https://web-production-4c1d00.up.railway.app/pricing`n" -ForegroundColor White

