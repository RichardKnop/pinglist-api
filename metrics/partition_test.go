package metrics

import (
	"time"

	"github.com/stretchr/testify/assert"
)

func (suite *MetricsTestSuite) TestPartitionRequestTime() {
	var (
		today                = time.Date(2016, time.February, 9, 0, 0, 0, 0, time.UTC)
		todaySubTableName    = "metrics_request_times_2016_02_09"
		tomorrowSubTableName = "metrics_request_times_2016_02_10"
		err                  error
	)

	// Sub tables should not exist yet
	assert.False(suite.T(), suite.service.db.HasTable(todaySubTableName))
	assert.False(suite.T(), suite.service.db.HasTable(tomorrowSubTableName))

	// This should create new sub tables for today and tomorrow
	err = suite.service.PartitionRequestTime(RequestTimeParentTableName, today)

	// Error should be nil and the today's sub table should have been created
	if assert.Nil(suite.T(), err) {
		assert.True(suite.T(), suite.service.db.HasTable(todaySubTableName))
		assert.True(suite.T(), suite.service.db.HasTable(tomorrowSubTableName))
	}
}

func (suite *MetricsTestSuite) TestCreateRequestTimeSubTable() {
	var (
		today                = time.Date(2016, time.February, 9, 0, 0, 0, 0, time.UTC)
		tomorrow             = today.Add(24 * time.Hour)
		dayAfterTomorrow     = tomorrow.Add(24 * time.Hour)
		todaySubTableName    = "metrics_request_times_2016_02_09"
		tomorrowSubTableName = "metrics_request_times_2016_02_10"
		subTable             *SubTable
		err                  error
		subTables            []*SubTable
		requestTimes         []*RequestTime
	)

	// Sub tables should not exist yet
	assert.False(suite.T(), suite.service.db.HasTable(todaySubTableName))
	assert.False(suite.T(), suite.service.db.HasTable(tomorrowSubTableName))

	// No sub table records should exists
	err = suite.db.Order("id").Find(&subTables).Error
	assert.NoError(suite.T(), err, "Fetching data failed")
	assert.Equal(suite.T(), 0, len(subTables))

	subTable, err = suite.service.createRequestTimeSubTable(
		RequestTimeParentTableName,
		todaySubTableName,
		today,
		tomorrow,
	)

	// Error should be nil and a today's sub table should have been created
	if assert.Nil(suite.T(), err) {
		assert.True(suite.T(), suite.service.db.HasTable(todaySubTableName))
		assert.False(suite.T(), suite.service.db.HasTable(tomorrowSubTableName))
	}

	// Correct sub table record should be returned
	if assert.NotNil(suite.T(), subTable) {
		assert.Equal(suite.T(), RequestTimeParentTableName, subTable.ParentTable)
		assert.Equal(suite.T(), todaySubTableName, subTable.Name)
	}

	// The today's sub table record should have been created
	err = suite.db.Order("id").Find(&subTables).Error
	assert.NoError(suite.T(), err, "Fetching data failed")
	assert.Equal(suite.T(), 1, len(subTables))
	assert.Equal(suite.T(), todaySubTableName, subTables[0].Name)

	subTable, err = suite.service.createRequestTimeSubTable(
		RequestTimeParentTableName,
		tomorrowSubTableName,
		tomorrow,
		dayAfterTomorrow,
	)

	// Error should be nil and a tomorrow's sub table should have been created
	if assert.Nil(suite.T(), err) {
		assert.True(suite.T(), suite.service.db.HasTable(todaySubTableName))
		assert.True(suite.T(), suite.service.db.HasTable(tomorrowSubTableName))
	}

	// Correct sub table record should be returned
	if assert.NotNil(suite.T(), subTable) {
		assert.Equal(suite.T(), RequestTimeParentTableName, subTable.ParentTable)
		assert.Equal(suite.T(), tomorrowSubTableName, subTable.Name)
	}

	// The tomorrow's sub table record should have been created
	err = suite.db.Order("id").Find(&subTables).Error
	assert.NoError(suite.T(), err, "Fetching data failed")
	assert.Equal(suite.T(), 2, len(subTables))
	assert.Equal(suite.T(), tomorrowSubTableName, subTables[1].Name)

	// These records should be inserted in the today's table
	requestTimes = []*RequestTime{
		newRequestTime(todaySubTableName, 1, today, 123),
		newRequestTime(todaySubTableName, 1, today.Add(1*time.Hour), 234),
	}
	for _, requestTime := range requestTimes {
		err := suite.db.Create(requestTime).Error
		assert.NoError(suite.T(), err, "Inserting test data failed")
	}

	// These records should be inserted in the tomorrow's table
	requestTimes = []*RequestTime{
		newRequestTime(tomorrowSubTableName, 1, tomorrow, 321),
		newRequestTime(tomorrowSubTableName, 1, tomorrow.Add(1*time.Hour), 432),
	}
	for _, requestTime := range requestTimes {
		err := suite.db.Create(requestTime).Error
		assert.NoError(suite.T(), err, "Inserting test data failed")
	}

	// Check data is being aggregated from all sub tables
	requestTimes = make([]*RequestTime, 0)
	err = suite.db.Order("timestamp").Find(&requestTimes).Error
	assert.NoError(suite.T(), err, "Fetching data failed")
	assert.Equal(suite.T(), 4, len(requestTimes))
	assert.Equal(suite.T(), int64(123), requestTimes[0].Value)
	assert.Equal(suite.T(), int64(234), requestTimes[1].Value)
	assert.Equal(suite.T(), int64(321), requestTimes[2].Value)
	assert.Equal(suite.T(), int64(432), requestTimes[3].Value)

	// Check data is correctly distributed to the today's table
	requestTimes = make([]*RequestTime, 0)
	err = suite.service.db.Table(todaySubTableName).
		Order("timestamp").Find(&requestTimes).Error
	assert.NoError(suite.T(), err, "Fetching data failed")
	assert.Equal(suite.T(), 2, len(requestTimes))
	assert.Equal(suite.T(), int64(123), requestTimes[0].Value)
	assert.Equal(suite.T(), int64(234), requestTimes[1].Value)

	// Check data is correctly distributed to the tomorrow's sub table
	requestTimes = make([]*RequestTime, 0)
	err = suite.service.db.Table(tomorrowSubTableName).
		Order("timestamp").Find(&requestTimes).Error
	assert.NoError(suite.T(), err, "Fetching data failed")
	assert.Equal(suite.T(), 2, len(requestTimes))
	assert.Equal(suite.T(), int64(321), requestTimes[0].Value)
	assert.Equal(suite.T(), int64(432), requestTimes[1].Value)
}
