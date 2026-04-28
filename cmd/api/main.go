package main

import (
	"fmt"
	"net/http"
	"logarda/internal/handlers"
)


func main() {
	fmt.Println("hello")
	http.HandleFunc("/health", handlers.HealthCheck)
	http.ListenAndServe(":8080", nil)
	
}


// must contain concurrency handling



// context for backend
// func main() {
//     go backgroundWorker()   // runs in background, attach go to create goroutines for the function (non-blocking) runs in background

//     http.ListenAndServe(":8080", nil) // blocks main thread for backend listening
// }

// await async yields control for other task in a single thread for other tasks

