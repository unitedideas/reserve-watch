# Reserve Watch - De-Dollarization Dashboard

Automated Go-based content engine that monitors de-dollarization trends—tracking the shift away from the U.S. dollar in global reserves and payments.

## Quick Start

### Prerequisites
- Go 1.22 or higher
- GCC compiler (for CGO - required by SQLite)
  - Linux: Usually pre-installed
  - macOS: Install Xcode Command Line Tools
  - Windows: Install [mingw-w64](https://www.mingw-w64.org/) or [TDM-GCC](https://jmeubank.github.io/tdm-gcc/)
- FRED API Key (get from https://fred.stlouisfed.org/docs/api/api_key.html)

### Installation

1. Clone the repository
```bash
git clone <repository-url>
cd reserve-watch
```

2. Install dependencies
```bash
go mod download
```

3. Configure environment variables
```bash
# Copy and edit the example
cp .env.example .env
```

Required environment variables:
```bash
FRED_API_KEY=your_fred_api_key_here
```

Optional variables for publishing:
```bash
LINKEDIN_ACCESS_TOKEN=
LINKEDIN_ORG_URN=
MAILCHIMP_API_KEY=
MAILCHIMP_SERVER_PREFIX=us1
MAILCHIMP_LIST_ID=
PUBLISH_LINKEDIN=false
PUBLISH_MAILCHIMP=false
AUTOPUBLISH=false
DRY_RUN=true
```

### Build

**On Linux/macOS:**
```bash
CGO_ENABLED=1 go build -o bin/reserve-watch cmd/runner/main.go
```

**On Windows (requires mingw-w64 or TDM-GCC):**
```bash
set CGO_ENABLED=1
go build -o bin/reserve-watch.exe cmd/runner/main.go
```

**Cross-compile for Linux from Windows:**
```bash
set GOOS=linux
set GOARCH=amd64
set CGO_ENABLED=0
go build -o bin/reserve-watch cmd/runner/main.go
```
Note: Cross-compilation with CGO disabled will require SQLite to be available on the target system.

### Run

```bash
./bin/reserve-watch
```

Or run directly:
```bash
go run cmd/runner/main.go
```

## Deployment to Production Server

### Option 1: Build and Deploy Binary

1. **Build for Linux (if building from Windows/Mac)**
```bash
GOOS=linux GOARCH=amd64 go build -o reserve-watch cmd/runner/main.go
```

2. **Transfer to server**
```bash
scp reserve-watch user@your-server:/opt/reserve-watch/
scp -r templates user@your-server:/opt/reserve-watch/
scp -r migrations user@your-server:/opt/reserve-watch/
```

3. **Set up environment on server**
```bash
ssh user@your-server
cd /opt/reserve-watch
nano .env  # Add your API keys
```

4. **Create systemd service**
```bash
sudo nano /etc/systemd/system/reserve-watch.service
```

Add:
```ini
[Unit]
Description=Reserve Watch Service
After=network.target

[Service]
Type=simple
User=your-user
WorkingDirectory=/opt/reserve-watch
ExecStart=/opt/reserve-watch/reserve-watch
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

5. **Enable and start**
```bash
sudo systemctl enable reserve-watch
sudo systemctl start reserve-watch
sudo systemctl status reserve-watch
```

### Option 2: Docker Deployment

Create `Dockerfile`:
```dockerfile
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=1 GOOS=linux go build -o reserve-watch cmd/runner/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates sqlite
WORKDIR /root/
COPY --from=builder /app/reserve-watch .
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/migrations ./migrations
CMD ["./reserve-watch"]
```

Build and run:
```bash
docker build -t reserve-watch .
docker run -d --name reserve-watch \
  -e FRED_API_KEY=your_key \
  -v $(pwd)/data:/root/data \
  reserve-watch
```

## Scheduled Jobs

The application runs a cron scheduler that checks FRED data daily at 9:00 AM:
- Fetches latest US Dollar Index data
- Detects changes and saves to database
- Generates charts and content
- Publishes to configured platforms (if enabled)

## Development

### Run linters
```bash
go vet ./...
go fmt ./...
```

### Run tests
```bash
# Run all tests (requires CGO)
CGO_ENABLED=1 go test -v ./...

# Run tests without SQLite/chart generation (no CGO required)
go test -v ./internal/config ./internal/publish
```

### Project Structure
```
/cmd/runner                 # Main entrypoint
/internal/config            # Environment configuration
/internal/ingest            # FRED data fetching
/internal/compose           # Content generation and charts
/internal/publish           # LinkedIn and Mailchimp publishers
/internal/store             # SQLite database layer
/internal/util              # Logging utilities
/templates                  # Content templates
/migrations                 # Database migrations
```

## Features

- ✅ FRED data ingestion
- ✅ SQLite storage with migrations
- ✅ Chart generation (PNG)
- ✅ Open Graph image generation
- ✅ LinkedIn publishing
- ✅ Mailchimp campaign creation
- ✅ Template-based content generation
- ✅ Cron scheduling
- ✅ Dry-run mode for testing

## Roadmap

See `README_ORIGINAL.md` for the full 14-day build plan including:
- Additional data sources (COFER, SWIFT RMB Tracker, CIPS, WGC)
- Indices calculation (RMB Penetration Score, Reserve Diversification Pressure)
- Alert webhooks
- YouTube integration
- Admin UI

## License

See LICENSE file for details.

