package commands

import (
	"time"

	"github.com/RichardKnop/pinglist-api/scheduler"
)

// RunAll runs the both the scheduler and the app
func RunAll() error {
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

	// Init the app
	app, err := initApp(cnf, db)
	if err != nil {
		return err
	}

	// Run the scheduling goroutines
	alarmsInterval := time.Duration(10)     // alarms check interval = 10s
	partitionInterval := time.Duration(600) // partition / rotate interval = 10m
	stoppedChan := theScheduler.Start(alarmsInterval, partitionInterval)

	// Run the server on port 8080
	app.Run(":8080")

	// Stop the scheduler
	stoppedChan <- true

	return nil
}
