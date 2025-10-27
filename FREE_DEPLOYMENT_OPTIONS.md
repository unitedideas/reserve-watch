# Free Deployment Options for Reserve Watch

## 🏆 Best Free Options (Ranked)

### 1. Railway.app ⭐ **RECOMMENDED**
**Why:** Easiest setup, generous free tier, great Go support, persistent storage

**Free Tier:**
- $5/month credit (enough for small apps)
- 500 hours/month execution
- 1GB RAM, 1GB storage
- Persistent volumes (database won't be lost!)

**Deploy in 3 Minutes:**

```bash
# 1. Install Railway CLI
npm install -g @railway/cli

# Or using PowerShell
iwr https://railway.app/install.ps1 | iex

# 2. Login
railway login

# 3. Initialize project
railway init

# 4. Set environment variables
railway variables set FRED_API_KEY=your_key_here

# 5. Deploy
railway up

# 6. View logs
railway logs
```

**Or Deploy via Web (No CLI needed):**
1. Go to https://railway.app/new
2. Connect your GitHub repo
3. Railway auto-detects Go and deploys
4. Add FRED_API_KEY in settings
5. Done! ✅

---

### 2. Render.com 🚀
**Why:** Free tier is truly free (no credit card), good persistent storage

**Free Tier:**
- Free forever (with credit card)
- 750 hours/month
- Sleeps after 15 min inactivity
- Persistent disk (database survives restarts!)

**Deploy:**

```bash
# 1. Go to https://render.com
# 2. Sign up (free)
# 3. Click "New +" → "Background Worker"
# 4. Connect GitHub repo
# 5. Render auto-detects using render.yaml
# 6. Add FRED_API_KEY environment variable
# 7. Deploy!
```

**Manual Setup:**
1. Build Command: `go build -o reserve-watch cmd/runner/main.go`
2. Start Command: `./reserve-watch`
3. Environment: Go
4. Plan: Free

---

### 3. Fly.io ✈️
**Why:** Great for Go apps, real persistent storage, global CDN

**Free Tier:**
- 3 shared-cpu-1x VMs
- 3GB persistent storage
- 160GB bandwidth/month

**Deploy:**

```bash
# 1. Install flyctl
# PowerShell
iwr https://fly.io/install.ps1 -useb | iex

# 2. Login
fly auth signup  # or fly auth login

# 3. Launch app
fly launch --no-deploy

# 4. Create volume for database
fly volumes create data --region ord --size 1

# 5. Set secrets
fly secrets set FRED_API_KEY=your_key_here

# 6. Deploy
fly deploy

# 7. Check status
fly status

# 8. View logs
fly logs
```

---

### 4. DigitalOcean App Platform
**Why:** $200 free credit for 60 days (requires credit card)

**Free Trial:**
- $200 credit (60 days)
- Then $5/month for basic plan
- Full featured, production ready

**Deploy:**
1. Go to https://cloud.digitalocean.com/apps
2. Click "Create App"
3. Connect GitHub
4. Select repository
5. DigitalOcean auto-configures
6. Add environment variables
7. Deploy

---

### 5. Google Cloud Run 🌐
**Why:** Very generous free tier, scales to zero

**Free Tier:**
- 2 million requests/month
- 360,000 GB-seconds memory
- 180,000 vCPU-seconds

**Deploy:**

```bash
# 1. Install gcloud CLI
# Download from: https://cloud.google.com/sdk/docs/install

# 2. Login
gcloud auth login

# 3. Create project
gcloud projects create reserve-watch-$(date +%s)

# 4. Build container
gcloud builds submit --tag gcr.io/PROJECT_ID/reserve-watch

# 5. Deploy
gcloud run deploy reserve-watch \
  --image gcr.io/PROJECT_ID/reserve-watch \
  --platform managed \
  --region us-central1 \
  --set-env-vars FRED_API_KEY=your_key

# Note: May need Dockerfile (I can create one)
```

---

## ⚡ Fastest Option: Railway (No Setup)

Since you want something NOW, use Railway:

### Ultra-Quick Railway Deploy:

1. **Go to**: https://railway.app/new
2. **Click**: "Deploy from GitHub repo"
3. **Select**: your reserve-watch repo (or upload code)
4. **Railway auto-detects Go** and deploys
5. **Add environment variable**:
   - Settings → Variables → Add `FRED_API_KEY`
6. **Done!** Your app is live in ~2 minutes

Railway will:
- ✅ Auto-build your Go app
- ✅ Run it continuously
- ✅ Provide persistent storage
- ✅ Auto-restart on failure
- ✅ Give you logs and metrics

---

## 📊 Comparison Table

| Platform | Free Tier | DB Storage | Setup Time | Best For |
|----------|-----------|------------|------------|----------|
| **Railway** | $5 credit/mo | ✅ Persistent | 2 min | Easiest |
| **Render** | 750 hrs/mo | ✅ Persistent | 5 min | No CC |
| **Fly.io** | 3 VMs | ✅ Persistent | 10 min | Control |
| **DO App** | $200/60d | ✅ Persistent | 10 min | Trial |
| **GCP Run** | 2M req/mo | ❌ Ephemeral* | 15 min | Scale |
| **Heroku** | 550 hrs/mo | ❌ Ephemeral | 10 min | Legacy |

*GCP Run needs Cloud SQL for persistence

---

## 🎯 My Recommendation: Railway

**Why Railway wins for your use case:**

1. **Literally 2 minutes to deploy**
2. **No credit card needed** (for $5/month credit)
3. **Database persists** (no data loss)
4. **Auto-deploys** from Git
5. **Great Go/CGO support**
6. **Simple environment variables**
7. **Good free tier** ($5 credit = ~150 hours)

### Deploy to Railway RIGHT NOW:

```bash
# If you have Node.js installed:
npx @railway/cli login
npx @railway/cli init
npx @railway/cli up

# Or just use the website (even easier):
# 1. https://railway.app/new
# 2. Import from GitHub
# 3. Add FRED_API_KEY
# 4. Deploy button
# Done!
```

---

## 🔧 If You Don't Have Git Installed

I noticed you don't have Git. Here's how to deploy without it:

### Option 1: GitHub Desktop (Easiest)
1. Download: https://desktop.github.com/
2. Create repo from this folder
3. Push to GitHub
4. Connect Railway/Render to GitHub repo
5. Auto-deploys!

### Option 2: Direct Upload (Railway)
1. Go to railway.app
2. Click "New Project"
3. Choose "Empty Project"
4. Railway CLI: `railway up` (uploads directly)

### Option 3: I'll create a deployment ZIP
Would you like me to create a deployment package you can upload directly?

---

## 💾 Database Consideration

**Important:** Your app uses SQLite, which works great on:
- ✅ Railway (persistent volumes)
- ✅ Render (persistent disk)
- ✅ Fly.io (persistent volumes)
- ❌ Heroku free tier (ephemeral!)
- ❌ Google Cloud Run (needs Cloud SQL)

---

## 🚀 What Do You Want to Do?

**Option A: Fastest (2 minutes)**
```
Use Railway web interface - no CLI needed
```

**Option B: CLI Deploy**
```
Pick Railway, Render, or Fly - I'll guide you
```

**Option C: Need Git First**
```
Install GitHub Desktop, then deploy
```

**Option D: Create Deployment Package**
```
I'll zip everything for manual upload
```

Which option sounds best? Or just tell me your FRED API key and preferred platform, and I'll give you the exact commands!


