package metrics

import (
	"time"

	"github.com/stretchr/testify/assert"
)

func (suite *MetricsTestSuite) TestPaginatedResponseTimesCount() {
	var (
		today                 = time.Date(2016, time.February, 9, 0, 0, 0, 0, time.UTC)
		todaySubTableName     = "metrics_response_times_2016_02_09"
		yesterday             = time.Date(2016, time.February, 8, 0, 0, 0, 0, time.UTC)
		yesterdaySubTableName = "metrics_response_times_2016_02_08"
		from, to              time.Time
		count                 int
		err                   error
	)

	// Partition the response time table
	err = suite.service.PartitionResponseTime(ResponseTimeParentTableName, today)
	assert.NoError(suite.T(), err, "Partitioning table failed")
	err = suite.service.PartitionResponseTime(ResponseTimeParentTableName, yesterday)
	assert.NoError(suite.T(), err, "Partitioning table failed")

	// Insert some test records
	testRecords := []*ResponseTime{
		NewResponseTime(yesterdaySubTableName, 1, yesterday, 123),
		NewResponseTime(yesterdaySubTableName, 1, yesterday.Add(1*time.Hour), 234),
		NewResponseTime(yesterdaySubTableName, 2, yesterday.Add(2*time.Hour), 345),
		NewResponseTime(todaySubTableName, 1, today, 123),
		NewResponseTime(todaySubTableName, 1, today.Add(1*time.Hour), 234),
		NewResponseTime(todaySubTableName, 2, today.Add(2*time.Hour), 345),
	}
	for _, testRecord := range testRecords {
		err := suite.db.Create(testRecord).Error
		assert.NoError(suite.T(), err, "Inserting test data failed")
	}

	// No filtering at all
	count, err = suite.service.PaginatedResponseTimesCount(
		0,   // reference ID
		"",  // date trunc
		nil, // from
		nil, // to
	)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 6, count)
	}

	// Filter by a valid reference ID
	count, err = suite.service.PaginatedResponseTimesCount(
		1,   // reference ID
		"",  // date trunc
		nil, // from
		nil, // to
	)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 4, count)
	}

	// Filter by a valid reference ID
	count, err = suite.service.PaginatedResponseTimesCount(
		2,   // reference ID
		"",  // date trunc
		nil, // from
		nil, // to
	)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 2, count)
	}

	// Filter by a bogus reference ID
	count, err = suite.service.PaginatedResponseTimesCount(
		3,   // reference ID
		"",  // date trunc
		nil, // from
		nil, // to
	)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 0, count)
	}

	// Filter by "from" and "to" timestamp
	from = yesterday.Add(2 * time.Hour)
	to = today
	count, err = suite.service.PaginatedResponseTimesCount(
		0,     // reference ID
		"",    // date trunc
		&from, // from
		&to,   // to
	)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 2, count)
	}

	// Filter by "date_trunc" timestamp
	count, err = suite.service.PaginatedResponseTimesCount(
		0,     // reference ID
		"day", // date trunc
		nil,   // from
		nil,   // to
	)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 2, count)
	}
}

func (suite *MetricsTestSuite) TestFindPaginatedResponseTimes() {
	var (
		today                 = time.Date(2016, time.February, 9, 0, 0, 0, 0, time.UTC)
		todaySubTableName     = "metrics_response_times_2016_02_09"
		yesterday             = time.Date(2016, time.February, 8, 0, 0, 0, 0, time.UTC)
		yesterdaySubTableName = "metrics_response_times_2016_02_08"
		from, to              time.Time
		ResponseTimes         []*ResponseTime
		err                   error
	)

	// Partition the request time table
	err = suite.service.PartitionResponseTime(ResponseTimeParentTableName, today)
	assert.NoError(suite.T(), err, "Partitioning table failed")
	err = suite.service.PartitionResponseTime(ResponseTimeParentTableName, yesterday)
	assert.NoError(suite.T(), err, "Partitioning table failed")

	// Insert some test records
	testRecords := []*ResponseTime{
		NewResponseTime(yesterdaySubTableName, 1, yesterday, 123),
		NewResponseTime(yesterdaySubTableName, 1, yesterday.Add(1*time.Hour), 234),
		NewResponseTime(yesterdaySubTableName, 2, yesterday.Add(2*time.Hour), 345),
		NewResponseTime(todaySubTableName, 1, today, 456),
		NewResponseTime(todaySubTableName, 1, today.Add(1*time.Hour), 567),
		NewResponseTime(todaySubTableName, 2, today.Add(2*time.Hour), 678),
	}
	for _, testRecord := range testRecords {
		err := suite.db.Create(testRecord).Error
		assert.NoError(suite.T(), err, "Inserting test data failed")
	}

	// No filtering at all
	ResponseTimes, err = suite.service.FindPaginatedResponseTimes(
		0,   // offset
		25,  // limit
		"",  // order by
		0,   // reference ID
		"",  // date trunc
		nil, // from
		nil, // to
	)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 6, len(ResponseTimes))
	}

	// Filter by a valid reference ID
	ResponseTimes, err = suite.service.FindPaginatedResponseTimes(
		0,   // offset
		25,  // limit
		"",  // order by
		1,   // reference ID
		"",  // date trunc
		nil, // from
		nil, // to
	)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 4, len(ResponseTimes))
		assert.Equal(suite.T(), testRecords[0].Timestamp.Unix(), ResponseTimes[0].Timestamp.Unix())
		assert.Equal(suite.T(), testRecords[1].Timestamp.Unix(), ResponseTimes[1].Timestamp.Unix())
		assert.Equal(suite.T(), testRecords[3].Timestamp.Unix(), ResponseTimes[2].Timestamp.Unix())
		assert.Equal(suite.T(), testRecords[4].Timestamp.Unix(), ResponseTimes[3].Timestamp.Unix())
	}

	// Filter by a valid reference ID
	ResponseTimes, err = suite.service.FindPaginatedResponseTimes(
		0,   // offset
		25,  // limit
		"",  // order by
		2,   // reference ID
		"",  // date trunc
		nil, // from
		nil, // to
	)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 2, len(ResponseTimes))
		assert.Equal(suite.T(), testRecords[2].Timestamp.Unix(), ResponseTimes[0].Timestamp.Unix())
		assert.Equal(suite.T(), testRecords[5].Timestamp.Unix(), ResponseTimes[1].Timestamp.Unix())
	}

	// Filter by a bogus reference ID
	ResponseTimes, err = suite.service.FindPaginatedResponseTimes(
		0,   // offset
		25,  // limit
		"",  // order by
		3,   // reference ID
		"",  // date trunc
		nil, // from
		nil, // to
	)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 0, len(ResponseTimes))
	}

	// Filter by "from" and "to" timestamp
	from = yesterday.Add(2 * time.Hour)
	to = today
	ResponseTimes, err = suite.service.FindPaginatedResponseTimes(
		0,     // offset
		25,    // limit
		"",    // order by
		0,     // reference ID
		"",    // date trunc
		&from, // from
		&to,   // to
	)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 2, len(ResponseTimes))
		assert.Equal(suite.T(), 2, len(ResponseTimes))
		assert.Equal(suite.T(), testRecords[2].Timestamp.Unix(), ResponseTimes[0].Timestamp.Unix())
		assert.Equal(suite.T(), testRecords[3].Timestamp.Unix(), ResponseTimes[1].Timestamp.Unix())
	}

	// Filter by "date_trunc" timestamp
	ResponseTimes, err = suite.service.FindPaginatedResponseTimes(
		0,     // offset
		25,    // limit
		"",    // order by
		0,     // reference ID
		"day", // date trunc
		nil,   // from
		nil,   // to
	)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 2, len(ResponseTimes))
		assert.Equal(
			suite.T(),
			yesterday.UTC().Format(time.RFC3339),
			ResponseTimes[0].Timestamp.UTC().Format(time.RFC3339),
		)
		assert.Equal(
			suite.T(),
			today.UTC().Format(time.RFC3339),
			ResponseTimes[1].Timestamp.UTC().Format(time.RFC3339),
		)
		assert.Equal(suite.T(), int64(234), ResponseTimes[0].Value) // (123 + 234 + 345) / 3
		assert.Equal(suite.T(), int64(567), ResponseTimes[1].Value) // (456 + 567 + 678) / 3
	}

	// This should return all records
	ResponseTimes, err = suite.service.FindPaginatedResponseTimes(
		0,   // offset
		25,  // limit
		"",  // order by
		0,   // reference ID
		"",  // date trunc
		nil, // from
		nil, // to
	)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 6, len(ResponseTimes))
		assert.Equal(suite.T(), testRecords[0].Timestamp.Unix(), ResponseTimes[0].Timestamp.Unix())
		assert.Equal(suite.T(), testRecords[1].Timestamp.Unix(), ResponseTimes[1].Timestamp.Unix())
		assert.Equal(suite.T(), testRecords[2].Timestamp.Unix(), ResponseTimes[2].Timestamp.Unix())
		assert.Equal(suite.T(), testRecords[3].Timestamp.Unix(), ResponseTimes[3].Timestamp.Unix())
		assert.Equal(suite.T(), testRecords[4].Timestamp.Unix(), ResponseTimes[4].Timestamp.Unix())
		assert.Equal(suite.T(), testRecords[5].Timestamp.Unix(), ResponseTimes[5].Timestamp.Unix())
	}

	// This should return all records ordered by timestamp desc
	ResponseTimes, err = suite.service.FindPaginatedResponseTimes(
		0,                // offset
		25,               // limit
		"timestamp desc", // order by
		0,                // reference ID
		"",               // date trunc
		nil,              // from
		nil,              // to
	)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 6, len(ResponseTimes))
		assert.Equal(suite.T(), testRecords[5].Timestamp.Unix(), ResponseTimes[0].Timestamp.Unix())
		assert.Equal(suite.T(), testRecords[4].Timestamp.Unix(), ResponseTimes[1].Timestamp.Unix())
		assert.Equal(suite.T(), testRecords[3].Timestamp.Unix(), ResponseTimes[2].Timestamp.Unix())
		assert.Equal(suite.T(), testRecords[2].Timestamp.Unix(), ResponseTimes[3].Timestamp.Unix())
		assert.Equal(suite.T(), testRecords[1].Timestamp.Unix(), ResponseTimes[4].Timestamp.Unix())
		assert.Equal(suite.T(), testRecords[0].Timestamp.Unix(), ResponseTimes[5].Timestamp.Unix())
	}

	// Test offset
	ResponseTimes, err = suite.service.FindPaginatedResponseTimes(
		3,   // offset
		25,  // limit
		"",  // order by
		0,   // reference ID
		"",  // date trunc
		nil, // from
		nil, // to
	)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 3, len(ResponseTimes))
		assert.Equal(suite.T(), testRecords[3].Timestamp.Unix(), ResponseTimes[0].Timestamp.Unix())
		assert.Equal(suite.T(), testRecords[4].Timestamp.Unix(), ResponseTimes[1].Timestamp.Unix())
		assert.Equal(suite.T(), testRecords[5].Timestamp.Unix(), ResponseTimes[2].Timestamp.Unix())
	}

	// Test limit
	ResponseTimes, err = suite.service.FindPaginatedResponseTimes(
		2,   // offset
		1,   // limit
		"",  // order by
		0,   // reference ID
		"",  // date trunc
		nil, // from
		nil, // to
	)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 1, len(ResponseTimes))
		assert.Equal(suite.T(), testRecords[2].Timestamp.Unix(), ResponseTimes[0].Timestamp.Unix())
	}
}
