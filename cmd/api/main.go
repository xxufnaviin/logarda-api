package main

import (
	"logarda/internal/db"
	"logarda/internal/handlers"
	"net/http"
)


func main() {
	handlers.Init()
	defer db.DB.Close() // calls before function closes
	
	http.HandleFunc("/health", handlers.HealthCheck)
	http.HandleFunc("/login", handlers.Login)

	http.ListenAndServe(":8080", nil)
	
}



// must contain concurrency handling



// context for backend
// func main() {
//     go backgroundWorker()   // runs in background, attach go to create goroutines for the function (non-blocking) runs in background

//     http.ListenAndServe(":8080", nil) // blocks main thread for backend listening
// }

// await async yields control for other task in a single thread for other tasks

