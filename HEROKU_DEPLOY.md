# Deploying Reserve Watch to Heroku

## Important Notice About Heroku

⚠️ **Heroku Limitations for This App:**
- Heroku has an **ephemeral filesystem** - files are deleted on restart
- SQLite database will be **lost on every restart** (dyno sleep/restart)
- **Better alternatives**: Use Heroku Postgres add-on OR deploy to traditional VPS

### Recommended Solution
Use **Heroku Postgres** instead of SQLite. See "Database Migration" section below.

## Prerequisites

1. **Heroku Account** - Sign up at https://signup.heroku.com/
2. **Heroku CLI** - Install from https://devcenter.heroku.com/articles/heroku-cli
3. **Git** - Heroku deploys via Git
4. **FRED API Key** - Get from https://fred.stlouisfed.org/docs/api/api_key.html

## Quick Deployment Steps

### 1. Install Heroku CLI (if not installed)

**Windows:**
```powershell
# Download and run installer from:
# https://devcenter.heroku.com/articles/heroku-cli#download-and-install
```

**macOS:**
```bash
brew tap heroku/brew && brew install heroku
```

**Linux:**
```bash
curl https://cli-assets.heroku.com/install.sh | sh
```

### 2. Login to Heroku

```bash
heroku login
```

This will open a browser for authentication.

### 3. Initialize Git (if not already done)

```bash
# Check if already initialized
git status

# If not, initialize
git init
git add .
git commit -m "Initial commit for Heroku deployment"
```

### 4. Create Heroku App

```bash
# Create app with a specific name
heroku create reserve-watch-YOUR-NAME

# Or let Heroku generate a random name
heroku create

# This creates a Git remote named 'heroku'
```

### 5. Set Buildpack for Go with CGO

```bash
# Use Heroku Go buildpack that supports CGO
heroku buildpacks:set https://github.com/heroku/heroku-buildpack-go.git
```

### 6. Configure Environment Variables

```bash
# Required
heroku config:set FRED_API_KEY=your_fred_api_key_here

# App configuration
heroku config:set APP_ENV=production
heroku config:set LOG_LEVEL=info
heroku config:set DRY_RUN=false
heroku config:set AUTOPUBLISH=false

# Database (SQLite - will be lost on restart!)
heroku config:set DB_DSN="file:./data/reserve_watch.db?_fk=1"

# Publishing (optional)
heroku config:set PUBLISH_LINKEDIN=false
heroku config:set PUBLISH_MAILCHIMP=false

# LinkedIn (if using)
# heroku config:set LINKEDIN_ACCESS_TOKEN=your_token
# heroku config:set LINKEDIN_ORG_URN=your_urn

# Mailchimp (if using)
# heroku config:set MAILCHIMP_API_KEY=your_key
# heroku config:set MAILCHIMP_SERVER_PREFIX=us1
# heroku config:set MAILCHIMP_LIST_ID=your_list_id
```

### 7. Deploy to Heroku

```bash
# Push to Heroku
git push heroku main

# Or if you're on master branch
git push heroku master
```

### 8. Scale the Dyno

```bash
# Make sure at least one dyno is running
heroku ps:scale web=1
```

### 9. Open Your App

```bash
heroku open
```

### 10. View Logs

```bash
# Stream logs
heroku logs --tail

# View last 100 lines
heroku logs -n 100
```

## Database Migration to PostgreSQL (Recommended)

Since Heroku has ephemeral storage, use PostgreSQL:

### Install Postgres Add-on

```bash
# Add free Postgres database
heroku addons:create heroku-postgresql:mini

# Get database URL
heroku config:get DATABASE_URL
```

### Update Application for PostgreSQL

You'll need to modify `internal/store/sqlite.go` to support PostgreSQL:

1. Add PostgreSQL driver to `go.mod`:
```go
require (
    github.com/lib/pq v1.10.9
)
```

2. Update connection string handling in config:
```go
// Check if using Postgres
if strings.HasPrefix(cfg.DBDsn, "postgres://") {
    db, err := sql.Open("postgres", cfg.DBDsn)
} else {
    db, err := sql.Open("sqlite3", cfg.DBDsn)
}
```

3. Update migrations for PostgreSQL compatibility

### Set Postgres URL

```bash
heroku config:set DB_DSN=$(heroku config:get DATABASE_URL)
```

## Heroku Scheduler for Cron Jobs

Since the built-in cron won't work well on free dynos, use Heroku Scheduler:

```bash
# Add Scheduler add-on (free)
heroku addons:create scheduler:standard

# Open scheduler dashboard
heroku addons:open scheduler
```

In the dashboard:
- Add job: `/app/reserve-watch`
- Frequency: Daily at 9:00 AM (choose timezone)

## Important Heroku-Specific Notes

### 1. Free Dyno Sleep
- Free dynos sleep after 30 minutes of inactivity
- First request after sleep takes longer
- Solution: Upgrade to Hobby dyno ($7/month) or use a ping service

### 2. Ephemeral Filesystem
- **Problem**: SQLite data lost on restart/deploy
- **Solution**: Use PostgreSQL (recommended) or external storage

### 3. Port Binding
The app should listen on the PORT environment variable. Currently, it doesn't expose HTTP, so this is fine for background jobs.

If you add a web interface later:
```go
port := os.Getenv("PORT")
if port == "" {
    port = "8080"
}
http.ListenAndServe(":"+port, handler)
```

### 4. Build Time
Go apps with CGO take longer to build on Heroku (5-10 minutes).

## Cost Breakdown

### Free Tier
- 1 web dyno (sleeps after 30 min)
- Heroku Postgres Mini (10,000 rows)
- Heroku Scheduler
- **Cost**: $0/month

### Recommended Tier
- Hobby dyno (always on): $7/month
- Heroku Postgres Mini: $5/month
- Heroku Scheduler: Free
- **Total**: ~$12/month

## Troubleshooting

### Build Fails
```bash
# Check build logs
heroku logs --tail

# Common issues:
# - CGO not enabled: Check heroku.yml
# - Missing dependencies: Run go mod tidy
# - Go version: Check go.mod
```

### App Crashes
```bash
# Check logs
heroku logs --tail

# Restart dyno
heroku restart

# Check dyno status
heroku ps
```

### Database Issues
```bash
# Check database connection
heroku pg:info

# Access database
heroku pg:psql

# Check tables
heroku pg:psql -c "\dt"
```

### Environment Variables
```bash
# List all config
heroku config

# Set variable
heroku config:set VAR_NAME=value

# Unset variable
heroku config:unset VAR_NAME
```

## Commands Cheat Sheet

```bash
# Deploy
git push heroku main

# View logs
heroku logs --tail

# Restart app
heroku restart

# Scale dynos
heroku ps:scale web=1

# SSH into dyno
heroku run bash

# Run migrations manually
heroku run ./reserve-watch migrate

# Check running processes
heroku ps

# Open app
heroku open

# Check app info
heroku info

# Delete app
heroku apps:destroy --confirm app-name
```

## Alternative: Docker on Heroku

You can also deploy using Docker:

1. Create `heroku.yml`:
```yaml
build:
  docker:
    web: Dockerfile
```

2. Set stack to container:
```bash
heroku stack:set container
```

3. Deploy:
```bash
git push heroku main
```

## Complete Deployment Script

Here's a complete script to deploy:

```bash
#!/bin/bash

# Login to Heroku
heroku login

# Create app
heroku create reserve-watch-production

# Set buildpack
heroku buildpacks:set https://github.com/heroku/heroku-buildpack-go.git

# Configure environment
heroku config:set FRED_API_KEY=your_key_here
heroku config:set APP_ENV=production
heroku config:set LOG_LEVEL=info
heroku config:set DRY_RUN=false
heroku config:set AUTOPUBLISH=false

# Add PostgreSQL (recommended)
heroku addons:create heroku-postgresql:mini

# Add Scheduler
heroku addons:create scheduler:standard

# Deploy
git add .
git commit -m "Deploy to Heroku"
git push heroku main

# Scale
heroku ps:scale web=1

# View logs
heroku logs --tail
```

## Next Steps After Deployment

1. **Set up Scheduler**:
   ```bash
   heroku addons:open scheduler
   ```
   Add daily job at 9 AM

2. **Monitor logs**:
   ```bash
   heroku logs --tail
   ```

3. **Check database**:
   ```bash
   heroku pg:psql -c "SELECT * FROM series_points LIMIT 5;"
   ```

4. **Test manually**:
   ```bash
   heroku run ./reserve-watch
   ```

5. **Upgrade to Hobby dyno** (if needed):
   ```bash
   heroku ps:resize web=hobby
   ```

## Support

- Heroku Dev Center: https://devcenter.heroku.com/
- Heroku CLI: https://devcenter.heroku.com/articles/heroku-cli
- Go on Heroku: https://devcenter.heroku.com/articles/getting-started-with-go

## Warning: SQLite on Heroku

Remember that SQLite **will not persist** on Heroku's free tier. Every time your dyno restarts (which happens at least once per day), you'll lose your database.

**Strongly recommend** migrating to PostgreSQL before production use!


