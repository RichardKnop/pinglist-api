package scheduler

import (
	"sync"
	"time"

	"github.com/RichardKnop/pinglist-api/alarms"
)

// Scheduler ...
type Scheduler struct {
	alarmsService alarms.ServiceInterface
}

// New starts a new Scheduler instance
func New(alarmsService alarms.ServiceInterface) *Scheduler {
	return &Scheduler{
		alarmsService: alarmsService,
	}
}

// Run opens individial goroutines to:
// - watch for scheduled alarms
// - partition alarm_results table & rotate old sub tables
func (s *Scheduler) Run(alarmsInterval, partitionInterval time.Duration) {
	var wg sync.WaitGroup

	// Watch for scheduled alarms
	wg.Add(1)
	go func() {
		for {
			// Wait
			time.Sleep(time.Second * alarmsInterval)

			// Get alarms to check
			alarms, err := s.alarmsService.GetAlarmsToCheck(time.Now())
			if err != nil {
				logger.Error(err)
				continue
			}

			// Any alarms to check
			if len(alarms) < 1 {
				logger.Info("No alarms to check")
				continue
			}

			// Create a new time object as a watermark
			now := time.Now()

			// Iterate over alarms and fire check goroutines
			logger.Infof("Triggerring %d alarm checks", len(alarms))
			for _, alarm := range alarms {
				go func(alarmID uint, watermark time.Time) {
					if err := s.alarmsService.CheckAlarm(alarmID, watermark); err != nil {
						logger.Error(err)
					}
				}(alarm.ID, now)
			}
		}
	}()

	// Partition alarm_results table and rotate old sub tables
	wg.Add(1)
	go func() {

		// Wait
		time.Sleep(time.Second * partitionInterval)

		// Partition the alarm_results table
		if err := s.alarmsService.PartitionTable(alarms.ResultParentTableName, time.Now()); err != nil {
			logger.Error(err)
		}

		// Rotate old sub tables
		if err := s.alarmsService.RotateSubTables(); err != nil {
			logger.Error(err)
		}
	}()

	wg.Wait()
}
