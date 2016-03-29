package commands

import "time"

// RunAll runs the both the scheduler and the app
func RunAll() error {
	// Init config and database
	cnf, db, err := initConfigDB(true, true)
	if err != nil {
		return err
	}
	defer db.Close()

	// Init the scheduler
	theScheduler, err := initScheduler(cnf, db)
	if err != nil {
		return err
	}

	// Init the app
	app, err := initApp(cnf, db)
	if err != nil {
		return err
	}

	// Run the scheduling goroutines
	_ = theScheduler.Run(
		time.Duration(10),  // alarms check interval = 10s
		time.Duration(600), // partition / rotate interval = 10m
	)

	// Run the server on port 8080
	app.Run(":8080")

	return nil
}
