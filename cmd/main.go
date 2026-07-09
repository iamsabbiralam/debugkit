package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/iamsabbiralam/debugkit"
)

func main() {
	// Initialize the global Collector instance
	debugkit.Instance = &debugkit.Collector{
		Requests:   make([]debugkit.RequestLog, 0),
		SQLQueries: make([]debugkit.SQLQueryLog, 0),
	}

	mux := http.NewServeMux()

	// 1. Simple dummy route
	mux.HandleFunc("/api/v1/users", func(w http.ResponseWriter, r *http.Request) {
		queryStr := "SELECT id, name, email FROM users WHERE status = 'active' LIMIT 10;"

		debugkit.TrackQuery(queryStr, func() error {
			time.Sleep(120 * time.Millisecond)
			return nil
		})

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "Hello Sabbir, Users list fetched from DB!"}`))
	})

	// Dedicated route to view goroutine stack
	mux.HandleFunc("/debugkit/goroutines", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		stackTrace := debugkit.Instance.GetGoroutineStackTrace()
		w.Write([]byte(stackTrace))
	})

	// 2. Dangerous route! (Leak)
	mux.HandleFunc("/api/v1/leak", func(w http.ResponseWriter, r *http.Request) {
		ch := make(chan int)
		go func() {
			val := <-ch
			fmt.Println(val)
		}()

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "Oops! One goroutine leaked successfully."}`))
	})

	// 3. DebugKit API route
	mux.HandleFunc("/debugkit/api", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		stats := debugkit.Instance.GetFullStats()
		json.NewEncoder(w).Encode(stats)
	})

	mux.HandleFunc("/debugkit", debugkit.DebugDashboardUIHandler)

	finalHandler := debugkit.DebuggerMiddleware(mux)

	fmt.Println("🔥 DebugKit Engine with Memory & Goroutine Tracker Started!")
	fmt.Println("📌 normal route: http://localhost:8080/api/v1/users")
	fmt.Println("⚠️ leak route (goroutine leak): http://localhost:8080/api/v1/leak")
	fmt.Println("📊 debug data API: http://localhost:8080/debugkit/api")

	http.ListenAndServe(":8080", finalHandler)
}
