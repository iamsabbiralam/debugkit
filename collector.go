package debugkit

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

type SQLQueryLog struct {
	Query     string        `json:"query"`
	Duration  time.Duration `json:"duration_ms"`
	Timestamp time.Time     `json:"timestamp"`
}

// combined response structure (for dashboard)
type DebugDataResponse struct {
	System   SystemSnapshot `json:"system"`
	Requests []RequestLog   `json:"requests"`
	SQLs     []SQLQueryLog  `json:"sql_queries"`
}

type Collector struct {
	mu         sync.Mutex
	Requests   []RequestLog
	SQLQueries []SQLQueryLog
}

// global Instance variable to hold the Collector instance, shared across the debugkit package
var Instance *Collector

func (c *Collector) AddRequest(log RequestLog) {
	if c == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.Requests) > 100 {
		c.Requests = c.Requests[1:]
	}
	c.Requests = append(c.Requests, log)
}

func (c *Collector) AddSQLQuery(log SQLQueryLog) {
	if c == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.SQLQueries) > 100 {
		c.SQLQueries = c.SQLQueries[1:]
	}
	c.SQLQueries = append(c.SQLQueries, log)
}

func (c *Collector) GetFullStats() DebugDataResponse {
	if c == nil {
		return DebugDataResponse{}
	}
	c.mu.Lock()
	reqs := make([]RequestLog, len(c.Requests))
	copy(reqs, c.Requests)
	sqls := make([]SQLQueryLog, len(c.SQLQueries))
	copy(sqls, c.SQLQueries)
	c.mu.Unlock()

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
		SQLs:     sqls,
	}
}

func (c *Collector) GetGoroutineStackTrace() string {
	if c == nil {
		return "Collector is not initialized"
	}
	var buf bytes.Buffer
	pprof.Lookup("goroutine").WriteTo(&buf, 1)
	return buf.String()
}

func byteToMB(b uint64) string {
	mb := b / 1024 / 1024
	return strconv.FormatUint(mb, 10) + " MB"
}
