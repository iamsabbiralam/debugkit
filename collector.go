package main

import (
	"bytes"
	"runtime"
	"runtime/pprof"
	"strconv"
	"sync"
	"time"
)

// Each request's metrics
type RequestLog struct {
	Method           string        `json:"method"`
	Path             string        `json:"path"`
	Duration         time.Duration `json:"duration_ms"`
	Status           int           `json:"status"`
	GoroutinesAtTime int           `json:"goroutines_at_time"` // request time goroutine count
	Timestamp        time.Time     `json:"timestamp"`
}

// SystemSnapshot will hold the live system metrics
type SystemSnapshot struct {
	ActiveGoroutines int    `json:"active_goroutines"`
	AllocatedMemory  string `json:"allocated_memory"`   // MB size
	TotalSystemAlloc string `json:"total_system_alloc"` // Total memory allocated by the OS
	HeapObjects      uint64 `json:"heap_objects"`       // Number of objects in the heap
}

// combined response structure (for dashboard)
type DebugDataResponse struct {
	System   SystemSnapshot `json:"system"`
	Requests []RequestLog   `json:"requests"`
}

type Collector struct {
	mu       sync.Mutex
	Requests []RequestLog
}

var Instance = &Collector{
	Requests: make([]RequestLog, 0),
}

func (c *Collector) AddRequest(log RequestLog) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.Requests) > 100 {
		c.Requests = c.Requests[1:]
	}
	c.Requests = append(c.Requests, log)
}

// this is the main function that will be called to get the full debug data
func (c *Collector) GetFullStats() DebugDataResponse {
	c.mu.Lock()
	reqs := make([]RequestLog, len(c.Requests))
	copy(reqs, c.Requests) // for thread safety, we copy the slice to
	c.mu.Unlock()

	// live memory data reading
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return DebugDataResponse{
		System: SystemSnapshot{
			ActiveGoroutines: runtime.NumGoroutine(),
			AllocatedMemory:  byteToMB(m.Alloc),
			TotalSystemAlloc: byteToMB(m.Sys),
			HeapObjects:      m.HeapObjects,
		},
		Requests: reqs,
	}
}

// Helper function: convert bytes to readable MB
func byteToMB(b uint64) string {
	mb := b / 1024 / 1024
	return strconv.FormatUint(mb, 10) + " MB"
}

// Method to get detailed information about running goroutines
func (c *Collector) GetGoroutineStackTrace() string {
	var buf bytes.Buffer
	// debug=1 will give a concise stack trace of all goroutines
	pprof.Lookup("goroutine").WriteTo(&buf, 1)
	return buf.String()
}
