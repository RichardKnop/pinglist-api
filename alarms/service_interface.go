package alarms

import (
	"net/http"
	"time"

	"github.com/RichardKnop/pinglist-api/accounts"
)

// ServiceInterface defines exported methods
type ServiceInterface interface {
	// Exported methods
	GetAccountsService() accounts.ServiceInterface
	FindAlarmByID(alarmID uint) (*Alarm, error)
	GetAlarmsToCheck(now time.Time) ([]*Alarm, error)
	CheckAlarm(alarmID uint, watermark time.Time) error
	PartitionTable(parentTableName string, now time.Time) error
	RotateSubTables() error

	// Needed for the newRoutes to be able to register handlers
	listRegionsHandler(w http.ResponseWriter, r *http.Request)
	createAlarmHandler(w http.ResponseWriter, r *http.Request)
	updateAlarmHandler(w http.ResponseWriter, r *http.Request)
	deleteAlarmHandler(w http.ResponseWriter, r *http.Request)
	listAlarmsHandler(w http.ResponseWriter, r *http.Request)
	listAlarmIncidentsHandler(w http.ResponseWriter, r *http.Request)
	listAlarmResultsHandler(w http.ResponseWriter, r *http.Request)
}
