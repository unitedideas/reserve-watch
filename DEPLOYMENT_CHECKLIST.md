# Deployment Checklist for Reserve Watch

## Pre-Deployment Requirements

### 1. Infrastructure Ready? ✓/✗
- [ ] Linux server provisioned (Ubuntu 20.04+ recommended)
- [ ] SSH access configured
- [ ] Domain name pointed to server IP (optional but recommended)
- [ ] Firewall rules allow ports 80 and 443

### 2. API Keys & Credentials ✓/✗
- [ ] **FRED API Key** (REQUIRED) - Get from https://fred.stlouisfed.org/docs/api/api_key.html
- [ ] LinkedIn Access Token (optional) - Only if publishing to LinkedIn
- [ ] Mailchimp API Key (optional) - Only if publishing to Mailchimp

### 3. Local Preparation ✓/✗
- [ ] Application tested locally with DRY_RUN=true
- [ ] Binary built for Linux
- [ ] Templates directory ready
- [ ] Migrations directory ready
- [ ] .env file prepared with credentials

## Quick Deploy Commands

### Option 1: Deploy from Windows to Linux Server

```powershell
# 1. Build for Linux (if on Windows without CGO)
# Note: This may require building on the Linux server itself due to SQLite/CGO
$env:GOOS="linux"
$env:GOARCH="amd64"
$env:CGO_ENABLED="1"
go build -o reserve-watch cmd/runner/main.go

# 2. Prepare deployment package
mkdir deploy-package
copy reserve-watch deploy-package/
xcopy /E /I templates deploy-package/templates
xcopy /E /I migrations deploy-package/migrations

# 3. Create .env file
@"
APP_ENV=production
LOG_LEVEL=info
DB_DSN=file:/opt/reserve-watch/data/reserve_watch.db?_fk=1
FRED_API_KEY=YOUR_FRED_API_KEY_HERE
DRY_RUN=false
AUTOPUBLISH=false
PUBLISH_LINKEDIN=false
PUBLISH_MAILCHIMP=false
"@ | Out-File -Encoding utf8 deploy-package/.env

# 4. Upload to server
scp -r deploy-package user@your-server:/tmp/reserve-watch
```

### Option 2: Build on Linux Server (Recommended)

```bash
# SSH into your server
ssh user@your-server

# Install Go and dependencies
sudo apt update
sudo apt install -y golang-go gcc sqlite3 git

# Clone or upload code
cd /opt
sudo mkdir reserve-watch
sudo chown $USER:$USER reserve-watch
cd reserve-watch

# Upload code here or use git

# Build on server
CGO_ENABLED=1 go build -o reserve-watch cmd/runner/main.go

# Create .env file
nano .env
# Add your credentials

# Set up systemd service (see DEPLOY.md)
sudo cp /path/to/systemd/file /etc/systemd/system/reserve-watch.service
sudo systemctl daemon-reload
sudo systemctl enable reserve-watch
sudo systemctl start reserve-watch
```

## Deployment Steps

### Step 1: Verify Prerequisites
```bash
# Check server connectivity
ping your-server-ip

# Check SSH access
ssh user@your-server "uname -a"
```

### Step 2: Prepare Server
```bash
ssh user@your-server

# Update system
sudo apt update && sudo apt upgrade -y

# Install required packages
sudo apt install -y gcc sqlite3 nginx certbot python3-certbot-nginx

# Create application directory
sudo mkdir -p /opt/reserve-watch/{data,output}
sudo useradd -r -s /bin/bash -d /opt/reserve-watch reserve-watch
sudo chown -R reserve-watch:reserve-watch /opt/reserve-watch
```

### Step 3: Deploy Application

**Option A: Build on Server**
```bash
# Install Go
wget https://go.dev/dl/go1.22.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin

# Upload source code
scp -r . user@your-server:/opt/reserve-watch/

# Build
cd /opt/reserve-watch
CGO_ENABLED=1 go build -o reserve-watch cmd/runner/main.go
```

**Option B: Upload Pre-built Binary**
```bash
# From local machine
scp reserve-watch user@your-server:/opt/reserve-watch/
scp -r templates user@your-server:/opt/reserve-watch/
scp -r migrations user@your-server:/opt/reserve-watch/
```

### Step 4: Configure Environment
```bash
ssh user@your-server
cd /opt/reserve-watch

# Create .env file
cat > .env << 'EOF'
APP_ENV=production
LOG_LEVEL=info
DB_DSN=file:/opt/reserve-watch/data/reserve_watch.db?_fk=1
FRED_API_KEY=your_actual_fred_api_key
DRY_RUN=false
AUTOPUBLISH=false
PUBLISH_LINKEDIN=false
PUBLISH_MAILCHIMP=false
EOF

sudo chmod 600 .env
sudo chown reserve-watch:reserve-watch .env
```

### Step 5: Create Systemd Service
```bash
sudo tee /etc/systemd/system/reserve-watch.service > /dev/null << 'EOF'
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

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl daemon-reload
sudo systemctl enable reserve-watch
```

### Step 6: Start Service
```bash
sudo systemctl start reserve-watch
sudo systemctl status reserve-watch

# Watch logs
sudo journalctl -u reserve-watch -f
```

### Step 7: Verify Deployment
```bash
# Check if process is running
sudo systemctl status reserve-watch

# Check logs for errors
sudo journalctl -u reserve-watch -n 50 --no-pager

# Verify database was created
ls -la /opt/reserve-watch/data/

# Check output directory
ls -la /opt/reserve-watch/output/

# Test manually (optional)
sudo -u reserve-watch /opt/reserve-watch/reserve-watch
```

## Post-Deployment

### Monitor First Run
```bash
# Watch logs in real-time
sudo journalctl -u reserve-watch -f

# Check for:
# - Successful FRED API connection
# - Database initialization
# - Cron scheduler started
# - First data fetch (may take until 9 AM)
```

### Verify Data Flow
```bash
# Check database
sqlite3 /opt/reserve-watch/data/reserve_watch.db "SELECT * FROM series_points ORDER BY ingested_at DESC LIMIT 5;"

# Check generated content
ls -la /opt/reserve-watch/output/

# Check posts (if publishing enabled)
sqlite3 /opt/reserve-watch/data/reserve_watch.db "SELECT * FROM posts ORDER BY published_at DESC LIMIT 5;"
```

### Setup HTTPS (if web interface needed)
```bash
# Install Nginx
sudo apt install -y nginx certbot python3-certbot-nginx

# Get SSL certificate
sudo certbot --nginx -d your-domain.com
```

## Troubleshooting

### Service won't start
```bash
# Check service status
sudo systemctl status reserve-watch

# Check logs
sudo journalctl -u reserve-watch -n 100 --no-pager

# Check permissions
ls -la /opt/reserve-watch/

# Test manually
sudo -u reserve-watch /opt/reserve-watch/reserve-watch
```

### Database errors
```bash
# Check database file
ls -la /opt/reserve-watch/data/*.db

# Check permissions
sudo chown -R reserve-watch:reserve-watch /opt/reserve-watch/data/

# Verify migrations
sqlite3 /opt/reserve-watch/data/reserve_watch.db ".schema"
```

### API connection issues
```bash
# Test FRED API
curl "https://api.stlouisfed.org/fred/series/observations?series_id=DTWEXBGS&api_key=YOUR_KEY&file_type=json&limit=1"

# Check network
ping api.stlouisfed.org
```

## Rollback Plan

If something goes wrong:
```bash
# Stop service
sudo systemctl stop reserve-watch

# Restore previous binary (if you kept a backup)
sudo cp /opt/reserve-watch/reserve-watch.backup /opt/reserve-watch/reserve-watch

# Restore database (if you have backup)
sudo cp /opt/reserve-watch/data/backup.db /opt/reserve-watch/data/reserve_watch.db

# Start service
sudo systemctl start reserve-watch
```

## What Information Do You Need?

Before deploying, please provide:

1. **Server Information:**
   - IP address or hostname
   - SSH username
   - SSH key path (or password)

2. **API Keys:**
   - FRED API Key (required)
   - LinkedIn credentials (if publishing)
   - Mailchimp credentials (if publishing)

3. **Deployment Preference:**
   - [ ] Build on server (recommended)
   - [ ] Upload pre-built binary
   - [ ] Use Docker

4. **Domain/DNS (optional):**
   - Domain name (if you have one)
   - DNS configured? (Y/N)

Once you provide this information, I can create the specific deployment commands for your setup!

