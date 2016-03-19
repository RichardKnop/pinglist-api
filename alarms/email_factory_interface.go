package alarms

import (
	"github.com/RichardKnop/pinglist-api/email"
)

// EmailFactoryInterface defines exported methods
type EmailFactoryInterface interface {
	NewAlarmDownEmail(alarm *Alarm) *email.Email
	NewAlarmUpEmail(alarm *Alarm) *email.Email
}
