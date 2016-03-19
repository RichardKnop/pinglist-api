package alarms

import (
	"fmt"

	"github.com/RichardKnop/pinglist-api/config"
	"github.com/RichardKnop/pinglist-api/email"
)

// TimeFormat specifies how the time will be parsed in emails
const TimeFormat = "Mon Jan _2 15:04:05 2006"

var alarmDownEmailTemplate = `
Hello %s,

Our system has noticed a downtime for one of your alarms:

%s is down since %s.

Kind Regards,

%s Team
`

var alarmUpEmailTemplate = `
Hello %s,

Our system has noticed the downtime of one of your alarms has been resolved.

%s is up again since %s after %d downtime.

Kind Regards,

%s Team
`

// EmailFactory facilitates construction of email.Email objects
type EmailFactory struct {
	cnf *config.Config
}

// NewEmailFactory starts a new emailFactory instance
func NewEmailFactory(cnf *config.Config) *EmailFactory {
	return &EmailFactory{cnf: cnf}
}

// NewAlarmDownEmail returns an alarm down notification email
func (f *EmailFactory) NewAlarmDownEmail(alarm *Alarm) *email.Email {
	// Define a greetings name for the user
	name := alarm.User.GetName()
	if name == "" {
		name = "friend"
	}

	// The email subject
	subject := fmt.Sprintf("ALERT: %s is down", alarm.EndpointURL)

	// Replace placeholders in the email template
	emailText := fmt.Sprintf(
		alarmDownEmailTemplate,
		name,
		alarm.EndpointURL,
		alarm.LastDowntimeStartedAt.Time.Format(TimeFormat),
		f.cnf.Web.Host,
	)

	return &email.Email{
		Subject: subject,
		Recipients: []*email.Recipient{&email.Recipient{
			Email: alarm.User.OauthUser.Username,
			Name:  alarm.User.GetName(),
		}},
		From: fmt.Sprintf("noreply@%s", f.cnf.Web.Host),
		Text: emailText,
	}
}

// NewAlarmUpEmail returns an alarm up notification email
func (f *EmailFactory) NewAlarmUpEmail(alarm *Alarm) *email.Email {
	// Define a greetings name for the user
	name := alarm.User.GetName()
	if name == "" {
		name = "friend"
	}

	// The email subject
	subject := fmt.Sprintf("ALERT: %s is up again", alarm.EndpointURL)

	// Replace placeholders in the email template
	downtime := alarm.LastUptimeStartedAt.Time.Sub(alarm.LastDowntimeStartedAt.Time)
	emailText := fmt.Sprintf(
		alarmUpEmailTemplate,
		name,
		alarm.EndpointURL,
		alarm.LastUptimeStartedAt.Time.Format(TimeFormat),
		fmt.Sprintf("%.2f minutes", downtime.Minutes()),
		f.cnf.Web.Host,
	)

	return &email.Email{
		Subject: subject,
		Recipients: []*email.Recipient{&email.Recipient{
			Email: alarm.User.OauthUser.Username,
			Name:  alarm.User.GetName(),
		}},
		From: fmt.Sprintf("noreply@%s", f.cnf.Web.Host),
		Text: emailText,
	}
}
