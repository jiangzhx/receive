package date

import (
	// "fmt"
	// "strings"
	"time"
)

const (
	TimeFormat = "2006-01-02 15:04:05"
	HourFormat = "2006-01-02T15Z"
	DayFormat  = "2006-01-02"
)

func Format(when string, format string) string {
	t, _ := time.Parse(TimeFormat, when)
	return t.Format(format)
}

func PreDay(when string, days int) time.Time {
	t, _ := time.Parse(TimeFormat, when)
	return t.AddDate(0, 0, days)
}

func Parse(when string) time.Time {
	t, _ := time.Parse(TimeFormat, when)
	return t
}
