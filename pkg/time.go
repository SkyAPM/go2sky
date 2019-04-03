package pkg

import "time"

func Millisecond(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}
