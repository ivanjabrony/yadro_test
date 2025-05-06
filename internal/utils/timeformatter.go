package utils

import (
	"time"
)

const clockLayout = "15:04:05.000"
const durationLayout = "15:04:05"

func UnmarshallTimeStamp(timestamp string) (time.Time, error) {
	return time.Parse(clockLayout, timestamp)
}

func MarshallTimeStamp(formattedTime time.Time) string {
	return formattedTime.Format(clockLayout)
}

func UnmarshallDuration(duration string) (time.Duration, error) {
	var zTime, _ = time.Parse(durationLayout, "00:00:00")
	t, err := time.Parse(durationLayout, duration)
	if err != nil {
		return zTime.Sub(zTime), err
	}

	return t.Sub(zTime), nil
}

func MarshallDuration(formattedDuration time.Duration) string {
	var zTime, _ = time.Parse(durationLayout, "00:00:00")
	return zTime.Add(formattedDuration).Format(durationLayout)
}

func MarshallDurationToTimestamp(formattedDuration time.Duration) string {
	var zTime, _ = time.Parse(clockLayout, "00:00:00")
	return zTime.Add(formattedDuration).Format(clockLayout)
}
