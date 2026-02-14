# 🔒 Secret Vault - Anniversary Gift App

A secure, romantic web application where you can create password-protected vaults containing secret letters for your loved ones.

## 🎁 Features

- **Secure Vault Creation**: Create vaults protected by strong keys
- **Security Questions**: Add 2 custom MCQ-style questions
- **Scoring System**: Points deducted for wrong answers (starts at 100)
- **Anti-Cheat Protection**: Multiple security measures to prevent cheating
- **Self-Destructing Vaults**: Creator can delete vaults
- **Beautiful UI**: Clean, romantic interface with animations

## 🔐 Security Features

1. **Password Hashing**: Uses Argon2id (industry-standard, more secure than bcrypt)
2. **Constant-time Comparison**: Prevents timing attacks
3. **Rate Limiting**: Limits attempts per IP address
4. **Case-insensitive Answers**: Answers normalized to prevent case issues
5. **No Client-side Secrets**: Answers and letters never sent to browser before unlock
6. **Strong Key Requirements**: Enforces uppercase, lowercase, numbers, and symbols
7. **Console Protection**: Prevents basic console-based cheating attempts
8. **Copy-Paste Disabled**: On answer fields to prevent cheating

## 📋 Prerequisites

Before you start, install these:

1. **Go Programming Language**
   - Download from: https://golang.org/dl/
   - Version 1.21 or higher
   - Installation guide: https://golang.org/doc/install

2. **Git** (for version control - optional but recommended)
   - Download from: https://git-scm.com/downloads

3. **Code Editor** (recommended)
   - VS Code: https://code.visualstudio.com/
   - Or any text editor you prefer

## 🚀 Local Setup (Run on Your Computer)

### Step 1: Set Up the Project

1. Open Terminal (Mac/Linux) or Command Prompt (Windows)
2. Navigate to the project folder:
```bash
cd secret-vault
```

### Step 2: Install Dependencies

```bash
go mod download
```

### Step 3: Run the Application

```bash
go run main.go
```

You should see:
```
🔒 Secret Vault server starting on port 8080...
```

### Step 4: Open in Browser

Open your browser and go to:
```
http://localhost:8080
```

That's it! The app is now running locally on your computer.

## 🌐 Deploy to the Internet (FREE!)

### Option 1: Fly.io (Recommended for Go apps)

1. **Install Fly CLI**
```bash
# Mac
brew install flyctl

# Windows (PowerShell)
powershell -Command "iwr https://fly.io/install.ps1 -useb | iex"

# Linux
curl -L https://fly.io/install.sh | sh
```

2. **Sign up for Fly.io**
```bash
fly auth signup
```

3. **Create fly.toml configuration file** (I'll create this for you below)

4. **Deploy**
```bash
fly launch
# Follow the prompts, accept defaults
# Choose a unique app name

fly deploy
```

Your app will be live at: `https://your-app-name.fly.dev`

### Option 2: Railway.app

1. Go to https://railway.app
2. Sign up with GitHub
3. Click "New Project" → "Deploy from GitHub repo"
4. Connect your GitHub account and select this repository
5. Railway will auto-detect Go and deploy!
6. Your app will be live at the provided URL

### Option 3: Render.com

1. Go to https://render.com
2. Sign up (free)
3. Click "New +" → "Web Service"
4. Connect your GitHub repository
5. Configure:
   - **Build Command**: `go build -o main`
   - **Start Command**: `./main`
6. Click "Create Web Service"

## 📁 Project Structure

```
secret-vault/
├── main.go              # Backend server (Go)
├── go.mod              # Go dependencies
├── static/             # Frontend files
│   ├── index.html      # Landing page
│   ├── create.html     # Create vault page
│   ├── unlock.html     # Unlock vault page
│   └── style.css       # All styling
└── data/               # Database (auto-created)
    └── vaults.db       # SQLite database
```

## 🎨 How to Use

### For the Creator (You):

1. Click "Create Vault"
2. Create a strong vault key (remember it!)
3. Set 2 security questions your girlfriend knows
4. Write your heartfelt letter
5. Click "Create Vault"
6. Share the **Vault ID** with her (via text, email, etc.)
7. Tell her the vault key separately (don't send both together!)

### For the Recipient (Your Girlfriend):

1. Open the website
2. Click "Open Vault"
3. Enter the Vault ID and Vault Key you received
4. Answer the 2 security questions
5. If correct → Read the beautiful letter! 💌
6. If wrong → Lose 10 points and try again

## ⚙️ Customization

### Change Starting Score
In `main.go`, find this line:
```go
score INTEGER DEFAULT 100,
```
Change `100` to any number you want.

### Change Attempts Limit
In `main.go`, find:
```go
attempts_left INTEGER DEFAULT 10,
```
Change `10` to your desired limit.

### Change Point Deduction
In `main.go`, find:
```go
newScore := score - 10
```
Change `10` to any deduction amount.

### Change Colors/Theme
Edit `static/style.css` at the top:
```css
:root {
    --primary: #6366f1;      /* Main color */
    --secondary: #ec4899;    /* Secondary color */
    --success: #10b981;      /* Success color */
    --danger: #ef4444;       /* Danger color */
}
```

## 🐛 Troubleshooting

### "Port already in use"
Change the port in `main.go` or set environment variable:
```bash
PORT=3000 go run main.go
```

### Database locked
Stop the server (Ctrl+C) and delete `data/vaults.db`, then restart.

### Dependencies not installing
Make sure you're in the project directory:
```bash
cd secret-vault
go mod tidy
```

## 📝 Important Notes

- **Vault IDs** are randomly generated 32-character strings
- **Answers are case-insensitive** (automatically converted to lowercase)
- **Keys are case-sensitive** for extra security
- **Database resets** when you redeploy (on free hosting)
- For permanent storage, use a hosted database (PostgreSQL on Render/Railway)

## 💝 Tips for Your Anniversary

1. **Create the vault a few days early** to test it works
2. **Send the Vault ID** via one channel (text message)
3. **Send the Key** via another channel (handwritten note, in person)
4. **Make the questions personal** - inside jokes, special moments
5. **Write from the heart** - she'll only see it after unlocking!

## 🎉 What Makes This Special

- It's **personalized** - only she can unlock it
- It's **gamified** - adds fun challenge to reading your letter
- It's **secure** - demonstrates you care about protecting your message
- It's **unique** - you built it yourself!
- It's a **keepsake** - she can screenshot and save the letter

## 🔒 Privacy & Security

- All passwords are hashed with Argon2id
- Letters are encrypted in the database
- No analytics or tracking
- Runs entirely on your infrastructure
- You control all the data

## 📞 Need Help?

If you run into issues:
1. Check the troubleshooting section above
2. Make sure all prerequisites are installed
3. Verify you're in the right directory
4. Check the terminal for error messages

## ❤️ Final Message

This is more than just code - it's a thoughtful, secure way to share your feelings. The effort you put into building and deploying this will mean so much more than a store-bought card.

Good luck with your anniversary! 🎊

---

**Made with ❤️ for your special someone**
