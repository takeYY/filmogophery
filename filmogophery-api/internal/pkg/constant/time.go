package constant

import "time"

const (
	DateFormat  string = "2006-01-02"
	DefaultDate string = "1895-12-28"
)

type (
	Date string
)

func GetDefaultDate() time.Time {
	result, _ := time.ParseInLocation(DateFormat, DefaultDate, time.Local)
	return result
}
