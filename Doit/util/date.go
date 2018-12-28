package util

import "time"

const (
	DateTimeFmt = "2006-01-02 15:04:05"
	DateFmt     = "2006-01-02"
)

func DateTimeStd() string {
	return time.Now().Format(DateTimeFmt)
}

func DateStd() string {
	return time.Now().Format(DateFmt)
}
