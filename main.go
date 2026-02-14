package main

import (
	"crypto/rand"
	"crypto/subtle"
	"database/sql"
	"encoding/hex"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/crypto/argon2"
	_ "modernc.org/sqlite" // Pure Go driver for Railway
)

// GLOBAL VARIABLE DECLARATION
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
	// FIXED: Changed hex.ToString to hex.EncodeToString
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

func initDB() {
	var err error
	if err := os.MkdirAll("./data", 0755); err != nil {
		log.Fatal("Failed to create data directory:", err)
	}

	// FIXED: Using '=' instead of ':=' to use the global 'db' variable
	db, err = sql.Open("sqlite", "./data/vaults.db")
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	log.Println("✅ Database initialized successfully")
}

// ... include your handlers (createVaultHandler, etc.) here ...

func main() {
	initDB()
	defer db.Close()

	// Static files and routes
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

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
