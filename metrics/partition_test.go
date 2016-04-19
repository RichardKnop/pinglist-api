package metrics

import (
	"time"

	"github.com/stretchr/testify/assert"
)

func (suite *MetricsTestSuite) TestPartitionResponseTime() {
	var (
		today                = time.Date(2016, time.February, 9, 0, 0, 0, 0, time.UTC)
		todaySubTableName    = "metrics_response_times_2016_02_09"
		tomorrowSubTableName = "metrics_response_times_2016_02_10"
		err                  error
	)

	// Sub tables should not exist yet
	assert.False(suite.T(), suite.service.db.HasTable(todaySubTableName))
	assert.False(suite.T(), suite.service.db.HasTable(tomorrowSubTableName))

	// This should create new sub tables for today and tomorrow
	err = suite.service.PartitionResponseTime(ResponseTimeParentTableName, today)

	// Error should be nil and the today's sub table should have been created
	if assert.Nil(suite.T(), err) {
		assert.True(suite.T(), suite.service.db.HasTable(todaySubTableName))
		assert.True(suite.T(), suite.service.db.HasTable(tomorrowSubTableName))
	}
}

func (suite *MetricsTestSuite) TestCreateResponseTimeSubTable() {
	var (
		today                = time.Date(2016, time.February, 9, 0, 0, 0, 0, time.UTC)
		tomorrow             = today.Add(24 * time.Hour)
		dayAfterTomorrow     = tomorrow.Add(24 * time.Hour)
		todaySubTableName    = "metrics_response_times_2016_02_09"
		tomorrowSubTableName = "metrics_response_times_2016_02_10"
		subTable             *SubTable
		err                  error
		subTables            []*SubTable
		ResponseTimes        []*ResponseTime
	)

	// Sub tables should not exist yet
	assert.False(suite.T(), suite.service.db.HasTable(todaySubTableName))
	assert.False(suite.T(), suite.service.db.HasTable(tomorrowSubTableName))

	// No sub table records should exists
	err = suite.db.Order("id").Find(&subTables).Error
	assert.NoError(suite.T(), err, "Fetching data failed")
	assert.Equal(suite.T(), 0, len(subTables))

	subTable, err = suite.service.createResponseTimeSubTable(
		ResponseTimeParentTableName,
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
		assert.Equal(suite.T(), ResponseTimeParentTableName, subTable.ParentTable)
		assert.Equal(suite.T(), todaySubTableName, subTable.Name)
	}

	// The today's sub table record should have been created
	err = suite.db.Order("id").Find(&subTables).Error
	assert.NoError(suite.T(), err, "Fetching data failed")
	assert.Equal(suite.T(), 1, len(subTables))
	assert.Equal(suite.T(), todaySubTableName, subTables[0].Name)

	subTable, err = suite.service.createResponseTimeSubTable(
		ResponseTimeParentTableName,
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
		assert.Equal(suite.T(), ResponseTimeParentTableName, subTable.ParentTable)
		assert.Equal(suite.T(), tomorrowSubTableName, subTable.Name)
	}

	// The tomorrow's sub table record should have been created
	err = suite.db.Order("id").Find(&subTables).Error
	assert.NoError(suite.T(), err, "Fetching data failed")
	assert.Equal(suite.T(), 2, len(subTables))
	assert.Equal(suite.T(), tomorrowSubTableName, subTables[1].Name)

	// These records should be inserted in the today's table
	ResponseTimes = []*ResponseTime{
		NewResponseTime(todaySubTableName, 1, today, 123),
		NewResponseTime(todaySubTableName, 1, today.Add(1*time.Hour), 234),
	}
	for _, ResponseTime := range ResponseTimes {
		err := suite.db.Create(ResponseTime).Error
		assert.NoError(suite.T(), err, "Inserting test data failed")
	}

	// These records should be inserted in the tomorrow's table
	ResponseTimes = []*ResponseTime{
		NewResponseTime(tomorrowSubTableName, 1, tomorrow, 321),
		NewResponseTime(tomorrowSubTableName, 1, tomorrow.Add(1*time.Hour), 432),
	}
	for _, ResponseTime := range ResponseTimes {
		err := suite.db.Create(ResponseTime).Error
		assert.NoError(suite.T(), err, "Inserting test data failed")
	}

	// Check data is being aggregated from all sub tables
	ResponseTimes = make([]*ResponseTime, 0)
	err = suite.db.Order("timestamp").Find(&ResponseTimes).Error
	assert.NoError(suite.T(), err, "Fetching data failed")
	assert.Equal(suite.T(), 4, len(ResponseTimes))
	assert.Equal(suite.T(), int64(123), ResponseTimes[0].Value)
	assert.Equal(suite.T(), int64(234), ResponseTimes[1].Value)
	assert.Equal(suite.T(), int64(321), ResponseTimes[2].Value)
	assert.Equal(suite.T(), int64(432), ResponseTimes[3].Value)

	// Check data is correctly distributed to the today's table
	ResponseTimes = make([]*ResponseTime, 0)
	err = suite.service.db.Table(todaySubTableName).
		Order("timestamp").Find(&ResponseTimes).Error
	assert.NoError(suite.T(), err, "Fetching data failed")
	assert.Equal(suite.T(), 2, len(ResponseTimes))
	assert.Equal(suite.T(), int64(123), ResponseTimes[0].Value)
	assert.Equal(suite.T(), int64(234), ResponseTimes[1].Value)

	// Check data is correctly distributed to the tomorrow's sub table
	ResponseTimes = make([]*ResponseTime, 0)
	err = suite.service.db.Table(tomorrowSubTableName).
		Order("timestamp").Find(&ResponseTimes).Error
	assert.NoError(suite.T(), err, "Fetching data failed")
	assert.Equal(suite.T(), 2, len(ResponseTimes))
	assert.Equal(suite.T(), int64(321), ResponseTimes[0].Value)
	assert.Equal(suite.T(), int64(432), ResponseTimes[1].Value)
}
