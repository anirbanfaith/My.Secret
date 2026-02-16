package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type Vault struct {
	ID        string    `json:"id"`
	Question1 string    `json:"question1"`
	Answer1   string    `json:"answer1"`
	Question2 string    `json:"question2"`
	Answer2   string    `json:"answer2"`
	Letter    string    `json:"letter"`
	CreatedAt time.Time `json:"created_at"`
}

type Attempt struct {
	Name      string    `json:"name"`
	Score     int       `json:"score"`
	Success   bool      `json:"success"`
	CreatedAt time.Time `json:"created_at"`
}

type Storage struct {
	Vaults   map[string]Vault     `json:"vaults"`
	Attempts map[string][]Attempt `json:"attempts"`
	mu       sync.RWMutex
}

var storage = &Storage{
	Vaults:   make(map[string]Vault),
	Attempts: make(map[string][]Attempt),
}

func getStoragePath() string {
	if _, err := os.Stat("/data"); err == nil {
		return "/data/data.json"
	}
	return "data.json"
}

var storageFile = getStoragePath()

func loadStorage() {
	data, err := os.ReadFile(storageFile)
	if err != nil {
		log.Println("No existing data, starting fresh")
		return
	}

	storage.mu.Lock()
	defer storage.mu.Unlock()

	if err := json.Unmarshal(data, storage); err != nil {
		log.Printf("Error parsing storage: %v", err)
		return
	}
	log.Printf("✅ Loaded %d vaults from storage\n", len(storage.Vaults))
}

func saveStorage() error {
	storage.mu.RLock()
	defer storage.mu.RUnlock()

	data, err := json.MarshalIndent(storage, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(storageFile, data, 0644)
}

func generateID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func main() {
	loadStorage()

	http.HandleFunc("/api/create", createVault)
	http.HandleFunc("/api/vault/", getVault)
	http.HandleFunc("/api/check-attempts", checkAttempts)
	http.HandleFunc("/api/unlock", unlockVault)
	http.HandleFunc("/api/leaderboard", getLeaderboard)

	// Backup endpoint - uses BACKUP_KEY environment variable
	http.HandleFunc("/api/backup", func(w http.ResponseWriter, r *http.Request) {
		backupKey := os.Getenv("BACKUP_KEY")
		if backupKey == "" {
			http.Error(w, "Backup disabled", http.StatusForbidden)
			return
		}

		if r.URL.Query().Get("key") != backupKey {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		data, err := os.ReadFile(storageFile)
		if err != nil {
			http.Error(w, "Failed to read data", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Disposition", "attachment; filename=backup.json")
		w.Write(data)
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	http.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/robots.txt")
	})

	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/favicon.ico")
	})

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/create.html")
	})

	http.HandleFunc("/v/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/vault.html")
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, "./static/index.html")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("✨ My Secret starting on port %s\n", port)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, nil))
}

func createVault(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req Vault
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if req.Question1 == "" || req.Answer1 == "" || req.Question2 == "" || req.Answer2 == "" || req.Letter == "" {
		http.Error(w, "All fields required", http.StatusBadRequest)
		return
	}

	vault := Vault{
		ID:        generateID(),
		Question1: req.Question1,
		Answer1:   strings.ToLower(strings.TrimSpace(req.Answer1)),
		Question2: req.Question2,
		Answer2:   strings.ToLower(strings.TrimSpace(req.Answer2)),
		Letter:    req.Letter,
		CreatedAt: time.Now(),
	}

	storage.mu.Lock()
	storage.Vaults[vault.ID] = vault
	storage.mu.Unlock()

	saveStorage()

	scheme := "http"
	if r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}
	vaultURL := fmt.Sprintf("%s://%s/v/%s", scheme, r.Host, vault.ID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"vault_id":  vault.ID,
		"vault_url": vaultURL,
	})
}

func getVault(w http.ResponseWriter, r *http.Request) {
	vaultID := strings.TrimPrefix(r.URL.Path, "/api/vault/")

	storage.mu.RLock()
	vault, exists := storage.Vaults[vaultID]
	storage.mu.RUnlock()

	if !exists {
		http.Error(w, "Vault not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"vault_id":  vault.ID,
		"question1": vault.Question1,
		"question2": vault.Question2,
	})
}

func checkAttempts(w http.ResponseWriter, r *http.Request) {
	vaultID := r.URL.Query().Get("vault_id")
	name := strings.TrimSpace(r.URL.Query().Get("name"))

	storage.mu.RLock()
	attempts := storage.Attempts[vaultID]
	storage.mu.RUnlock()

	count := 0
	for _, a := range attempts {
		if a.Name == name && !a.Success {
			count++
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"attempts_used": count,
		"attempts_left": 5 - count,
		"can_try":       count < 5,
	})
}

func unlockVault(w http.ResponseWriter, r *http.Request) {
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

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Answer1 = strings.ToLower(strings.TrimSpace(req.Answer1))
	req.Answer2 = strings.ToLower(strings.TrimSpace(req.Answer2))

	storage.mu.RLock()
	vault, exists := storage.Vaults[req.VaultID]
	attempts := storage.Attempts[req.VaultID]
	storage.mu.RUnlock()

	if !exists {
		http.Error(w, "Vault not found", http.StatusNotFound)
		return
	}

	attemptCount := 0
	for _, a := range attempts {
		if a.Name == req.Name && !a.Success {
			attemptCount++
		}
	}

	if attemptCount >= 5 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success":     false,
			"max_reached": true,
		})
		return
	}

	if req.Answer1 == vault.Answer1 && req.Answer2 == vault.Answer2 {
		score := 100 - (attemptCount * 20)
		if score < 20 {
			score = 20
		}

		attempt := Attempt{
			Name:      req.Name,
			Score:     score,
			Success:   true,
			CreatedAt: time.Now(),
		}

		storage.mu.Lock()
		storage.Attempts[req.VaultID] = append(storage.Attempts[req.VaultID], attempt)
		storage.mu.Unlock()
		saveStorage()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"letter":  vault.Letter,
			"score":   score,
		})
		return
	}

	attempt := Attempt{
		Name:      req.Name,
		Score:     0,
		Success:   false,
		CreatedAt: time.Now(),
	}

	storage.mu.Lock()
	storage.Attempts[req.VaultID] = append(storage.Attempts[req.VaultID], attempt)
	storage.mu.Unlock()
	saveStorage()

	attemptCount++

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":       false,
		"attempts_left": 5 - attemptCount,
		"max_reached":   attemptCount >= 5,
	})
}

func getLeaderboard(w http.ResponseWriter, r *http.Request) {
	vaultID := r.URL.Query().Get("vault_id")

	storage.mu.RLock()
	attempts := storage.Attempts[vaultID]
	storage.mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(attempts)
}
