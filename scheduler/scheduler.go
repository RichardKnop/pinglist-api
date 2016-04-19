package scheduler

import (
	"sync"
	"time"

	"github.com/RichardKnop/pinglist-api/alarms"
	"github.com/RichardKnop/pinglist-api/metrics"
)

// Scheduler ...
type Scheduler struct {
	metricsService metrics.ServiceInterface
	alarmsService  alarms.ServiceInterface
}

// New starts a new Scheduler instance
func New(metricsService metrics.ServiceInterface, alarmsService alarms.ServiceInterface) *Scheduler {
	return &Scheduler{
		metricsService: metricsService,
		alarmsService:  alarmsService,
	}
}

// Run opens individial goroutines to:
// - watch for scheduled alarms
// - partition alarm_results table & rotate old sub tables
func (s *Scheduler) Run(alarmsInterval, partitionInterval time.Duration) sync.WaitGroup {
	var wg sync.WaitGroup

	// Watch for scheduled alarm checks
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			// Wait before repeating
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
		defer wg.Done()

		// Partition the request time metrics table
		err := s.metricsService.PartitionResponseTime(
			metrics.ResponseTimeParentTableName,
			time.Now(),
		)
		if err != nil {
			logger.Error(err)
		}

		// Rotate old sub tables
		if err := s.metricsService.RotateSubTables(); err != nil {
			logger.Error(err)
		}

		// Wait before repeating
		time.Sleep(time.Second * partitionInterval)
	}()

	return wg
}
