# DebugKit 🚀

[![Go Report Card](https://goreportcard.com/badge/github.com/iamsabbiralam/debugkit)](https://goreportcard.com/report/github.com/iamsabbiralam/debugkit)
[![Go Reference](https://pkg.go.dev/badge/github.com/iamsabbiralam/debugkit.svg)](https://pkg.go.dev/github.com/iamsabbiralam/debugkit)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

DebugKit is a lightweight, zero-dependency, real-time telemetry and goroutine leak tracking engine for Go (Golang) applications. It automatically monitors active goroutines, memory allocations, and SQL query performance, exposing everything through a beautiful embedded dashboard running directly inside your application.

---

## 🌟 Key Features

- **📊 Live Memory Telemetry**: Track allocated heap memory, total system allocation (OS-level memory), and live heap objects in real-time.
- **⚠️ Goroutine Leak Detection**: Instantly detect hanging background goroutines and inspect full goroutine stack traces directly in the browser.
- **🔍 SQL Query Interceptor & Benchmarking**: Measure database query execution latencies and logs historically via an isolated query wrapper.
- **🎨 Embedded Premium Dashboard**: Powered by Go's native `embed.FS` feature—no external asset CDNs, extra configuration, or external file requirements needed.
- **🛡️ Production-Safe (Nil Guards)**: Built with defensive software architecture. If initialization is accidentally omitted, internal guards prevent `nil pointer dereference` panics, ensuring your core web server never crashes.

---

## 📦 Installation

To add DebugKit to your Go project, run the following command in your terminal:

```bash
go get [github.com/iamsabbiralam/debugkit@v1.0.0](https://github.com/iamsabbiralam/debugkit@v1.0.0)
```

## 🛠️ Quick Start & Integration
Integrating DebugKit into any standard `net/http` application is extremely simple. Here is a complete, production-ready example:

```package main

import (
	"fmt"
	"net/http"
	"time"

	"[github.com/iamsabbiralam/debugkit](https://github.com/iamsabbiralam/debugkit)"
)

func main() {
	// 1. Initialize the global DebugKit engine
	debugkit.New()

	mux := http.NewServeMux()

	// A dummy endpoint tracking database latency
	mux.HandleFunc("/api/v1/users", func(w http.ResponseWriter, r *http.Request) {
		queryStr := "SELECT id, name, email FROM users WHERE status = 'active' LIMIT 10;"

		// Use TrackQuery interceptor to log and benchmark query latency
		debugkit.TrackQuery(queryStr, func() error {
			time.Sleep(120 * time.Millisecond) // Simulating database latency
			return nil
		})

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"message": "Hello World, Users list fetched successfully!"}`))
	})

	// Simulating a dangerous goroutine leak endpoint (Hitting this leaks 1 goroutine forever)
	mux.HandleFunc("/api/v1/leak", func(w http.ResponseWriter, r *http.Request) {
		ch := make(chan int)
		go func() {
			val := <-ch // Hangs forever because the channel never receives data
			fmt.Println(val)
		}()
		w.Write([]byte(`{"message": "One goroutine leaked successfully."}`))
	})

	// 2. Mount the embedded premium dashboard and telemetry endpoints
	// This opens up the /debugkit endpoint on your active router
	debugkit.RegisterUIHandlers(mux)

	// 3. Wrap your Multiplexer with the DebugKit Middleware to log active HTTP traffic
	finalHandler := debugkit.DebuggerMiddleware(mux)

	fmt.Println("🔥 DebugKit Engine with Memory & Goroutine Tracker Started!")
	fmt.Println("📊 Open Dashboard UI: http://localhost:8080/debugkit")
	fmt.Println("⚙️ Raw Telemetry API: http://localhost:8080/debugkit/api")

	http.ListenAndServe(":8080", finalHandler)
}
```
## 🗺️ Exposed HTTP Endpoints
Once `RegisterUIHandlers` is invoked on your router, the following paths are activated automatically:

| Endpoint | Content Type | Description |
|---|---|---|
| `/debugkit` | `text/html` | The main visual telemetry dashboard UI. |
| `/debugkit/api` | `application/json` | Raw system stats, request logs, and SQL histories in a JSON payload. |

## 🔒 Safety & Defensive Architecture
This package was designed with strict production stability in mind. In Go applications, a `nil pointer dereference` is a fatal error that can drop a live server. DebugKit mitigates this entirely by placing internal Nil Checks at the receiver level of all public methods:
```
func (c *Collector) AddRequest(log RequestLog) {
      if c == nil { 
            return // Completely eliminates runtime memory panics/crashes
      }
      // ...
}
```
If a developer forgets to invoke `debugkit.New()`, the package drops tracking operations silently without disrupting your core HTTP servers or raising runtime panics.

## 📄 License
This project is licensed under the MIT License - see the LICENSE file for details.