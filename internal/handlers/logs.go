package handlers

import (
	"encoding/json"
	"logarda/internal/db"
	"logarda/internal/model"
	"net/http"
)

func GetErrorLogs(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var request model.MetricsLogsRequest

	// get username and duration from query params
	request.Username = r.URL.Query().Get("username")
	request.Duration = r.URL.Query().Get("duration")

	if request.Username == "" {
		http.Error(w, "Invalid Parameters", http.StatusBadRequest) // return error if not able to get username
		return
	}
	if request.Duration == "" {
		request.Duration = "24" // if empty set duration to default (24 hours)
	}

	w.Header().Set("Content-Type", "application/json")

	// get all error logs here
	logs, err := db.GetErrorLogs(ctx, request.Username, request.Duration)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]any{
			"message": "failed to get error logs",
			"status":  400,
			"error":   err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]any{
		"data":    logs,
		"message": "success"})
}
