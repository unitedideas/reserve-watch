# 🚀 Local Development Setup

## Quick Start (2 minutes)

### 1. Create `.env` file:
```bash
# Copy this to .env in project root
APP_ENV=dev
LOG_LEVEL=info
PORT=8080
DB_DSN=file:./data/reserve_watch.db?_fk=1
FRED_API_KEY=b7cb42380ac4ab4708ff13b305755de5
DRY_RUN=true
AUTOPUBLISH=false
PUBLISH_LINKEDIN=false
PUBLISH_MAILCHIMP=false
```

### 2. Run locally:
```powershell
# PowerShell
cd C:\Users\unite\apps\reserve-watch
go run cmd/runner/main.go
```

### 3. Test in browser:
- Homepage: http://localhost:8080/
- Health: http://localhost:8080/health
- API: http://localhost:8080/api/latest
- Realtime API: http://localhost:8080/api/latest/realtime

---

## Development Workflow (FAST)

### Option 1: Quick Test Loop
```powershell
# 1. Make code changes
# 2. Run locally
go run cmd/runner/main.go

# 3. Test in browser (Ctrl+C when done)
# 4. Fix any errors
# 5. Repeat until working

# 6. ONLY THEN commit and push
git add .
git commit -m "Add feature X - tested locally"
git push origin main

# 7. Wait 3 min for Railway deploy
# 8. Verify on live site once
```

### Option 2: Auto-Reload Development
```powershell
# Install air for hot reload
go install github.com/cosmtrek/air@latest

# Run with auto-reload
air
# Now edits auto-restart the server!
```

---

## Testing Checklist (Before Deploy)

Run these locally:

```powershell
# 1. Build check
go build -o reserve-watch.exe cmd/runner/main.go

# 2. Test compile
go test ./...

# 3. Vet code
go vet ./...

# 4. Run app
go run cmd/runner/main.go

# 5. Test endpoints:
# - http://localhost:8080/
# - http://localhost:8080/health
# - http://localhost:8080/api/latest
# - http://localhost:8080/api/history
# - http://localhost:8080/api/latest/realtime
```

**If all pass → Commit + Push → Railway auto-deploys**

---

## Why This is Faster

| Approach | Time | Iterations |
|----------|------|-----------|
| **Deploy-Test Loop** | 3-5 min/change | Slow |
| **Local-Test Loop** | 10 sec/change | Fast |
| **Deploy when done** | 3-5 min once | Final check |

**10x faster development!**

---

## Common Local Testing

### Test Data Fetching:
```powershell
# Set FRED key
$env:FRED_API_KEY="b7cb42380ac4ab4708ff13b305755de5"

# Run
go run cmd/runner/main.go

# Check logs for:
# - "Fetching FRED series: DTWEXBGS"
# - "Fetching real-time DXY from Yahoo Finance..."
# - "Yahoo DXY: 98.9380"
```

### Test Web Interface:
```powershell
# Run app
go run cmd/runner/main.go

# Open browser to http://localhost:8080/
# Should see dashboard with data cards
```

### Debug Mode:
```powershell
# Verbose logging
$env:LOG_LEVEL="debug"
go run cmd/runner/main.go
```

---

## My New Workflow (AI Assistant)

### BEFORE:
1. Edit code
2. Commit + push
3. Wait 3-5 minutes for Railway
4. Check if it works
5. If broken, go to step 1 (repeat 10x = 30-50 minutes!)

### AFTER:
1. Edit code
2. Run locally (`go run cmd/runner/main.go`)
3. Test in browser (takes 10 seconds)
4. Fix errors immediately
5. Repeat steps 1-4 until working (10-15 minutes total)
6. **ONLY THEN** commit + push once
7. Wait 3-5 minutes for Railway
8. Quick verification on live site
9. DONE (total: 15-20 minutes vs 30-50 minutes)

---

## Install Air (Optional - Auto Reload)

```powershell
# Install
go install github.com/cosmtrek/air@latest

# Create .air.toml (config file)
# Then just run:
air

# Now every file save auto-restarts the server!
```

---

## Quick Commands

```powershell
# Run
go run cmd/runner/main.go

# Build
go build -o reserve-watch.exe cmd/runner/main.go

# Test
go test ./...

# Clean build
rm -r data/
go run cmd/runner/main.go

# Check for issues
go vet ./...
```

---

## Database Management

```powershell
# View local database
sqlite3 data/reserve_watch.db

# SQLite commands:
.tables                          # List tables
SELECT * FROM series_points;     # View data
SELECT * FROM posts;             # View posts
.exit                            # Exit
```

---

## Deployment Checklist

**Before pushing to Railway:**

- [ ] Runs locally without errors
- [ ] Homepage loads at http://localhost:8080/
- [ ] `/health` returns healthy status
- [ ] `/api/latest` returns data
- [ ] All data sources fetch successfully (check logs)
- [ ] No compilation errors (`go build`)
- [ ] Tests pass (`go test ./...`)

**If all ✅ → Push to GitHub → Railway deploys**

---

## Pro Tip: Test Script

Create `test-local.ps1`:
```powershell
# Quick local test script
Write-Host "Building..." -ForegroundColor Yellow
go build -o reserve-watch.exe cmd/runner/main.go

if ($LASTEXITCODE -eq 0) {
    Write-Host "Build successful!" -ForegroundColor Green
    Write-Host "Running tests..." -ForegroundColor Yellow
    go test ./...
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "Tests passed!" -ForegroundColor Green
        Write-Host "Starting server on http://localhost:8080" -ForegroundColor Cyan
        go run cmd/runner/main.go
    } else {
        Write-Host "Tests failed!" -ForegroundColor Red
    }
} else {
    Write-Host "Build failed!" -ForegroundColor Red
}
```

Run: `.\test-local.ps1`

---

**Bottom line: Develop locally (fast), deploy once (slow), verify (quick).**

