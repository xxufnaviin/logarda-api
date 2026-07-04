package main

import (
	"fmt"
	"logarda/internal/db"
	"logarda/internal/handlers"
	"logarda/internal/workers"
	"net/http"
)

func main() {
	handlers.Init()
	go workers.ErrorLogsWorker()
	go workers.MetricStreamWorker()
	go workers.PredictedMetricStreamWorker()

	defer db.DB.Close() // calls before function closes

	http.HandleFunc("/api/health", handlers.GetHealth)
	http.HandleFunc("/api/auth/register", handlers.Register)
	http.HandleFunc("/api/auth/login", handlers.Login)
	http.HandleFunc("/api/aws/credentials/update", handlers.SaveAWSCredentials)
	http.HandleFunc("/api/metrics", handlers.GetMetrics)

	http.HandleFunc("/api/logs", handlers.GetErrorLogs)
	http.HandleFunc("/api/logs/stats", handlers.GetErrorLogStats)

	http.HandleFunc("/api/user", handlers.GetUserDetails)
	http.HandleFunc("/api/user/collector", handlers.SetCollectorOn)

	http.HandleFunc("/websocket", handlers.WebsocketHandler)

	http.HandleFunc("/api/analytics/predict", handlers.PredictMetrics)
	http.HandleFunc("/api/models/performance", handlers.GetModelPerformance)

	fmt.Println("Logarda backend server started at http://localhost:8080")
	http.ListenAndServe(":8080", nil)

}

// must contain concurrency handling

// context for backend
// func main() {
//     go backgroundWorker()   // runs in background, attach go to create goroutines for the function (non-blocking) runs in background

//     http.ListenAndServe(":8080", nil) // blocks main thread for backend listening
// }

// await async yields control for other task in a single thread for other tasks
