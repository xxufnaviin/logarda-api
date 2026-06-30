package handlers

import (
	"encoding/json"
	"log"
	"logarda/internal/config"
	"logarda/internal/db"
	"net/http"
	"time"
)

func Init() {
	// load secrets
	config.LoadSecrets()
	// create postgresql connection pool
	db.CreatePostgreSQLPool()
	// create redis connection
	db.CreateRedisClient()

	time.Sleep(5 * time.Second) // init cooldown before starting workers
}

func GetHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]any{
		"message": "ok",
		"status":  http.StatusOK,
	})
	log.Println("Health Check Status: Online")
}
