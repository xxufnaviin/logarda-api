package handlers

import (
	"encoding/json"
	"fmt"
	"logarda/internal/config"
	"logarda/internal/model"
	"net/http"
)

var performanceEndpoint = "models/%s/performance"

func GetModelPerformance(w http.ResponseWriter, r *http.Request) {
	var request model.ModelPerformanceRequest
	var response any

	// get model from query params
	request.Model = r.URL.Query().Get("model")

	if request.Model == "" {
		http.Error(w, "Invalid Parameters", http.StatusBadRequest) // return error if not able to get model
		return
	}

	requestEndpoint := fmt.Sprintf(performanceEndpoint, request.Model)
	request_url := fmt.Sprintf("%s%s", config.ANALYTICS_API, requestEndpoint)

	w.Header().Set("Content-Type", "application/json")

	resp, err := http.Get(request_url)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]any{
			"message": "fetch failed",
			"status":  500,
			"error":   "internal server error",
		})
		return
	}
	defer resp.Body.Close()

	json.NewDecoder(resp.Body).Decode(&response)
	json.NewEncoder(w).Encode(response)
}
