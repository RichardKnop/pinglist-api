package alarmstates

const (
	// OK - The metric is within the defined threshold
	OK = "ok"
	// Alarm - The metric is outside of the defined threshold
	Alarm = "alarm"
	// InsufficientData - The alarm has just started, the metric is not available,
	// or not enough data is available for the metric to determine the alarm state
	InsufficientData = "insufficient_data"
)
