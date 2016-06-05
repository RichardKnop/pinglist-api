package alarms

import (
	"fmt"

	"github.com/RichardKnop/pinglist-api/alarms/incidenttypes"
	"github.com/RichardKnop/pinglist-api/config"
	"github.com/RichardKnop/pinglist-api/email"
)

// EmailTimeFormat specifies how the time will be parsed in emails
const EmailTimeFormat = "Mon Jan 2 15:04:05 2006"

var newIncidentEmailSubjectTemplates = map[string]string{
	incidenttypes.Slow:    "ALERT: %s returned slow response",
	incidenttypes.Timeout: "ALERT: %s timed out",
	incidenttypes.BadCode: "ALERT: %s returned bad status code",
	incidenttypes.Other:   "ALERT: %s failed for unknown reason",
}

var incidentResolvedEmailSubjectTemplate = "ALERT: %s is up and working correctly"

var newIncidentEmailTextTemplates = map[string]string{
	incidenttypes.Slow: `
Hello %s,

Our system has noticed a new incident with one of your alarms:

%s returned a slow response at %s [UTC].

Take a look at the incident dashboard: %s

Kind Regards,

%s Team
`,
	incidenttypes.Timeout: `
Hello %s,

Our system has noticed a new incident with one of your alarms:

%s timed out at %s [UTC].

Take a look at the incident dashboard: %s

Kind Regards,

%s Team
`,
	incidenttypes.BadCode: `
Hello %s,

Our system has noticed a new incident with one of your alarms:

%s returned a bad status code at %s [UTC].

Take a look at the incident dashboard: %s

Kind Regards,

%s Team
`,
	incidenttypes.Other: `
Hello %s,

Our system has noticed a new incident with one of your alarms:

%s failed for an unknown reason at %s [UTC].

Take a look at the incident dashboard: %s

Kind Regards,

%s Team
`,
}

var incidentResolvedEmailTextTemplate = `
Hello %s,

Our system has noticed a recent incident with one of your alarms has been resolved.

Since %s [UTC], %s is up and working correctly again after %s.

Take a look at the incident dashboard: %s

Kind Regards,

%s Team
`

// EmailFactory facilitates construction of email.Email objects
type EmailFactory struct {
	cnf *config.Config
}

// NewEmailFactory starts a new EmailFactory instance
func NewEmailFactory(cnf *config.Config) *EmailFactory {
	return &EmailFactory{cnf: cnf}
}

// NewIncidentEmail returns a new incident notification email
func (f *EmailFactory) NewIncidentEmail(incident *Incident) *email.Email {
	// Define a greetings name for the user
	name := incident.Alarm.User.GetName()
	if name == "" {
		name = "friend"
	}

	// The email subject
	subject := fmt.Sprintf(
		newIncidentEmailSubjectTemplates[incident.IncidentTypeID.String],
		incident.Alarm.EndpointURL,
	)

	// Dashboard incidents link
	incidentsLink := fmt.Sprintf(
		"%s://%s/alarms/%d/incidents/",
		f.cnf.Web.AppScheme,
		f.cnf.Web.AppHost,
		incident.Alarm.ID,
	)

	// Replace placeholders in the email template
	emailText := fmt.Sprintf(
		newIncidentEmailTextTemplates[incident.IncidentTypeID.String],
		name,
		incident.Alarm.EndpointURL,
		incident.Alarm.LastDowntimeStartedAt.Time.UTC().Format(EmailTimeFormat),
		incidentsLink,
		f.cnf.Web.AppHost,
	)

	return &email.Email{
		Subject: subject,
		Recipients: []*email.Recipient{&email.Recipient{
			Email: incident.Alarm.User.OauthUser.Username,
			Name:  incident.Alarm.User.GetName(),
		}},
		From: fmt.Sprintf("noreply@%s", f.cnf.Web.AppHost),
		Text: emailText,
	}
}

// NewIncidentsResolvedEmail returns an incidents resolved notification email
func (f *EmailFactory) NewIncidentsResolvedEmail(alarm *Alarm) *email.Email {
	// Define a greetings name for the user
	name := alarm.User.GetName()
	if name == "" {
		name = "friend"
	}

	// The email subject
	subject := fmt.Sprintf(incidentResolvedEmailSubjectTemplate, alarm.EndpointURL)

	// Downtime started at
	downtimeStartedAt := alarm.LastDowntimeStartedAt.Time.UTC().Format(EmailTimeFormat)

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

	// Replace placeholders in the email template
	emailText := fmt.Sprintf(
		incidentResolvedEmailTextTemplate,
		name,
		downtimeStartedAt,
		alarm.EndpointURL,
		downtime,
		incidentsLink,
		f.cnf.Web.AppHost,
	)

	return &email.Email{
		Subject: subject,
		Recipients: []*email.Recipient{&email.Recipient{
			Email: alarm.User.OauthUser.Username,
			Name:  alarm.User.GetName(),
		}},
		From: fmt.Sprintf("noreply@%s", f.cnf.Web.AppHost),
		Text: emailText,
	}
}
