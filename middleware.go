// middleware.go
package main

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

func DebuggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		wrapper := &responseWriterWrapper{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(wrapper, r)
		duration := time.Since(startTime)
		// Save the current goroutine count along with request data
		Instance.AddRequest(RequestLog{
			Method:           r.Method,
			Path:             r.URL.Path,
			Duration:         duration / time.Millisecond,
			Status:           wrapper.statusCode,
			GoroutinesAtTime: runtime.NumGoroutine(), // capturing the number of goroutines at the time of the request
			Timestamp:        time.Now(),
		})
	})
}
