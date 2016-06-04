package alarms

import (
	"github.com/RichardKnop/pinglist-api/alarms/incidenttypes"
)

var newIncidentPushNotificationTemplates = map[string]string{
	incidenttypes.Slow:    "ALERT: %s returned slow response",
	incidenttypes.Timeout: "ALERT: %s timed out",
	incidenttypes.BadCode: "ALERT: %s returned bad status code",
	incidenttypes.Other:   "ALERT: %s failed for unknown reason",
}

var incidentsResolvedPushNotificationTemplate = "ALERT: %s is up and working correctly"
