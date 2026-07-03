// main.go
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	// 1. Simple dummy route
	mux.HandleFunc("/api/v1/users", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "Users loaded successfully"}`))
	})

	// Dedicated route to view goroutine stack
	mux.HandleFunc("/debugkit/goroutines", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain") // This will display as readable text in the browser
		w.Header().Set("Access-Control-Allow-Origin", "*")

		stackTrace := Instance.GetGoroutineStackTrace()
		w.Write([]byte(stackTrace))
	})

	// 2. Dangerous route! Hitting this will cause a goroutine to be stuck forever (Leak)
	mux.HandleFunc("/api/v1/leak", func(w http.ResponseWriter, r *http.Request) {
		// This channel will never be closed or receive data
		ch := make(chan int)

		// Create a goroutine leak
		go func() {
			val := <-ch // This will cause the goroutine to hang forever
			fmt.Println(val)
		}()

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "Oops! One goroutine leaked successfully."}`))
	})

	// 3. Updated DebugKit API route
	mux.HandleFunc("/debugkit/api", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// Now we will send the complete combined stats
		stats := Instance.GetFullStats()
		json.NewEncoder(w).Encode(stats)
	})

	mux.HandleFunc("/debugkit", DebugDashboardUIHandler)

	finalHandler := DebuggerMiddleware(mux)

	fmt.Println("🔥 DebugKit Engine with Memory & Goroutine Tracker Started!")
	fmt.Println("📌 normal route: http://localhost:8080/api/v1/users")
	fmt.Println("⚠️ leak route (goroutine leak): http://localhost:8080/api/v1/leak")
	fmt.Println("📊 debug data API: http://localhost:8080/debugkit/api")

	http.ListenAndServe(":8080", finalHandler)
}
