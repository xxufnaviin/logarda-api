package handlers

import (
	"encoding/json"
	"logarda/internal/config"
	"logarda/internal/db"
	"net/http"
)

func Init() {
	// load secrets
	config.LoadSecrets()
	// create postgresql connection pool
	db.CreatePostgreSQLPool()
	// create redis connection
	db.CreateRedisClient()
}

func GetHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(map[string]any{
		"message": "ok",
		"status":  http.StatusOK,
	})

}
