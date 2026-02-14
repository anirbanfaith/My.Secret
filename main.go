package main

import (
	"crypto/rand"
	"crypto/subtle"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	// UPDATED: Using pure Go SQLite driver to fix Railway CGO errors
	"golang.org/x/crypto/argon2"
	_ "modernc.org/sqlite"
)

type Attempt struct {
	Name      string    `json:"name"`
	Score     int       `json:"score"`
	Success   bool      `json:"success"`
	CreatedAt time.Time `json:"created_at"`
}

var db *sql.DB

func hashPassword(password string, salt []byte) string {
	if salt == nil {
		salt = make([]byte, 16)
		rand.Read(salt)
	}
	hash := argon2.IDKey([]byte(strings.ToLower(strings.TrimSpace(password))), salt, 1, 64*1024, 4, 32)
	return hex.EncodeToString(salt) + ":" + hex.EncodeToString(hash)
}

func verifyPassword(password, hashedPassword string) bool {
	parts := strings.Split(hashedPassword, ":")
	if len(parts) != 2 {
		return false
	}
	salt, _ := hex.DecodeString(parts[0])
	storedHash, _ := hex.DecodeString(parts[1])
	newHash := argon2.IDKey([]byte(strings.ToLower(strings.TrimSpace(password))), salt, 1, 64*1024, 4, 32)
	return subtle.ConstantTimeCompare(newHash, storedHash) == 1
}

func generateVaultID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func initDB() {
	var err error

	// Create data directory
	if err := os.MkdirAll("./data", 0755); err != nil {
		log.Fatal("Failed to create data directory:", err)
	}

	// UPDATED: Driver name changed from "sqlite3" to "sqlite"
	db, err = sql.Open("sqlite", "./data/vaults.db")
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Create vaults table
	_, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS vaults (
        id TEXT PRIMARY KEY,
        question1 TEXT NOT NULL,
        answer1_hash TEXT NOT NULL,
        question2 TEXT NOT NULL,
        answer2_hash TEXT NOT NULL,
        letter TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        is_locked BOOLEAN DEFAULT 1
    )`)
	if err != nil {
		log.Fatal("Failed to create vaults table:", err)
	}

	// Create attempts table
	_, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS attempts (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        vault_id TEXT NOT NULL,
        name TEXT NOT NULL,
        score INTEGER DEFAULT 100,
        success BOOLEAN DEFAULT 0,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    )`)
	if err != nil {
		log.Fatal("Failed to create attempts table:", err)
	}

	log.Println("✅ Database initialized successfully")
}

func createVaultHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Question1 string `json:"question1"`
		Answer1   string `json:"answer1"`
		Question2 string `json:"question2"`
		Answer2   string `json:"answer2"`
		Letter    string `json:"letter"`
	}

	json.NewDecoder(r.Body).Decode(&req)

	if req.Question1 == "" || req.Answer1 == "" || req.Question2 == "" || req.Answer2 == "" || req.Letter == "" {
		log.Println("❌ Create vault failed: missing fields")
		http.Error(w, "All fields required", http.StatusBadRequest)
		return
	}

	vaultID := generateVaultID()
	answer1Hash := hashPassword(req.Answer1, nil)
	answer2Hash := hashPassword(req.Answer2, nil)

	_, err := db.Exec(`INSERT INTO vaults (id, question1, answer1_hash, question2, answer2_hash, letter)
        VALUES (?, ?, ?, ?, ?, ?)`, vaultID, req.Question1, answer1Hash, req.Question2, answer2Hash, req.Letter)

	if err != nil {
		log.Printf("❌ Database insert failed: %v\n", err)
		http.Error(w, "Failed to create vault", http.StatusInternalServerError)
		return
	}

	scheme := "http"
	if r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}
	vaultURL := fmt.Sprintf("%s://%s/v/%s", scheme, r.Host, vaultID)

	log.Printf("✅ Vault created: %s\n", vaultID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"vault_id":  vaultID,
		"vault_url": vaultURL,
	})
}

func getVaultHandler(w http.ResponseWriter, r *http.Request) {
	vaultID := strings.TrimPrefix(r.URL.Path, "/api/vault/")
	log.Printf("📖 Getting vault: %s\n", vaultID)

	var question1, question2 string
	var isLocked bool

	err := db.QueryRow(`SELECT question1, question2, is_locked FROM vaults WHERE id = ?`, vaultID).
		Scan(&question1, &question2, &isLocked)

	if err == sql.ErrNoRows {
		log.Printf("❌ Vault not found: %s\n", vaultID)
		http.Error(w, "Vault not found", http.StatusNotFound)
		return
	}

	if err != nil {
		log.Printf("❌ Database error: %v\n", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	log.Printf("✅ Vault found: %s (locked: %v)\n", vaultID, isLocked)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"vault_id":  vaultID,
		"question1": question1,
		"question2": question2,
		"is_locked": isLocked,
	})
}

func checkAttemptsHandler(w http.ResponseWriter, r *http.Request) {
	vaultID := r.URL.Query().Get("vault_id")
	name := strings.TrimSpace(r.URL.Query().Get("name"))

	var count int
	db.QueryRow(`SELECT COUNT(*) FROM attempts WHERE vault_id = ? AND name = ? AND success = 0`,
		vaultID, name).Scan(&count)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"attempts_used": count,
		"attempts_left": 5 - count,
		"can_try":       count < 5,
	})
}

func unlockVaultHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		VaultID string `json:"vault_id"`
		Name    string `json:"name"`
		Answer1 string `json:"answer1"`
		Answer2 string `json:"answer2"`
	}

	json.NewDecoder(r.Body).Decode(&req)
	req.Name = strings.TrimSpace(req.Name)

	var attemptCount int
	db.QueryRow(`SELECT COUNT(*) FROM attempts WHERE vault_id = ? AND name = ? AND success = 0`,
		req.VaultID, req.Name).Scan(&attemptCount)

	if attemptCount >= 5 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success":       false,
			"max_reached":   true,
			"attempts_left": 0,
		})
		return
	}

	var answer1Hash, answer2Hash, letter string
	err := db.QueryRow(`SELECT answer1_hash, answer2_hash, letter FROM vaults WHERE id = ?`,
		req.VaultID).Scan(&answer1Hash, &answer2Hash, &letter)

	if err != nil {
		http.Error(w, "Vault not found", http.StatusNotFound)
		return
	}

	answer1Correct := verifyPassword(req.Answer1, answer1Hash)
	answer2Correct := verifyPassword(req.Answer2, answer2Hash)

	if answer1Correct && answer2Correct {
		score := 100 - (attemptCount * 20)
		if score < 20 {
			score = 20
		}

		db.Exec(`INSERT INTO attempts (vault_id, name, score, success) VALUES (?, ?, ?, 1)`,
			req.VaultID, req.Name, score)
		db.Exec(`UPDATE vaults SET is_locked = 0 WHERE id = ?`, req.VaultID)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"letter":  letter,
			"score":   score,
		})
		return
	}

	db.Exec(`INSERT INTO attempts (vault_id, name, score, success) VALUES (?, ?, 0, 0)`,
		req.VaultID, req.Name)

	attemptCount++
	attemptsLeft := 5 - attemptCount

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":       false,
		"attempts_left": attemptsLeft,
		"max_reached":   attemptsLeft == 0,
	})
}

func getLeaderboardHandler(w http.ResponseWriter, r *http.Request) {
	vaultID := r.URL.Query().Get("vault_id")

	rows, err := db.Query(`
        SELECT name, score, success, created_at 
        FROM attempts 
        WHERE vault_id = ? 
        ORDER BY success DESC, score DESC, created_at ASC`, vaultID)

	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var attempts []Attempt
	for rows.Next() {
		var a Attempt
		rows.Scan(&a.Name, &a.Score, &a.Success, &a.CreatedAt)
		attempts = append(attempts, a)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(attempts)
}

func main() {
	initDB()
	defer db.Close()

	// API endpoints
	http.HandleFunc("/api/create", createVaultHandler)
	http.HandleFunc("/api/vault/", getVaultHandler)
	http.HandleFunc("/api/check-attempts", checkAttemptsHandler)
	http.HandleFunc("/api/unlock", unlockVaultHandler)
	http.HandleFunc("/api/leaderboard", getLeaderboardHandler)

	// Static files
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Page routes
	http.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/create.html")
	})

	http.HandleFunc("/v/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/vault.html")
	})

	// Homepage
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFile(w, r, "./static/index.html")
			return
		}
		http.NotFound(w, r)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("🔒 Secret Vault starting on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
