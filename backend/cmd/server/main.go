package main

import (
	"bufio"
	"log"
	"net/http"
	"os"
	"strings"

	"dashboard-builder/backend/internal/db"
	"dashboard-builder/backend/internal/httpapi"
)

// loadEnv reads KEY=VALUE lines from a .env file (if present) into the process
// environment. Real environment variables take precedence over file values.
// ponytail: minimal dotenv, no dependency; add quoting/export rules if needed.
func loadEnv(path string) {
	f, err := os.Open(path)
	if err != nil {
		return // no .env file is fine
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, val, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		val = strings.Trim(strings.TrimSpace(val), `"'`)
		if _, exists := os.LookupEnv(key); !exists {
			_ = os.Setenv(key, val)
		}
	}
}

func main() {
	loadEnv(".env")

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./data/dashboard-builder.db"
	}

	conn, err := db.Open(dbPath)
	if err != nil {
		log.Fatalf("open database: %v", err)
	}
	defer conn.Close()

	router := httpapi.NewRouter(conn)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("dashboard-builder backend listening on :%s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal(err)
	}
}
