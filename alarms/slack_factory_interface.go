package alarms

// SlackFactoryInterface defines exported methods
type SlackFactoryInterface interface {
	NewIncidentMessage(incident *Incident) string
	NewIncidentsResolvedMessage(alarm *Alarm) string
}
