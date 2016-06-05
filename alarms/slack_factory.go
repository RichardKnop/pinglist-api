package alarms

import (
	"fmt"

	"github.com/RichardKnop/pinglist-api/alarms/incidenttypes"
	"github.com/RichardKnop/pinglist-api/config"
)

var (
	slackNotificationsUsername = "webhookbot"
	slackNotificationsEmoji    = ""
)

// SlackTimeFormat specifies how the time will be parsed in Slack messages
const SlackTimeFormat = "Mon Jan 2 15:04:05 2006"

var newIncidentSlackMessageTemplates = map[string]string{
	incidenttypes.Slow: `
Our system has noticed a new incident with one of your alarms:

%s returned a slow response at %s [UTC].

Take a look at the incident dashboard: %s
`,
	incidenttypes.Timeout: `
Our system has noticed a new incident with one of your alarms:

%s timed out at %s [UTC].

Take a look at the incident dashboard: %s
`,
	incidenttypes.BadCode: `
Our system has noticed a new incident with one of your alarms:

%s returned a bad status code at %s [UTC].

Take a look at the incident dashboard: %s
`,
	incidenttypes.Other: `
Our system has noticed a new incident with one of your alarms:

%s failed for an unknown reason at %s [UTC].

Take a look at the incident dashboard: %s
`,
}

var incidentResolvedSlackMessageTemplate = `
Our system has noticed a recent incident with one of your alarms has been resolved.

Since %s [UTC], %s is up and working correctly again after %s.

Take a look at the incident dashboard: %s
`

// SlackFactory facilitates construction of Slack messages
type SlackFactory struct {
	cnf *config.Config
}

// NewSlackFactory starts a new SlackFactory instance
func NewSlackFactory(cnf *config.Config) *SlackFactory {
	return &SlackFactory{cnf: cnf}
}

// NewIncidentMessage returns a new incident notification message
func (f *SlackFactory) NewIncidentMessage(incident *Incident) string {
	// Dashboard incidents link
	incidentsLink := fmt.Sprintf(
		"%s://%s/alarms/%d/incidents/",
		f.cnf.Web.AppScheme,
		f.cnf.Web.AppHost,
		incident.Alarm.ID,
	)

	return fmt.Sprintf(
		newIncidentSlackMessageTemplates[incident.IncidentTypeID.String],
		incident.Alarm.EndpointURL,
		incident.Alarm.LastDowntimeStartedAt.Time.UTC().Format(SlackTimeFormat),
		incidentsLink,
	)
}

// NewIncidentsResolvedMessage returns an incidents resolved notification message
func (f *SlackFactory) NewIncidentsResolvedMessage(alarm *Alarm) string {
	// Downtime started at
	downtimeStartedAt := alarm.LastDowntimeStartedAt.Time.UTC().Format(SlackTimeFormat)

	// Downtime
	downtime := fmt.Sprintf(
		"%.2f minutes",
		alarm.LastUptimeStartedAt.Time.Sub(alarm.LastDowntimeStartedAt.Time.UTC()).Minutes(),
	)

	// Dashboard incidents link
	incidentsLink := fmt.Sprintf(
		"%s://%s/alarms/%d/incidents/",
		f.cnf.Web.AppScheme,
		f.cnf.Web.AppHost,
		alarm.ID,
	)

	return fmt.Sprintf(
		incidentResolvedSlackMessageTemplate,
		downtimeStartedAt,
		alarm.EndpointURL,
		downtime,
		incidentsLink,
	)
}
