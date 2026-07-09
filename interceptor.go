// interceptor.go
package debugkit

import (
	"time"
)

// TrackQuery takes a SQL query string and a function that executes the query, measures the execution time, and logs it to the global Collector instance.
func TrackQuery(query string, executeFn func() error) error {
	startTime := time.Now()
	// Execute the provided function that runs the SQL query
	err := executeFn()
	duration := time.Since(startTime)
	// Add the SQL query log to the global Collector instance
	Instance.AddSQLQuery(SQLQueryLog{
		Query:     query,
		Duration:  duration / time.Millisecond,
		Timestamp: time.Now(),
	})

	return err
}
