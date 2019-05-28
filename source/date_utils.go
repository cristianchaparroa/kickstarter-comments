package main

import "time"

func unixMilli(t time.Time) int64 {
	return t.Round(time.Millisecond).UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}

// MakeTimestampMilli retrieves the time since 1970 to now in milliseconds
func MakeTimestampMilli() int64 {
	return unixMilli(time.Now())
}
