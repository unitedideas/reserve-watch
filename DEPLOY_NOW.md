# Deploy Your App NOW - 3 Easy Steps

## You Already Have the Railway Variables Page Open!

### Step 1: Update the FRED_API_KEY
In the Railway browser window:

1. **Find the field that says**: `FRED_API_KEY`
2. **It currently has**: `your_fred_api_key_here`
3. **Click in that field and replace with**: `b7cb42380ac4ab4708ff13b305755de5`
4. **Press Enter** or click out of the field

### Step 2: Click "Add All" (Optional but Recommended)
At the bottom of the suggested variables, click the **"Add All"** button
- This will add all the environment variables your app needs
- The FRED_API_KEY you just entered will be included

OR just scroll to the bottom and variables will auto-save.

### Step 3: Redeploy
The app should **automatically redeploy** when you save variables.

If it doesn't:
- Go back to the **"Deployments"** tab (at the top)
- Click the **"Deploy"** button

---

## ⏱️ What Happens Next (2-3 minutes)

1. ✅ Railway starts building your app (shows "Building...")
2. ✅ Compiles the Go code
3. ✅ Creates container
4. ✅ Deploys to: `https://web-production-4c1d00.up.railway.app`
5. ✅ **Your app is LIVE!** 🎉

---

## 🎯 Your FRED API Key

```
b7cb42380ac4ab4708ff13b305755de5
```

Copy this and paste it into the `FRED_API_KEY` field!

---

## 📺 Watch the Deployment

You can watch the logs by clicking:
- **"Logs"** at the top navigation
- Or **"View logs"** on the deployment card

You'll see:
- Building...
- Running migrations...
- Starting cron scheduler...
- **Success!**

---

## ✅ How to Know It's Working

Once deployed, your app will:
- Run daily at 9 AM to check FRED data
- Store data in SQLite database
- Generate charts when data changes
- Log everything to Railway logs

Check the logs to see: "Starting Reserve Watch..." and "Cron scheduler started"

---

## 🚀 DO THIS NOW:

1. Switch to Railway browser tab
2. Find `FRED_API_KEY` field
3. Replace with: `b7cb42380ac4ab4708ff13b305755de5`
4. Click "Add All"
5. Wait 2 minutes
6. **DONE!**

