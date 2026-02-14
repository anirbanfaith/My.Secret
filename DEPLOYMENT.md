# 🚀 DEPLOYMENT GUIDE - Step by Step

## Quick Start (Get it Running in 5 Minutes!)

### On Your Computer (Local Testing)

1. **Open Terminal/Command Prompt**
   - Windows: Press `Win + R`, type `cmd`, press Enter
   - Mac: Press `Cmd + Space`, type `Terminal`, press Enter

2. **Navigate to Project**
   ```bash
   cd secret-vault
   ```

3. **Install Dependencies**
   ```bash
   go mod download
   ```

4. **Run the Server**
   ```bash
   go run main.go
   ```

5. **Open Your Browser**
   - Go to: `http://localhost:8080`
   - You should see the Secret Vault homepage!

---

## 🌐 Deploy to Internet (FREE Hosting)

### OPTION 1: Fly.io (Best for Go Apps) ⭐ RECOMMENDED

**Step 1: Install Fly CLI**

For Mac:
```bash
brew install flyctl
```

For Windows (PowerShell as Administrator):
```powershell
powershell -Command "iwr https://fly.io/install.ps1 -useb | iex"
```

For Linux:
```bash
curl -L https://fly.io/install.sh | sh
```

**Step 2: Sign Up**
```bash
flyctl auth signup
```
- Opens browser
- Create free account with email

**Step 3: Login**
```bash
flyctl auth login
```

**Step 4: Launch Your App**
```bash
cd secret-vault
flyctl launch
```

You'll be asked:
- **App name**: Choose a unique name (e.g., `my-secret-vault-2024`)
- **Region**: Choose closest to you
- **Database**: No (we're using SQLite)
- **Deploy now**: Yes

**Step 5: Deploy**
```bash
flyctl deploy
```

**Step 6: Open Your App**
```bash
flyctl open
```

✅ **Your app is now LIVE!** Share the URL with your girlfriend!

Your URL will be: `https://your-app-name.fly.dev`

---

### OPTION 2: Railway.app (Easiest!)

**Step 1: Create GitHub Repository (if you haven't)**
1. Go to https://github.com
2. Sign in/Sign up
3. Click "New Repository"
4. Name it: `secret-vault`
5. Make it private (optional, recommended for privacy)
6. Don't initialize with README (we already have one)

**Step 2: Push Code to GitHub**
```bash
cd secret-vault
git init
git add .
git commit -m "Initial commit - Secret Vault app"
git branch -M main
git remote add origin https://github.com/YOUR-USERNAME/secret-vault.git
git push -u origin main
```

**Step 3: Deploy on Railway**
1. Go to https://railway.app
2. Click "Start a New Project"
3. Select "Deploy from GitHub repo"
4. Authorize Railway to access your GitHub
5. Select the `secret-vault` repository
6. Railway will auto-detect Go and deploy!

**Step 4: Get Your URL**
- Click on your deployment
- Go to "Settings" → "Domains"
- Click "Generate Domain"
- Your app is live at: `your-app.up.railway.app`

---

### OPTION 3: Render.com

**Step 1: Push to GitHub** (same as Railway Option 2 above)

**Step 2: Deploy on Render**
1. Go to https://render.com
2. Sign up with GitHub
3. Click "New +" → "Web Service"
4. Select your `secret-vault` repository
5. Configure:
   - **Name**: secret-vault
   - **Environment**: Go
   - **Build Command**: `go build -o main`
   - **Start Command**: `./main`
   - **Plan**: Free
6. Click "Create Web Service"

Wait 2-3 minutes for deployment.

Your URL: `https://secret-vault.onrender.com`

---

## 🎯 After Deployment Checklist

✅ **Test the App**
1. Visit your live URL
2. Click "Create Vault"
3. Fill in test data
4. Save the Vault ID
5. Try unlocking it

✅ **Prepare for Your Girlfriend**
1. Create the REAL vault with your letter
2. Copy the Vault ID
3. Remember the vault key
4. Test unlocking once to make sure it works
5. Delete the test vault (optional)

✅ **Share with Her**
1. Send Vault ID via text/email
2. Tell her the vault key (in person or separate message)
3. Send her the website URL
4. Wait for her reaction! 💝

---

## 🔧 Troubleshooting

### "go: command not found"
**Problem**: Go is not installed
**Solution**: Download from https://golang.org/dl/

### "Port 8080 already in use"
**Problem**: Another app is using port 8080
**Solution**: 
```bash
PORT=3000 go run main.go
```
Then open `http://localhost:3000`

### "Cannot connect to database"
**Problem**: Database file is locked
**Solution**: 
1. Stop the server (Ctrl+C)
2. Delete `data/vaults.db`
3. Restart server

### Deployment fails
**Problem**: Build errors on hosting platform
**Solution**:
1. Make sure all files are committed
2. Check `go.mod` and `go.sum` are present
3. Try local build first: `go build`

### App is slow on free tier
**Problem**: Free hosting can be slow
**Solution**: 
- Fly.io: Might sleep after inactivity (wakes in ~5 seconds)
- Railway: 500 hours/month free
- Render: Spins down after 15 min inactivity
- **Tip**: Visit the site once before sharing with her to wake it up

---

## 💡 Pro Tips

### Make it Special
1. **Custom Domain**: Buy a domain like `ourspecialvault.com` for $12/year
2. **Custom Colors**: Edit `style.css` to match her favorite colors
3. **Add Photos**: Modify the letter to include `<img src="URL">` tags
4. **Voice Message**: Host audio file and link it in the letter

### Keep it Running
- Free tiers are perfect for this use case
- The vault will stay online as long as the hosting is active
- Export the database periodically to keep a backup

### Privacy
- Use a strong vault key (she should never guess it)
- Don't share the vault key publicly
- After she unlocks it, you can delete the vault for privacy

---

## 📊 Hosting Comparison

| Platform | Pros | Cons | Best For |
|----------|------|------|----------|
| **Fly.io** | Fast, Go-optimized, global | Sleeps if inactive | Production apps |
| **Railway** | Easiest setup, nice UI | 500hrs/month limit | Quick deploys |
| **Render** | Simple, reliable | Slower cold starts | Long-running apps |
| **Local** | Total control, free | Only accessible on your network | Testing |

---

## 🎊 You're Done!

Your Secret Vault is ready! 

**Next Steps:**
1. Test everything works
2. Create your romantic letter
3. Share the magic with her
4. Enjoy her reaction! 💝

**Remember**: This isn't just an app - it's a unique, thoughtful gift that shows you care enough to build something special just for her.

Good luck! 🍀
