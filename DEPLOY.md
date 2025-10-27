# Deployment Guide - Reserve Watch

## Overview
This guide covers deploying the Reserve Watch De-Dollarization Dashboard to a production Linux server.

## Prerequisites on Server
- Ubuntu/Debian Linux 20.04+ (or similar)
- Root or sudo access
- Domain name pointed to server IP (optional but recommended)

## Step-by-Step Deployment

### 1. Prepare the Server

```bash
# Update system
sudo apt update && sudo apt upgrade -y

# Install required packages
sudo apt install -y sqlite3 gcc git

# Create application user
sudo useradd -r -s /bin/bash -d /opt/reserve-watch reserve-watch
sudo mkdir -p /opt/reserve-watch
sudo chown reserve-watch:reserve-watch /opt/reserve-watch
```

### 2. Build the Application

**Option A: Build on server**
```bash
# Install Go on server
wget https://go.dev/dl/go1.22.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin

# Clone and build
cd /opt/reserve-watch
git clone <your-repo-url> .
CGO_ENABLED=1 go build -o reserve-watch cmd/runner/main.go
```

**Option B: Cross-compile and upload (from dev machine)**
```bash
# On your dev machine
GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o reserve-watch cmd/runner/main.go

# Upload to server
scp reserve-watch user@your-server:/opt/reserve-watch/
scp -r templates user@your-server:/opt/reserve-watch/
scp -r migrations user@your-server:/opt/reserve-watch/
```

### 3. Configure Environment

```bash
sudo -u reserve-watch bash
cd /opt/reserve-watch

# Create .env file
cat > .env << 'EOF'
APP_ENV=production
LOG_LEVEL=info
DB_DSN=file:/opt/reserve-watch/data/reserve_watch.db?_fk=1

# FRED API (Required)
FRED_API_KEY=your_actual_fred_api_key_here

# LinkedIn (Optional)
LINKEDIN_ACCESS_TOKEN=
LINKEDIN_ORG_URN=

# Mailchimp (Optional)
MAILCHIMP_API_KEY=
MAILCHIMP_SERVER_PREFIX=us1
MAILCHIMP_LIST_ID=

# Publishing
PUBLISH_LINKEDIN=false
PUBLISH_MAILCHIMP=false
AUTOPUBLISH=false
DRY_RUN=true
EOF

# Create directories
mkdir -p data output
chmod 600 .env
```

### 4. Create Systemd Service

```bash
sudo cat > /etc/systemd/system/reserve-watch.service << 'EOF'
[Unit]
Description=Reserve Watch - De-Dollarization Dashboard
After=network.target

[Service]
Type=simple
User=reserve-watch
Group=reserve-watch
WorkingDirectory=/opt/reserve-watch
ExecStart=/opt/reserve-watch/reserve-watch
Restart=always
RestartSec=10
StandardOutput=journal
StandardError=journal

# Security settings
NoNewPrivileges=true
PrivateTmp=true
ReadOnlyPaths=/etc
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/opt/reserve-watch/data /opt/reserve-watch/output

[Install]
WantedBy=multi-user.target
EOF
```

### 5. Start the Service

```bash
# Enable and start
sudo systemctl daemon-reload
sudo systemctl enable reserve-watch
sudo systemctl start reserve-watch

# Check status
sudo systemctl status reserve-watch

# View logs
sudo journalctl -u reserve-watch -f
```

### 6. Setup Nginx Reverse Proxy (Optional)

If you want to expose a web interface later:

```bash
sudo apt install -y nginx

sudo cat > /etc/nginx/sites-available/reserve-watch << 'EOF'
server {
    listen 80;
    server_name your-domain.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
EOF

sudo ln -s /etc/nginx/sites-available/reserve-watch /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

### 7. Setup SSL with Let's Encrypt (Optional)

```bash
sudo apt install -y certbot python3-certbot-nginx
sudo certbot --nginx -d your-domain.com
```

## Monitoring and Maintenance

### Check Application Status
```bash
sudo systemctl status reserve-watch
```

### View Logs
```bash
# Real-time logs
sudo journalctl -u reserve-watch -f

# Last 100 lines
sudo journalctl -u reserve-watch -n 100

# Logs from specific time
sudo journalctl -u reserve-watch --since "2024-01-01" --until "2024-01-02"
```

### Restart Service
```bash
sudo systemctl restart reserve-watch
```

### Update Application
```bash
sudo systemctl stop reserve-watch

# Update code (if using git)
cd /opt/reserve-watch
sudo -u reserve-watch git pull
sudo -u reserve-watch CGO_ENABLED=1 go build -o reserve-watch cmd/runner/main.go

# Or upload new binary
# scp reserve-watch user@server:/opt/reserve-watch/

sudo systemctl start reserve-watch
```

### Backup Database
```bash
# Create backup
sudo -u reserve-watch sqlite3 /opt/reserve-watch/data/reserve_watch.db ".backup '/opt/reserve-watch/data/backup-$(date +%Y%m%d).db'"

# Setup daily backup cron
sudo -u reserve-watch crontab -e
# Add: 0 2 * * * sqlite3 /opt/reserve-watch/data/reserve_watch.db ".backup '/opt/reserve-watch/data/backup-$(date +\%Y\%m\%d).db'"
```

## Docker Deployment (Alternative)

### Build and Run with Docker

```dockerfile
# Dockerfile
FROM golang:1.22-alpine AS builder
RUN apk add --no-cache gcc musl-dev sqlite-dev
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o reserve-watch cmd/runner/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates sqlite-libs
WORKDIR /app
COPY --from=builder /app/reserve-watch .
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/migrations ./migrations
RUN mkdir -p /app/data /app/output
CMD ["./reserve-watch"]
```

```bash
# Build
docker build -t reserve-watch:latest .

# Run
docker run -d \
  --name reserve-watch \
  --restart unless-stopped \
  -v $(pwd)/data:/app/data \
  -v $(pwd)/output:/app/output \
  -e FRED_API_KEY=your_key \
  -e DRY_RUN=true \
  reserve-watch:latest

# View logs
docker logs -f reserve-watch
```

## Troubleshooting

### Service won't start
```bash
# Check service status
sudo systemctl status reserve-watch

# Check logs for errors
sudo journalctl -u reserve-watch -n 50 --no-pager

# Verify binary permissions
ls -la /opt/reserve-watch/reserve-watch

# Test manually
sudo -u reserve-watch /opt/reserve-watch/reserve-watch
```

### Database errors
```bash
# Check database file permissions
ls -la /opt/reserve-watch/data/

# Verify migrations
sqlite3 /opt/reserve-watch/data/reserve_watch.db ".schema"

# Check database integrity
sqlite3 /opt/reserve-watch/data/reserve_watch.db "PRAGMA integrity_check;"
```

### API connection issues
```bash
# Test FRED API key
curl "https://api.stlouisfed.org/fred/series/observations?series_id=DTWEXBGS&api_key=YOUR_KEY&file_type=json&limit=1"

# Check network connectivity
ping api.stlouisfed.org
```

## Security Best Practices

1. **API Keys**: Store in `.env` file with 600 permissions
2. **Firewall**: Configure UFW or iptables
3. **Updates**: Regularly update system and dependencies
4. **Monitoring**: Set up monitoring for service uptime
5. **Backups**: Regular database backups to external storage
6. **Logs**: Rotate logs to prevent disk space issues

## Performance Tuning

For high-frequency updates:
```bash
# Increase file descriptors
echo "reserve-watch soft nofile 4096" | sudo tee -a /etc/security/limits.conf
echo "reserve-watch hard nofile 8192" | sudo tee -a /etc/security/limits.conf
```

## Support

For issues, check:
1. Application logs: `sudo journalctl -u reserve-watch`
2. System logs: `sudo dmesg`
3. Database integrity
4. API connectivity

