package durations

import "time"

// GetDurationWithDefault gets either the value configured without changing anything
// OR the default provided in case the value is 0
func GetDurationWithDefault(value time.Duration, defaultToSet time.Duration) time.Duration {
	var warmupToSet time.Duration
	if value > 0 {
		warmupToSet = value
	} else {
		warmupToSet = defaultToSet
	}
	return warmupToSet
}
