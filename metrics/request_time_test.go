package metrics

import (
	"time"

	"github.com/stretchr/testify/assert"
)

func (suite *MetricsTestSuite) TestPaginatedRequestTimesCount() {
	var (
		today                 = time.Date(2016, time.February, 9, 0, 0, 0, 0, time.UTC)
		todaySubTableName     = "metrics_request_times_2016_02_09"
		yesterday             = time.Date(2016, time.February, 8, 0, 0, 0, 0, time.UTC)
		yesterdaySubTableName = "metrics_request_times_2016_02_08"
		from, to              time.Time
		count                 int
		err                   error
	)

	// Partition the request time table
	err = suite.service.PartitionRequestTime(RequestTimeParentTableName, today)
	assert.NoError(suite.T(), err, "Partitioning table failed")
	err = suite.service.PartitionRequestTime(RequestTimeParentTableName, yesterday)
	assert.NoError(suite.T(), err, "Partitioning table failed")

	// Insert some test records
	testRecords := []*RequestTime{
		NewRequestTime(yesterdaySubTableName, 1, yesterday, 123),
		NewRequestTime(yesterdaySubTableName, 1, yesterday.Add(1*time.Hour), 234),
		NewRequestTime(yesterdaySubTableName, 2, yesterday.Add(2*time.Hour), 345),
		NewRequestTime(todaySubTableName, 1, today, 123),
		NewRequestTime(todaySubTableName, 1, today.Add(1*time.Hour), 234),
		NewRequestTime(todaySubTableName, 2, today.Add(2*time.Hour), 345),
	}
	for _, testRecord := range testRecords {
		err := suite.db.Create(testRecord).Error
		assert.NoError(suite.T(), err, "Inserting test data failed")
	}

	// No filtering at all
	count, err = suite.service.PaginatedRequestTimesCount(
		0,   // reference ID
		"",  // date trunc
		nil, // from
		nil, // to
	)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 6, count)
	}

	// Filter by a valid reference ID
	count, err = suite.service.PaginatedRequestTimesCount(
		1,   // reference ID
		"",  // date trunc
		nil, // from
		nil, // to
	)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 4, count)
	}

	// Filter by a valid reference ID
	count, err = suite.service.PaginatedRequestTimesCount(
		2,   // reference ID
		"",  // date trunc
		nil, // from
		nil, // to
	)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 2, count)
	}

	// Filter by a bogus reference ID
	count, err = suite.service.PaginatedRequestTimesCount(
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
	count, err = suite.service.PaginatedRequestTimesCount(
		0,     // reference ID
		"",    // date trunc
		&from, // from
		&to,   // to
	)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 2, count)
	}

	// Filter by "date_trunc" timestamp
	count, err = suite.service.PaginatedRequestTimesCount(
		0,     // reference ID
		"day", // date trunc
		nil,   // from
		nil,   // to
	)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 2, count)
	}
}

func (suite *MetricsTestSuite) TestFindPaginatedRequestTimes() {
	var (
		today                 = time.Date(2016, time.February, 9, 0, 0, 0, 0, time.UTC)
		todaySubTableName     = "metrics_request_times_2016_02_09"
		yesterday             = time.Date(2016, time.February, 8, 0, 0, 0, 0, time.UTC)
		yesterdaySubTableName = "metrics_request_times_2016_02_08"
		from, to              time.Time
		requestTimes          []*RequestTime
		err                   error
	)

	// Partition the request time table
	err = suite.service.PartitionRequestTime(RequestTimeParentTableName, today)
	assert.NoError(suite.T(), err, "Partitioning table failed")
	err = suite.service.PartitionRequestTime(RequestTimeParentTableName, yesterday)
	assert.NoError(suite.T(), err, "Partitioning table failed")

	// Insert some test records
	testRecords := []*RequestTime{
		NewRequestTime(yesterdaySubTableName, 1, yesterday, 123),
		NewRequestTime(yesterdaySubTableName, 1, yesterday.Add(1*time.Hour), 234),
		NewRequestTime(yesterdaySubTableName, 2, yesterday.Add(2*time.Hour), 345),
		NewRequestTime(todaySubTableName, 1, today, 456),
		NewRequestTime(todaySubTableName, 1, today.Add(1*time.Hour), 567),
		NewRequestTime(todaySubTableName, 2, today.Add(2*time.Hour), 678),
	}
	for _, testRecord := range testRecords {
		err := suite.db.Create(testRecord).Error
		assert.NoError(suite.T(), err, "Inserting test data failed")
	}

	// No filtering at all
	requestTimes, err = suite.service.FindPaginatedRequestTimes(
		0,   // offset
		25,  // limit
		"",  // order by
		0,   // reference ID
		"",  // date trunc
		nil, // from
		nil, // to
	)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 6, len(requestTimes))
	}

	// Filter by a valid reference ID
	requestTimes, err = suite.service.FindPaginatedRequestTimes(
		0,   // offset
		25,  // limit
		"",  // order by
		1,   // reference ID
		"",  // date trunc
		nil, // from
		nil, // to
	)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 4, len(requestTimes))
		assert.Equal(suite.T(), testRecords[0].Timestamp.Unix(), requestTimes[0].Timestamp.Unix())
		assert.Equal(suite.T(), testRecords[1].Timestamp.Unix(), requestTimes[1].Timestamp.Unix())
		assert.Equal(suite.T(), testRecords[3].Timestamp.Unix(), requestTimes[2].Timestamp.Unix())
		assert.Equal(suite.T(), testRecords[4].Timestamp.Unix(), requestTimes[3].Timestamp.Unix())
	}

	// Filter by a valid reference ID
	requestTimes, err = suite.service.FindPaginatedRequestTimes(
		0,   // offset
		25,  // limit
		"",  // order by
		2,   // reference ID
		"",  // date trunc
		nil, // from
		nil, // to
	)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 2, len(requestTimes))
		assert.Equal(suite.T(), testRecords[2].Timestamp.Unix(), requestTimes[0].Timestamp.Unix())
		assert.Equal(suite.T(), testRecords[5].Timestamp.Unix(), requestTimes[1].Timestamp.Unix())
	}

	// Filter by a bogus reference ID
	requestTimes, err = suite.service.FindPaginatedRequestTimes(
		0,   // offset
		25,  // limit
		"",  // order by
		3,   // reference ID
		"",  // date trunc
		nil, // from
		nil, // to
	)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 0, len(requestTimes))
	}

	// Filter by "from" and "to" timestamp
	from = yesterday.Add(2 * time.Hour)
	to = today
	requestTimes, err = suite.service.FindPaginatedRequestTimes(
		0,     // offset
		25,    // limit
		"",    // order by
		0,     // reference ID
		"",    // date trunc
		&from, // from
		&to,   // to
	)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 2, len(requestTimes))
		assert.Equal(suite.T(), 2, len(requestTimes))
		assert.Equal(suite.T(), testRecords[2].Timestamp.Unix(), requestTimes[0].Timestamp.Unix())
		assert.Equal(suite.T(), testRecords[3].Timestamp.Unix(), requestTimes[1].Timestamp.Unix())
	}

	// Filter by "date_trunc" timestamp
	requestTimes, err = suite.service.FindPaginatedRequestTimes(
		0,     // offset
		25,    // limit
		"",    // order by
		0,     // reference ID
		"day", // date trunc
		nil,   // from
		nil,   // to
	)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 2, len(requestTimes))
		assert.Equal(
			suite.T(),
			yesterday.UTC().Format(time.RFC3339),
			requestTimes[0].Timestamp.UTC().Format(time.RFC3339),
		)
		assert.Equal(
			suite.T(),
			today.UTC().Format(time.RFC3339),
			requestTimes[1].Timestamp.UTC().Format(time.RFC3339),
		)
		assert.Equal(suite.T(), int64(234), requestTimes[0].Value) // (123 + 234 + 345) / 3
		assert.Equal(suite.T(), int64(567), requestTimes[1].Value) // (456 + 567 + 678) / 3
	}

	// This should return all records
	requestTimes, err = suite.service.FindPaginatedRequestTimes(
		0,   // offset
		25,  // limit
		"",  // order by
		0,   // reference ID
		"",  // date trunc
		nil, // from
		nil, // to
	)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 6, len(requestTimes))
		assert.Equal(suite.T(), testRecords[0].Timestamp.Unix(), requestTimes[0].Timestamp.Unix())
		assert.Equal(suite.T(), testRecords[1].Timestamp.Unix(), requestTimes[1].Timestamp.Unix())
		assert.Equal(suite.T(), testRecords[2].Timestamp.Unix(), requestTimes[2].Timestamp.Unix())
		assert.Equal(suite.T(), testRecords[3].Timestamp.Unix(), requestTimes[3].Timestamp.Unix())
		assert.Equal(suite.T(), testRecords[4].Timestamp.Unix(), requestTimes[4].Timestamp.Unix())
		assert.Equal(suite.T(), testRecords[5].Timestamp.Unix(), requestTimes[5].Timestamp.Unix())
	}

	// This should return all records ordered by timestamp desc
	requestTimes, err = suite.service.FindPaginatedRequestTimes(
		0,                // offset
		25,               // limit
		"timestamp desc", // order by
		0,                // reference ID
		"",               // date trunc
		nil,              // from
		nil,              // to
	)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 6, len(requestTimes))
		assert.Equal(suite.T(), testRecords[5].Timestamp.Unix(), requestTimes[0].Timestamp.Unix())
		assert.Equal(suite.T(), testRecords[4].Timestamp.Unix(), requestTimes[1].Timestamp.Unix())
		assert.Equal(suite.T(), testRecords[3].Timestamp.Unix(), requestTimes[2].Timestamp.Unix())
		assert.Equal(suite.T(), testRecords[2].Timestamp.Unix(), requestTimes[3].Timestamp.Unix())
		assert.Equal(suite.T(), testRecords[1].Timestamp.Unix(), requestTimes[4].Timestamp.Unix())
		assert.Equal(suite.T(), testRecords[0].Timestamp.Unix(), requestTimes[5].Timestamp.Unix())
	}

	// Test offset
	requestTimes, err = suite.service.FindPaginatedRequestTimes(
		3,   // offset
		25,  // limit
		"",  // order by
		0,   // reference ID
		"",  // date trunc
		nil, // from
		nil, // to
	)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 3, len(requestTimes))
		assert.Equal(suite.T(), testRecords[3].Timestamp.Unix(), requestTimes[0].Timestamp.Unix())
		assert.Equal(suite.T(), testRecords[4].Timestamp.Unix(), requestTimes[1].Timestamp.Unix())
		assert.Equal(suite.T(), testRecords[5].Timestamp.Unix(), requestTimes[2].Timestamp.Unix())
	}

	// Test limit
	requestTimes, err = suite.service.FindPaginatedRequestTimes(
		2,   // offset
		1,   // limit
		"",  // order by
		0,   // reference ID
		"",  // date trunc
		nil, // from
		nil, // to
	)
	if assert.Nil(suite.T(), err) {
		assert.Equal(suite.T(), 1, len(requestTimes))
		assert.Equal(suite.T(), testRecords[2].Timestamp.Unix(), requestTimes[0].Timestamp.Unix())
	}
}
