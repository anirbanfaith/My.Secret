# My Secret

A web application for creating and sharing encrypted secret messages with loved ones through link-based access and question verification.

## Overview

My Secret allows users to create private vaults containing secret messages that can only be unlocked by answering two custom security questions. Each vault is accessible via a unique URL and tracks all unlock attempts with a built-in leaderboard system.

## Features

- **Link-Based Access**: Each vault generates a unique shareable URL
- **Dual Question Security**: Two custom security questions protect each message
- **Attempt Tracking**: Five attempts per unique name
- **Leaderboard System**: Track and display all unlock attempts with scores
- **Multiple Unlocks**: Different people can independently unlock the same vault
- **Score Calculation**: 100 points minus 20 points per failed attempt (minimum 20)
- **Simple JSON Storage**: No database required - uses file-based storage

## Technology Stack

- **Backend**: Go (Golang)
- **Storage**: JSON file-based persistence
- **Frontend**: Vanilla JavaScript, HTML5, CSS3
- **Fonts**: Libre Baskerville (serif), Courier New (monospace)
- **Hosting**: Railway.app (or any Go-compatible platform)

## Project Structure

```
My.Secret/
├── main.go              # Go server application
├── go.mod              # Go module definition
├── go.sum              # Go dependencies checksum
├── data.json           # Storage file (auto-created)
├── .gitignore          # Git ignore rules
└── static/             # Frontend assets
    ├── index.html      # Homepage
    ├── create.html     # Vault creation page
    ├── vault.html      # Vault unlock page
    └── style.css       # Global stylesheet
```

## Local Development

### Prerequisites

- Go 1.21 or higher
- Modern web browser

### Setup

1. Clone or download the repository
2. Navigate to project directory
3. Run the application:

```bash
go run main.go
```

4. Open browser to `http://localhost:3000`

### Configuration

The application uses the `PORT` environment variable. If not set, defaults to port 3000.

## Deployment

### Railway.app

1. Push code to GitHub repository
2. Connect repository to Railway
3. Railway auto-detects Go application
4. Application deploys automatically
5. Access via generated Railway URL

### Environment Variables

- `PORT`: Server port (automatically set by hosting platform)

## Usage

### Creating a Vault

1. Navigate to the homepage
2. Click "Create Secret"
3. Enter two security questions and answers
4. Write your secret message
5. Click "Create Secret"
6. Copy the generated unique URL
7. Share URL with intended recipient

### Opening a Vault

1. Open the shared vault URL
2. Enter your name
3. Answer the two security questions
4. View the secret message upon success
5. Check the leaderboard to see all attempts

### Leaderboard

- Displays all unlock attempts
- Shows success/failure status
- Displays scores (100 to 20 points)
- Accessible before or after attempting unlock
- Groups by name (shows best attempt per person)

## Security Features

- Case-insensitive answer matching
- Answer trimming (removes extra spaces)
- Five attempts per unique name
- Failed attempts are logged
- Copy-paste prevention on answer fields

## API Endpoints

- `POST /api/create` - Create new vault
- `GET /api/vault/{id}` - Retrieve vault questions
- `GET /api/check-attempts` - Check remaining attempts
- `POST /api/unlock` - Attempt to unlock vault
- `GET /api/leaderboard` - Retrieve attempt history
- `GET /health` - Health check endpoint

## Design Philosophy

The application uses a minimalist aesthetic inspired by vintage typography and stationery design:

- Serif fonts for romantic elements
- Monospace fonts for functional elements
- Black borders and cream backgrounds
- Clean, uncluttered interface
- Professional appearance without excessive decoration

## Data Persistence

All data is stored in `data.json` in the project root:

- Vault information (questions, answers, letters)
- Attempt history (names, scores, timestamps)
- Automatic save on each operation

### Backup

To backup your data, simply copy the `data.json` file. To restore, replace the file and restart the application.

## Browser Support

- Chrome/Edge (latest)
- Firefox (latest)
- Safari (latest)
- Opera (latest)

## License

This project is provided as-is for personal use.

## Contributing

This is a personal project. Feel free to fork and modify for your own use.

## Support

For issues or questions, please refer to the deployment platform's documentation or Go language resources.
