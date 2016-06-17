package commands

import (
	"time"

	"github.com/RichardKnop/pinglist-api/scheduler"
)

// RunScheduler runs the scheduler
func RunScheduler() error {
	cnf, db, err := initConfigDB(true, true)
	if err != nil {
		return err
	}
	defer db.Close()
	if err := initServices(cnf, db); err != nil {
		return err
	}

	// Init the scheduler
	theScheduler := scheduler.New(metricsService, alarmsService)

	// Run the scheduling goroutines
	alarmsInterval := time.Duration(10)     // alarms check interval = 10s
	partitionInterval := time.Duration(600) // partition / rotate interval = 10m
	<-theScheduler.Start(alarmsInterval, partitionInterval)

	return nil
}
