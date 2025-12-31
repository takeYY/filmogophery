package constant

import "time"

const (
	DateFormat  string = "2006-01-02"
	DefaultDate string = "1895-12-28"
)

type (
	Date     string
	Datetime string // UTC datetime like '2006-01-02T15:04:05Z'
)

func GetDefaultDate() time.Time {
	result, _ := time.ParseInLocation(DateFormat, DefaultDate, time.Local)
	return result
}

func ToDate(t time.Time) Date {
	return Date(t.Format(DateFormat))
}

func ToTime(date string) (time.Time, error) {
	return time.Parse(DateFormat, date)
}

func ToUTC(t time.Time) Datetime {
	return Datetime(t.In(time.UTC).Format(time.RFC3339))
}
