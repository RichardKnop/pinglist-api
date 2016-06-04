package alarms

import (
	"github.com/RichardKnop/pinglist-api/email"
)

// EmailFactoryInterface defines exported methods
type EmailFactoryInterface interface {
	NewIncidentEmail(incident *Incident) *email.Email
	NewIncidentsResolvedEmail(alarm *Alarm) *email.Email
}
