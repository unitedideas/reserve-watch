# Quick Local Test Script for Reserve Watch

Write-Host "`n🚀 Reserve Watch - Local Test Suite`n" -ForegroundColor Cyan

# Step 1: Build
Write-Host "Step 1: Building..." -ForegroundColor Yellow
go build -o reserve-watch.exe cmd/runner/main.go

if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ Build failed!" -ForegroundColor Red
    exit 1
}
Write-Host "✅ Build successful!`n" -ForegroundColor Green

# Step 2: Run tests
Write-Host "Step 2: Running tests..." -ForegroundColor Yellow
go test ./...

if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ Tests failed!" -ForegroundColor Red
    exit 1
}
Write-Host "✅ Tests passed!`n" -ForegroundColor Green

# Step 3: Check code
Write-Host "Step 3: Vetting code..." -ForegroundColor Yellow
go vet ./...

if ($LASTEXITCODE -ne 0) {
    Write-Host "⚠️ Code issues found!" -ForegroundColor Yellow
}
Write-Host "✅ Code check complete!`n" -ForegroundColor Green

# Step 4: Start server
Write-Host "Step 4: Starting server..." -ForegroundColor Cyan
Write-Host "📍 Homepage:      http://localhost:8080/" -ForegroundColor White
Write-Host "📍 Health:        http://localhost:8080/health" -ForegroundColor White
Write-Host "📍 API Latest:    http://localhost:8080/api/latest" -ForegroundColor White
Write-Host "📍 API Realtime:  http://localhost:8080/api/latest/realtime" -ForegroundColor White
Write-Host "`nPress Ctrl+C to stop`n" -ForegroundColor Yellow

$env:FRED_API_KEY = "b7cb42380ac4ab4708ff13b305755de5"
$env:PORT = "8080"
$env:DB_DSN = "file:./data/reserve_watch.db?_fk=1"
$env:DRY_RUN = "true"
$env:LOG_LEVEL = "info"

.\reserve-watch.exe

