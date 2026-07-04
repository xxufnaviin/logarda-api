package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"logarda/internal/config"
	"logarda/internal/db"
	"logarda/internal/model"
	"net/http"
	"net/url"
)

var predictionEndpoint = "analytics/predict"

func GetMetrics(w http.ResponseWriter, r *http.Request) {
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
		request.Duration = "3" // if empty set duration to default (3 hours)
	}

	w.Header().Set("Content-Type", "application/json")

	// get all metrics here
	metrics, err := db.GetMetrics(ctx, request.Username, request.Duration)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]any{
			"message": "failed to get metrics",
			"status":  400,
			"error":   err.Error(),
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]any{
		"data":    metrics,
		"message": "success"})
	log.Printf("Metrics data for the past %s hours fetched for %s\n", request.Duration, request.Username)
}

func PredictMetrics(w http.ResponseWriter, r *http.Request) {
	var request model.MetricsLogsRequest
	var response model.PredictMetricsResponse

	// get username and duration from query params
	request.Username = r.URL.Query().Get("username")
	request.Duration = r.URL.Query().Get("duration")

	if request.Username == "" {
		http.Error(w, "Invalid Parameters", http.StatusBadRequest) // return error if not able to get username
		return
	}
	if request.Duration == "" {
		http.Error(w, "Invalid Parameters", http.StatusBadRequest) // return error if not able to get duration
		return
	}

	params := url.Values{}
	params.Add("duration", request.Duration)
	params.Add("username", request.Username)

	// construct request url using endpoint name and request params
	request_url := fmt.Sprintf("%s%s?%s", config.ANALYTICS_API, predictionEndpoint, params.Encode())

	w.Header().Set("Content-Type", "application/json")

	resp, err := http.Get(request_url)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]any{
			"message": "prediction failed",
			"status":  500,
			"error":   "internal server error",
		})
		return
	}
	defer resp.Body.Close()

	// get response from analytical endpoint and return it to user
	json.NewDecoder(resp.Body).Decode(&response)
	json.NewEncoder(w).Encode(response)

	log.Printf("Metrics predicted for the next %s hours for %s\n", request.Duration, request.Username)
}
