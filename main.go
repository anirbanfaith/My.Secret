package main

import (
	"crypto/rand"
	"crypto/subtle"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	// FIXED: Pure Go driver for Railway compatibility
	"golang.org/x/crypto/argon2"
	_ "modernc.org/sqlite"
)

// FIXED: Declaring 'db' here at the package level so all functions can see it
var db *sql.DB

type Attempt struct {
	Name      string    `json:"name"`
	Score     int       `json:"score"`
	Success   bool      `json:"success"`
	CreatedAt time.Time `json:"created_at"`
}

func hashPassword(password string, salt []byte) string {
	if salt == nil {
		salt = make([]byte, 16)
		rand.Read(salt)
	}
	hash := argon2.IDKey([]byte(strings.ToLower(strings.TrimSpace(password))), salt, 1, 64*1024, 4, 32)
	// FIXED: Changed 'hex.ToString' (which doesn't exist) to 'hex.EncodeToString'
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

	// FIXED: Using '=' instead of ':=' to assign to the global 'db' variable
	db, err = sql.Open("sqlite", "./data/vaults.db")
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Create tables if they don't exist
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS vaults (
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

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS attempts (
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

// --- Handlers ---

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
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	vaultID := generateVaultID()
	answer1Hash := hashPassword(req.Answer1, nil)
	answer2Hash := hashPassword(req.Answer2, nil)
	_, err := db.Exec(`INSERT INTO vaults (id, question1, answer1_hash, question2, answer2_hash, letter) VALUES (?, ?, ?, ?, ?, ?)`,
		vaultID, req.Question1, answer1Hash, req.Question2, answer2Hash, req.Letter)
	if err != nil {
		http.Error(w, "Failed to create vault", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"vault_id": vaultID})
}

func getVaultHandler(w http.ResponseWriter, r *http.Request) {
	vaultID := strings.TrimPrefix(r.URL.Path, "/api/vault/")
	var q1, q2 string
	var isLocked bool
	err := db.QueryRow(`SELECT question1, question2, is_locked FROM vaults WHERE id = ?`, vaultID).Scan(&q1, &q2, &isLocked)
	if err != nil {
		http.Error(w, "Vault not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"vault_id": vaultID, "question1": q1, "question2": q2, "is_locked": isLocked,
	})
}

func checkAttemptsHandler(w http.ResponseWriter, r *http.Request) {
	vaultID := r.URL.Query().Get("vault_id")
	name := strings.TrimSpace(r.URL.Query().Get("name"))
	var count int
	db.QueryRow(`SELECT COUNT(*) FROM attempts WHERE vault_id = ? AND name = ? AND success = 0`, vaultID, name).Scan(&count)
	json.NewEncoder(w).Encode(map[string]interface{}{"attempts_used": count, "can_try": count < 5})
}

func unlockVaultHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		VaultID string `json:"vault_id"`
		Name    string `json:"name"`
		Answer1 string `json:"answer1"`
		Answer2 string `json:"answer2"`
	}
	json.NewDecoder(r.Body).Decode(&req)
	var a1h, a2h, letter string
	err := db.QueryRow(`SELECT answer1_hash, answer2_hash, letter FROM vaults WHERE id = ?`, req.VaultID).Scan(&a1h, &a2h, &letter)
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if verifyPassword(req.Answer1, a1h) && verifyPassword(req.Answer2, a2h) {
		db.Exec(`UPDATE vaults SET is_locked = 0 WHERE id = ?`, req.VaultID)
		json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "letter": letter})
	} else {
		db.Exec(`INSERT INTO attempts (vault_id, name, success) VALUES (?, ?, 0)`, req.VaultID, req.Name)
		http.Error(w, "Wrong answers", http.StatusUnauthorized)
	}
}

func getLeaderboardHandler(w http.ResponseWriter, r *http.Request) {
	vaultID := r.URL.Query().Get("vault_id")
	rows, err := db.Query(`SELECT name, score, success, created_at FROM attempts WHERE vault_id = ?`, vaultID)
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
	json.NewEncoder(w).Encode(attempts)
}

func main() {
	initDB()
	defer db.Close()

	http.HandleFunc("/api/create", createVaultHandler)
	http.HandleFunc("/api/vault/", getVaultHandler)
	http.HandleFunc("/api/check-attempts", checkAttemptsHandler)
	http.HandleFunc("/api/unlock", unlockVaultHandler)
	http.HandleFunc("/api/leaderboard", getLeaderboardHandler)

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "./static/create.html") })
	http.HandleFunc("/v/", func(w http.ResponseWriter, r *http.Request) { http.ServeFile(w, r, "./static/vault.html") })
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
	log.Printf("🔒 Server starting on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
