// interceptor.go
package main

import (
	"time"
)

// এটি একটি কাস্টম হেল্পার যা ডেভেলপাররা ডাটাবেজ কুয়েরি এক্সিকিউট করার সময় রিয়েল-টাইমে টাইম ট্র্যাক করতে ব্যবহার করতে পারবে।
// (ভবিষ্যতে আমরা এটাকে বিল্ট-ইন database/sql/driver দিয়ে অটোমেট করব, এখন PoC এর জন্য কাস্টম ট্র্যাকার বানাচ্ছি)
func TrackQuery(query string, executeFn func() error) error {
	startTime := time.Now()

	// কুয়েরি রান করা
	err := executeFn()

	duration := time.Since(startTime)

	// কালেক্টরে কুয়েরি লগা করা
	Instance.AddSQLQuery(SQLQueryLog{
		Query:     query,
		Duration:  duration / time.Millisecond,
		Timestamp: time.Now(),
	})

	return err
}
