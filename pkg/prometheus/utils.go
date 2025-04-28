package prometheus

import (
	"fmt"
	"strconv"
	"time"
)

func parseTime(timeStr string) (time.Time, error) {
	// Try parsing as RFC3339 first
	t, err := time.Parse(time.RFC3339, timeStr)
	if err == nil {
		return t, nil
	}

	// Try parsing as Unix timestamp (milliseconds)
	if unixTime, err := strconv.ParseInt(timeStr, 10, 64); err == nil {
		return time.Unix(unixTime/1000, (unixTime%1000)*1000000).UTC(), nil
	}

	return time.Time{}, fmt.Errorf("invalid time format, expected RFC3339 or Unix timestamp")
}
