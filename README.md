# 💝 My Secret - A Romantic Secret Message Vault

A beautiful web application for creating and sharing encrypted secret messages with someone special. Built with Go and designed with a minimalist, romantic aesthetic inspired by digibouquet.

![My Secret](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/license-MIT-green)
![Status](https://img.shields.io/badge/status-active-success)

## ✨ Features

### 🎯 Core Features
- **Link-Based Access** - Each vault gets a unique shareable URL
- **Challenge Questions** - Protect messages with 2 custom security questions
- **Per-Person Attempts** - Each person gets 5 independent attempts (tracked by name)
- **Multiple Winners** - Everyone who answers correctly can read the message
- **Live Leaderboard** - See who tried and their scores in real-time
- **Score System** - Points based on attempts: 100 → 80 → 60 → 40 → 20

### 🎨 Design
- **Digibouquet-Inspired Aesthetic** - Clean, minimalist, romantic
- **Typography** - Libre Baskerville (serif) + Courier New (monospace)
- **Color Palette** - Cream (#f5f0e8), White (#fefcf7), Black borders
- **Responsive** - Works on mobile, tablet, and desktop
- **No Database Required** - Uses simple JSON file storage

### 🔒 Privacy & Security
- **Case-Insensitive Answers** - Automatic normalization
- **No Copy-Paste** - Answer fields protected from copying
- **Rate Limiting** - 5 attempts per name prevents brute force
- **Local Storage** - Data stored in `data.json` file

---

## 🚀 Quick Start

### Prerequisites
- Go 1.21 or higher
- Git

### Local Development

```bash
# Clone the repository
git clone https://github.com/yourusername/My.Secret.git
cd My.Secret

# Run the server
go run main.go

# Open in browser
http://localhost:3000
```

That's it! The app will create a `data.json` file automatically.

---

## 📖 How It Works

### For Creators

1. **Visit the homepage** → Click "Create Secret"
2. **Fill in the form:**
   - Question 1 & Answer 1 (e.g., "What's my nickname?" → "Honey")
   - Question 2 & Answer 2 (e.g., "Where did we meet?" → "Paris")
   - Write your secret letter (handwritten-style text area)
3. **Get a unique link** → Share with your special someone
4. **Track attempts** → See who tried via the leaderboard

### For Recipients

1. **Open the secret link** (e.g., `yoursite.com/v/abc123xyz`)
2. **Check leaderboard** (optional) → See who already tried
3. **Enter your name** → Get 5 attempts
4. **Answer the questions** → Tries left shown in real-time
5. **Success!** → Read the secret letter + see your score
6. **View leaderboard** → See your ranking

---

## 🏗️ Project Structure

```
My.Secret/
├── main.go              # Go server (API + routing)
├── go.mod               # Go dependencies
├── data.json            # Auto-created data storage
├── static/
│   ├── style.css        # Beautiful digibouquet styling
│   ├── index.html       # Homepage
│   ├── create.html      # Create vault form
│   └── vault.html       # Unlock vault page
├── .gitignore          # Git ignore rules
└── README.md           # This file
```

---

## 🌐 Deployment

### Railway (Recommended - Free Tier)

```bash
# 1. Push to GitHub
git init
git add .
git commit -m "Initial commit"
git remote add origin https://github.com/yourusername/My.Secret.git
git push -u origin main

# 2. Deploy on Railway
# - Go to railway.app
# - Sign in with GitHub
# - New Project → Deploy from GitHub
# - Select your repository
# - Done! Auto-deploys on every push
```

**Free Tier:** $5 credit/month (enough for ~500 hours)

### Other Platforms

**Render.com** (Free with sleep mode)
```bash
# Build: go build -o main
# Start: ./main
```

**Fly.io** (Best free tier)
```bash
fly launch
fly deploy
```

**Environment Variables:**
- `PORT` - Automatically set by hosting platforms
- No manual configuration needed!

---

## 🎯 Use Cases

### Personal
- 💌 **Anniversary messages** - Create yearly vaults
- 🎂 **Birthday surprises** - Time-release secrets
- 💍 **Proposals** - Make it memorable
- 🎁 **Gift reveals** - Tease before the big day

### Creative
- 🎮 **Treasure hunts** - Clues in secret vaults
- 📚 **Story games** - Interactive narratives
- 🎓 **Education** - Gamified learning challenges
- 🎭 **Events** - Secret invitations with puzzles

### Professional
- 🔐 **Secure sharing** - Password-protected messages
- 🎉 **Team building** - Fun office games
- 📢 **Announcements** - Build anticipation

---

## 🛠️ API Endpoints

### Create Vault
```http
POST /api/create
Content-Type: application/json

{
  "question1": "What's my nickname?",
  "answer1": "honey",
  "question2": "Where did we meet?",
  "answer2": "paris",
  "letter": "Dear love, I wanted to tell you..."
}

Response: 200 OK
{
  "vault_id": "abc123xyz...",
  "vault_url": "https://yoursite.com/v/abc123xyz"
}
```

### Get Vault Info
```http
GET /api/vault/{vault_id}

Response: 200 OK
{
  "vault_id": "abc123xyz",
  "question1": "What's my nickname?",
  "question2": "Where did we meet?"
}
```

### Check Attempts
```http
GET /api/check-attempts?vault_id={id}&name={name}

Response: 200 OK
{
  "attempts_used": 2,
  "attempts_left": 3,
  "can_try": true
}
```

### Unlock Vault
```http
POST /api/unlock
Content-Type: application/json

{
  "vault_id": "abc123xyz",
  "name": "John",
  "answer1": "honey",
  "answer2": "paris"
}

Response: 200 OK (Success)
{
  "success": true,
  "letter": "Dear love, I wanted to tell you...",
  "score": 80
}

Response: 401 Unauthorized (Wrong Answer)
{
  "success": false,
  "attempts_left": 2,
  "max_reached": false
}
```

### Get Leaderboard
```http
GET /api/leaderboard?vault_id={id}

Response: 200 OK
[
  {
    "name": "John",
    "score": 100,
    "success": true,
    "created_at": "2026-02-15T05:00:00Z"
  },
  {
    "name": "Sarah",
    "score": 0,
    "success": false,
    "created_at": "2026-02-15T05:05:00Z"
  }
]
```

---

## ⚙️ Configuration

### Change Port (Local)
```bash
PORT=8080 go run main.go
```

### Storage Location
Data is stored in `data.json` in the project root. To backup:
```bash
cp data.json data.backup.json
```

### Customize Styling
Edit `static/style.css` to change:
- Colors (search for hex codes)
- Fonts (change `font-family` values)
- Layout (modify padding/margins)

---

## 🎨 Design Philosophy

### Aesthetic Principles
1. **Minimalism** - Clean, uncluttered interface
2. **Romanticism** - Handwritten fonts, soft colors
3. **Clarity** - Easy to understand, intuitive flow
4. **Elegance** - Professional yet warm

### Typography Choices
- **Libre Baskerville (Italic)** - Logo, handwritten letters
- **Courier New** - Forms, buttons, UI elements
- **Uppercase + Letter-spacing** - Headers, labels

### Color Palette
- Background: `#f5f0e8` (Warm beige)
- Cards: `#fefcf7` (Off-white)
- Text: `#2d2d2d` (Dark charcoal)
- Accents: `#f5f0e8` (Muted tan)
- Borders: `#2d2d2d` (2px solid)

---

## 🔒 Security Notes

### What's Protected
✅ Answers are case-insensitive and trimmed  
✅ 5-attempt limit per name prevents brute force  
✅ Copy-paste disabled on answer fields  
✅ Each person tracked independently  

### What's Not Protected
❌ No encryption at rest (data.json is plain text)  
❌ No HTTPS enforcement (use reverse proxy)  
❌ No rate limiting per IP (only per name)  
❌ No admin panel (manual data.json edits)  

**Recommendation:** This app is designed for fun, romantic messages between trusted people. Don't use it for highly sensitive information.

---

## 🐛 Troubleshooting

### "Site can't be reached" on Railway
- Check Deploy Logs for errors
- Verify `0.0.0.0` binding in main.go
- Clear browser DNS cache: `Ctrl+Shift+Delete`

### "Vault not found" error
- Ensure `data.json` exists and is readable
- Check vault ID in URL is correct
- Restart server: `Ctrl+C` → `go run main.go`

### Leaderboard shows duplicates
- This is expected! Shows ALL attempts
- Frontend filters to best score per person

### Spinner won't stop (create page)
- Hard refresh: `Ctrl+Shift+R`
- Check browser console for errors
- Verify API endpoint is reachable

---

## 🤝 Contributing

Contributions welcome! Please:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing`)
5. Open a Pull Request

### Ideas for Contributions
- [ ] Add timer/expiration for vaults
- [ ] Email notifications when vault is unlocked
- [ ] Admin panel for managing vaults
- [ ] Multiple choice questions option
- [ ] Image/video attachments in letters
- [ ] Dark mode support
- [ ] Multiple language support
- [ ] Custom themes/colors

---

## 📄 License

MIT License - feel free to use for personal or commercial projects!

---

## 💖 Acknowledgments

- **Design Inspiration:** [digibouquet](https://digibouquet.vercel.app) by Pauline
- **Fonts:** Google Fonts (Libre Baskerville)
- **Hosting:** Railway, Render, Fly.io
- **Built with:** Go, vanilla JavaScript, pure CSS

---

## 📬 Contact & Support

- **Issues:** [GitHub Issues](https://github.com/yourusername/My.Secret/issues)
- **Discussions:** [GitHub Discussions](https://github.com/yourusername/My.Secret/discussions)
- **Email:** your.email@example.com

---

## 🎉 Made with Love

Created for lovers, dreamers, and everyone who believes in the magic of secret messages.

**"Made with love for those who understands"**

---

### Quick Links

- [Live Demo](https://mysecret-production.up.railway.app)
- [Documentation](#-how-it-works)
- [API Reference](#-api-endpoints)
- [Deployment Guide](#-deployment)

---

**Star ⭐ this repo if you like it!**