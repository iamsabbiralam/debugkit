package debugkit

import (
	"net/http"
	"runtime"
	"time"
)

type responseWriterWrapper struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriterWrapper) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// DebuggerMiddleware intercepts HTTP requests to capture lifecycle metrics
func DebuggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Ignore internal DebugKit requests to avoid cluttering the logs
		if r.URL.Path == "/debugkit" || r.URL.Path == "/debugkit/api" {
			next.ServeHTTP(w, r)
			return
		}

		startTime := time.Now()
		wrapper := &responseWriterWrapper{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(wrapper, r)
		duration := time.Since(startTime)
		// Only push logs if the global Instance has been initialized with New()
		if Instance != nil {
			Instance.AddRequest(RequestLog{
				Method:           r.Method,
				Path:             r.URL.Path,
				Duration:         duration / time.Millisecond,
				Status:           wrapper.statusCode,
				GoroutinesAtTime: runtime.NumGoroutine(), // Capturing thread pool allocation
				Timestamp:        time.Now(),
			})
		}
	})
}
