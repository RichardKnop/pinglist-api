package scheduler

import (
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

// Start periodically runs goroutines to:
// - watch for scheduled alarms
// - partition alarm_results table & rotate old sub tables
func (s *Scheduler) Start(alarmsInterval, partitionInterval time.Duration) chan bool {
	// Partition / rotate metrics table once initially
	s.runPartitioningJob()

	// Stop channel
	stopped := make(chan bool, 1)

	// Tickers
	alarmsCheckTicker := time.NewTicker(alarmsInterval * time.Second)
	partitionTicker := time.NewTicker(partitionInterval * time.Second)

	go func() {
		for {
			select {
			case <-alarmsCheckTicker.C:
				go s.runAlarmCheckJob()
			case <-partitionTicker.C:
				go s.runPartitioningJob()
			case <-stopped:
				return
			}
		}
	}()

	return stopped
}

func (s *Scheduler) runAlarmCheckJob() {
	// Get alarms to check
	alarms, err := s.alarmsService.GetAlarmsToCheck(time.Now())
	if err != nil {
		logger.Error(err)
		return
	}

	// Any alarms to check
	if len(alarms) < 1 {
		logger.Info("No alarms to check")
		return
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

func (s *Scheduler) runPartitioningJob() {
	// Partition the request time metrics table
	err := s.metricsService.PartitionResponseTime(
		metrics.ResponseTimeParentTableName,
		time.Now(),
	)
	if err != nil {
		logger.Error(err)
		return
	}

	// Rotate old sub tables
	if err := s.metricsService.RotateSubTables(); err != nil {
		logger.Error(err)
		return
	}
}
