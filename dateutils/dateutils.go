package dateutils

import "time"

const (
	apiDateLayout   = "2006-01-02T15:04:05Z"
	apiDbDateLayout = "2006-01-02 15:04:05"
)

// GetNow is a function to encapsulate the recover of a time.Time as UTC (time.Now().UTC()).
func GetNow() time.Time {
	return time.Now().UTC()
}

// GetNowString is a function to encapsulate the result of GetNow with the pattern 2006-01-02T15:04:05Z.
func GetNowString() string {
	return GetNow().Format(apiDateLayout)
}

// GetNowDBString is a function to encapsulate the result of GetNow with the pattern prepared for DB.
func GetNowDBString() string {
	return GetNow().Format(apiDbDateLayout)
}
