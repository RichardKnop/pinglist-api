package alarms

import (
	"testing"
	"time"

	"github.com/RichardKnop/pinglist-api/accounts"
	"github.com/RichardKnop/pinglist-api/alarms/incidenttypes"
	"github.com/RichardKnop/pinglist-api/config"
	"github.com/RichardKnop/pinglist-api/oauth"
	"github.com/RichardKnop/pinglist-api/util"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func TestNewIncidentEmailSlow(t *testing.T) {
	emailFactory := NewEmailFactory(&config.Config{
		Web: config.WebConfig{
			Scheme:    "https",
			Host:      "api.pingli.st",
			AppScheme: "https",
			AppHost:   "pingli.st",
		},
	})

	lastDowntimeStartedAt := time.Date(
		2016, // year
		6,    // month
		4,    // day
		11,   // hour
		26,   // minute
		15,   // second
		1234, // nanosecond
		time.FixedZone("HKT", 8*3600), // timezone
	)

	incident := &Incident{
		IncidentTypeID: util.StringOrNull(incidenttypes.Slow),
		Alarm: &Alarm{
			Model: gorm.Model{ID: 123},
			User: &accounts.User{
				OauthUser: &oauth.User{
					Username: "john@reese",
				},
				FirstName: util.StringOrNull("John"),
				LastName:  util.StringOrNull("Reese"),
			},
			EndpointURL:           "http://endpoint-url",
			LastDowntimeStartedAt: util.TimeOrNull(&lastDowntimeStartedAt),
		},
	}
	email := emailFactory.NewIncidentEmail(incident)

	assert.Equal(t, "ALERT: http://endpoint-url returned slow response", email.Subject)
	assert.Equal(t, 1, len(email.Recipients))
	assert.Equal(t, "john@reese", email.Recipients[0].Email)
	assert.Equal(t, "John Reese", email.Recipients[0].Name)
	assert.Equal(t, "noreply@pingli.st", email.From)

	expectedText := `
Hello John Reese,

Our system has noticed a new incident with one of your alarms:

http://endpoint-url returned a slow response at Sat Jun 4 03:26:15 2016 [UTC].

Take a look at the incident dashboard: https://pingli.st/alarms/123/incidents/

Kind Regards,

pingli.st Team
`
	assert.Equal(t, expectedText, email.Text)
}

func TestNewIncidentEmailTimeout(t *testing.T) {
	emailFactory := NewEmailFactory(&config.Config{
		Web: config.WebConfig{
			Scheme:    "https",
			Host:      "api.pingli.st",
			AppScheme: "https",
			AppHost:   "pingli.st",
		},
	})

	lastDowntimeStartedAt := time.Date(
		2016, // year
		6,    // month
		4,    // day
		11,   // hour
		26,   // minute
		15,   // second
		1234, // nanosecond
		time.FixedZone("HKT", 8*3600), // timezone
	)

	incident := &Incident{
		IncidentTypeID: util.StringOrNull(incidenttypes.Timeout),
		Alarm: &Alarm{
			Model: gorm.Model{ID: 123},
			User: &accounts.User{
				OauthUser: &oauth.User{
					Username: "john@reese",
				},
				FirstName: util.StringOrNull("John"),
				LastName:  util.StringOrNull("Reese"),
			},
			EndpointURL:           "http://endpoint-url",
			LastDowntimeStartedAt: util.TimeOrNull(&lastDowntimeStartedAt),
		},
	}
	email := emailFactory.NewIncidentEmail(incident)

	assert.Equal(t, "ALERT: http://endpoint-url timed out", email.Subject)
	assert.Equal(t, 1, len(email.Recipients))
	assert.Equal(t, "john@reese", email.Recipients[0].Email)
	assert.Equal(t, "John Reese", email.Recipients[0].Name)
	assert.Equal(t, "noreply@pingli.st", email.From)

	expectedText := `
Hello John Reese,

Our system has noticed a new incident with one of your alarms:

http://endpoint-url timed out at Sat Jun 4 03:26:15 2016 [UTC].

Take a look at the incident dashboard: https://pingli.st/alarms/123/incidents/

Kind Regards,

pingli.st Team
`
	assert.Equal(t, expectedText, email.Text)
}

func TestNewIncidentEmailBadCode(t *testing.T) {
	emailFactory := NewEmailFactory(&config.Config{
		Web: config.WebConfig{
			Scheme:    "https",
			Host:      "api.pingli.st",
			AppScheme: "https",
			AppHost:   "pingli.st",
		},
	})

	lastDowntimeStartedAt := time.Date(
		2016, // year
		6,    // month
		4,    // day
		11,   // hour
		26,   // minute
		15,   // second
		1234, // nanosecond
		time.FixedZone("HKT", 8*3600), // timezone
	)

	incident := &Incident{
		IncidentTypeID: util.StringOrNull(incidenttypes.BadCode),
		Alarm: &Alarm{
			Model: gorm.Model{ID: 123},
			User: &accounts.User{
				OauthUser: &oauth.User{
					Username: "john@reese",
				},
				FirstName: util.StringOrNull("John"),
				LastName:  util.StringOrNull("Reese"),
			},
			EndpointURL:           "http://endpoint-url",
			LastDowntimeStartedAt: util.TimeOrNull(&lastDowntimeStartedAt),
		},
	}
	email := emailFactory.NewIncidentEmail(incident)

	assert.Equal(t, "ALERT: http://endpoint-url returned bad status code", email.Subject)
	assert.Equal(t, 1, len(email.Recipients))
	assert.Equal(t, "john@reese", email.Recipients[0].Email)
	assert.Equal(t, "John Reese", email.Recipients[0].Name)
	assert.Equal(t, "noreply@pingli.st", email.From)

	expectedText := `
Hello John Reese,

Our system has noticed a new incident with one of your alarms:

http://endpoint-url returned a bad status code at Sat Jun 4 03:26:15 2016 [UTC].

Take a look at the incident dashboard: https://pingli.st/alarms/123/incidents/

Kind Regards,

pingli.st Team
`
	assert.Equal(t, expectedText, email.Text)
}

func TestNewIncidentEmailOther(t *testing.T) {
	emailFactory := NewEmailFactory(&config.Config{
		Web: config.WebConfig{
			Scheme:    "https",
			Host:      "api.pingli.st",
			AppScheme: "https",
			AppHost:   "pingli.st",
		},
	})

	lastDowntimeStartedAt := time.Date(
		2016, // year
		6,    // month
		4,    // day
		11,   // hour
		26,   // minute
		15,   // second
		1234, // nanosecond
		time.FixedZone("HKT", 8*3600), // timezone
	)

	incident := &Incident{
		IncidentTypeID: util.StringOrNull(incidenttypes.Other),
		Alarm: &Alarm{
			Model: gorm.Model{ID: 123},
			User: &accounts.User{
				OauthUser: &oauth.User{
					Username: "john@reese",
				},
				FirstName: util.StringOrNull("John"),
				LastName:  util.StringOrNull("Reese"),
			},
			EndpointURL:           "http://endpoint-url",
			LastDowntimeStartedAt: util.TimeOrNull(&lastDowntimeStartedAt),
		},
	}
	email := emailFactory.NewIncidentEmail(incident)

	assert.Equal(t, "ALERT: http://endpoint-url failed for unknown reason", email.Subject)
	assert.Equal(t, 1, len(email.Recipients))
	assert.Equal(t, "john@reese", email.Recipients[0].Email)
	assert.Equal(t, "John Reese", email.Recipients[0].Name)
	assert.Equal(t, "noreply@pingli.st", email.From)

	expectedText := `
Hello John Reese,

Our system has noticed a new incident with one of your alarms:

http://endpoint-url failed for an unknown reason at Sat Jun 4 03:26:15 2016 [UTC].

Take a look at the incident dashboard: https://pingli.st/alarms/123/incidents/

Kind Regards,

pingli.st Team
`
	assert.Equal(t, expectedText, email.Text)
}

func TestIncidentsResolved(t *testing.T) {
	emailFactory := NewEmailFactory(&config.Config{
		Web: config.WebConfig{
			Scheme:    "https",
			Host:      "api.pingli.st",
			AppScheme: "https",
			AppHost:   "pingli.st",
		},
	})

	lastDowntimeStartedAt := time.Date(
		2016, // year
		6,    // month
		4,    // day
		11,   // hour
		26,   // minute
		15,   // second
		1234, // nanosecond
		time.FixedZone("HKT", 8*3600), // timezone
	)

	lastUptimeStartedAt := time.Date(
		2016, // year
		6,    // month
		4,    // day
		11,   // hour
		27,   // minute
		48,   // second
		8201, // nanosecond
		time.FixedZone("HKT", 8*3600), // timezone
	)

	alarm := &Alarm{
		Model: gorm.Model{ID: 123},
		User: &accounts.User{
			OauthUser: &oauth.User{
				Username: "john@reese",
			},
			FirstName: util.StringOrNull("John"),
			LastName:  util.StringOrNull("Reese"),
		},
		EndpointURL:           "http://endpoint-url",
		LastDowntimeStartedAt: util.TimeOrNull(&lastDowntimeStartedAt),
		LastUptimeStartedAt:   util.TimeOrNull(&lastUptimeStartedAt),
	}
	email := emailFactory.NewIncidentsResolvedEmail(alarm)

	assert.Equal(t, "ALERT: http://endpoint-url is up and working correctly", email.Subject)
	assert.Equal(t, 1, len(email.Recipients))
	assert.Equal(t, "john@reese", email.Recipients[0].Email)
	assert.Equal(t, "John Reese", email.Recipients[0].Name)
	assert.Equal(t, "noreply@pingli.st", email.From)

	expectedText := `
Hello John Reese,

Our system has noticed a recent incident with one of your alarms has been resolved.

Since Sat Jun 4 03:26:15 2016 [UTC], http://endpoint-url is up and working correctly again after 1.55 minutes.

Take a look at the incident dashboard: https://pingli.st/alarms/123/incidents/

Kind Regards,

pingli.st Team
`
	assert.Equal(t, expectedText, email.Text)
}
