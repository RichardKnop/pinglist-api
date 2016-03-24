package timeseries

import (
	"time"

	"github.com/stretchr/testify/assert"
)

func (suite *TimeseriesTestSuite) TestPaginatedRequestTimesCount() {
	var (
		today             = time.Date(2016, time.February, 9, 0, 0, 0, 0, time.UTC)
		todaySubTableName = "timeseries_request_times_2016_02_09"
		count             int
		err               error
	)

	// Partition the request time table
	err = suite.service.PartitionRequestTime(RequestTimeParentTableName, today)
	assert.NoError(suite.T(), err, "Partitioning table failed")

	// Insert some test records
	testRecords := []*RequestTime{
		newRequestTime(todaySubTableName, 1, today, 123),
		newRequestTime(todaySubTableName, 1, today.Add(1*time.Hour), 234),
		newRequestTime(todaySubTableName, 2, today, 345),
	}
	for _, testRecord := range testRecords {
		err := suite.db.Create(testRecord).Error
		assert.NoError(suite.T(), err, "Inserting test data failed")
	}

	count, err = suite.service.PaginatedRequestTimesCount(0)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 3, count)
	}

	count, err = suite.service.PaginatedRequestTimesCount(1)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 2, count)
	}

	count, err = suite.service.PaginatedRequestTimesCount(2)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 1, count)
	}

	count, err = suite.service.PaginatedRequestTimesCount(3)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 0, count)
	}
}

func (suite *TimeseriesTestSuite) TestFindPaginatedRequestTimes() {
	var (
		today             = time.Date(2016, time.February, 9, 0, 0, 0, 0, time.UTC)
		todaySubTableName = "timeseries_request_times_2016_02_09"
		requestTimes      []*RequestTime
		err               error
	)

	// Partition the request time table
	err = suite.service.PartitionRequestTime(RequestTimeParentTableName, today)
	assert.NoError(suite.T(), err, "Partitioning table failed")

	// Insert some test records
	testRecords := []*RequestTime{
		newRequestTime(todaySubTableName, 1, today, 123),
		newRequestTime(todaySubTableName, 1, today.Add(1*time.Hour), 234),
		newRequestTime(todaySubTableName, 1, today.Add(2*time.Hour), 345),
		newRequestTime(todaySubTableName, 1, today.Add(3*time.Hour), 456),
	}
	for _, testRecord := range testRecords {
		err := suite.db.Create(testRecord).Error
		assert.NoError(suite.T(), err, "Inserting test data failed")
	}

	// This should return all records
	requestTimes, err = suite.service.FindPaginatedRequestTimes(
		0,  // offset
		25, // limit
		"", // order by
		0,  // reference ID
	)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 4, len(requestTimes))
		assert.Equal(suite.T(), testRecords[0].Timestamp.Unix(), requestTimes[0].Timestamp.Unix())
		assert.Equal(suite.T(), testRecords[1].Timestamp.Unix(), requestTimes[1].Timestamp.Unix())
		assert.Equal(suite.T(), testRecords[2].Timestamp.Unix(), requestTimes[2].Timestamp.Unix())
		assert.Equal(suite.T(), testRecords[3].Timestamp.Unix(), requestTimes[3].Timestamp.Unix())
	}

	// This should return all records ordered by timestamp desc
	requestTimes, err = suite.service.FindPaginatedRequestTimes(
		0,                // offset
		25,               // limit
		"timestamp desc", // order by
		0,                // reference ID
	)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 4, len(requestTimes))
		assert.Equal(suite.T(), testRecords[3].Timestamp.Unix(), requestTimes[0].Timestamp.Unix())
		assert.Equal(suite.T(), testRecords[2].Timestamp.Unix(), requestTimes[1].Timestamp.Unix())
		assert.Equal(suite.T(), testRecords[1].Timestamp.Unix(), requestTimes[2].Timestamp.Unix())
		assert.Equal(suite.T(), testRecords[0].Timestamp.Unix(), requestTimes[3].Timestamp.Unix())
	}

	// Test offset
	requestTimes, err = suite.service.FindPaginatedRequestTimes(
		2,  // offset
		25, // limit
		"", // order by
		0,  // reference ID
	)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 2, len(requestTimes))
		assert.Equal(suite.T(), testRecords[2].Timestamp.Unix(), requestTimes[0].Timestamp.Unix())
		assert.Equal(suite.T(), testRecords[3].Timestamp.Unix(), requestTimes[1].Timestamp.Unix())
	}

	// Test limit
	requestTimes, err = suite.service.FindPaginatedRequestTimes(
		2,  // offset
		1,  // limit
		"", // order by
		0,  // reference ID
	)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 1, len(requestTimes))
		assert.Equal(suite.T(), testRecords[2].Timestamp.Unix(), requestTimes[0].Timestamp.Unix())
	}
}
