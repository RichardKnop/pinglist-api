package alarms

import (
	"time"

	"github.com/stretchr/testify/assert"
)

func (suite *AlarmsTestSuite) TestPaginatedResultsCount() {
	var (
		today             = time.Date(2016, time.February, 9, 0, 0, 0, 0, time.UTC)
		todaySubTableName = "alarm_results_2016_02_09"
		count             int
		err               error
	)

	// Partition the results table
	err = suite.service.PartitionTable(ResultParentTableName, today)
	assert.NoError(suite.T(), err, "Partitioning table failed")

	// Insert some test results
	testResults := []*Result{
		newResult(todaySubTableName, suite.alarms[0], today, 123),
		newResult(todaySubTableName, suite.alarms[0], today.Add(1*time.Hour), 234),
		newResult(todaySubTableName, suite.alarms[1], today, 345),
	}
	for _, result := range testResults {
		err := suite.db.Create(result).Error
		assert.NoError(suite.T(), err, "Inserting test data failed")
	}

	count, err = suite.service.paginatedResultsCount(nil)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 3, count)
	}

	count, err = suite.service.paginatedResultsCount(suite.alarms[0])
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 2, count)
	}

	count, err = suite.service.paginatedResultsCount(suite.alarms[1])
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 1, count)
	}

	count, err = suite.service.paginatedResultsCount(suite.alarms[2])
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 0, count)
	}
}

func (suite *AlarmsTestSuite) TestFindPaginatedResults() {
	var (
		today             = time.Date(2016, time.February, 9, 0, 0, 0, 0, time.UTC)
		todaySubTableName = "alarm_results_2016_02_09"
		results           []*Result
		err               error
	)

	// Partition the results table
	err = suite.service.PartitionTable(ResultParentTableName, today)
	assert.NoError(suite.T(), err, "Partitioning table failed")

	// Insert some test results
	testResults := []*Result{
		newResult(todaySubTableName, suite.alarms[0], today, 123),
		newResult(todaySubTableName, suite.alarms[0], today.Add(1*time.Hour), 234),
		newResult(todaySubTableName, suite.alarms[0], today.Add(2*time.Hour), 345),
		newResult(todaySubTableName, suite.alarms[0], today.Add(3*time.Hour), 456),
	}
	for _, result := range testResults {
		err := suite.db.Create(result).Error
		assert.NoError(suite.T(), err, "Inserting test data failed")
	}

	// This should return all results
	results, err = suite.service.findPaginatedResults(0, 25, "", nil)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 4, len(results))
		assert.Equal(suite.T(), testResults[0].ID, results[0].ID)
		assert.Equal(suite.T(), testResults[1].ID, results[1].ID)
		assert.Equal(suite.T(), testResults[2].ID, results[2].ID)
		assert.Equal(suite.T(), testResults[3].ID, results[3].ID)
	}

	// This should return all results ordered by ID desc
	results, err = suite.service.findPaginatedResults(0, 25, "id desc", nil)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 4, len(results))
		assert.Equal(suite.T(), testResults[3].ID, results[0].ID)
		assert.Equal(suite.T(), testResults[2].ID, results[1].ID)
		assert.Equal(suite.T(), testResults[1].ID, results[2].ID)
		assert.Equal(suite.T(), testResults[0].ID, results[3].ID)
	}

	// Test offset
	results, err = suite.service.findPaginatedResults(2, 25, "", nil)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 2, len(results))
		assert.Equal(suite.T(), testResults[2].ID, results[0].ID)
		assert.Equal(suite.T(), testResults[3].ID, results[1].ID)
	}

	// Test limit
	results, err = suite.service.findPaginatedResults(2, 1, "", nil)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 1, len(results))
		assert.Equal(suite.T(), testResults[2].ID, results[0].ID)
	}
}
