package pkg

import "time"

// Millisecond converts time to unix millisecond
func Millisecond(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}
