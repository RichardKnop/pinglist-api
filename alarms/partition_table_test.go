package alarms

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func (suite *AlarmsTestSuite) TestPartitionTable() {
	var (
		today                = time.Date(2016, time.February, 9, 0, 0, 0, 0, time.UTC)
		tomorrow             = today.Add(24 * time.Hour)
		todaySubTableName    = "alarm_results_2016_02_09"
		tomorrowSubTableName = "alarm_results_2016_02_10"
		err                  error
	)

	// Sub tables should not exist yet
	assert.False(suite.T(), suite.service.db.HasTable(todaySubTableName))
	assert.False(suite.T(), suite.service.db.HasTable(tomorrowSubTableName))

	// This should create a today's sub table
	err = suite.service.PartitionTable(ResultParentTableName, today)

	// Error should be nil and the today's sub table should have been created
	if assert.Nil(suite.T(), err) {
		assert.True(suite.T(), suite.service.db.HasTable(todaySubTableName))
		assert.False(suite.T(), suite.service.db.HasTable(tomorrowSubTableName))
	}

	// Since we are nearing the next day, a tomorrow's sub table should be created
	err = suite.service.PartitionTable(
		ResultParentTableName,
		tomorrow.Add(-59*time.Minute),
	)

	// Error should be nil and the tomorrow's sub table should have been created
	if assert.Nil(suite.T(), err) {
		assert.True(suite.T(), suite.service.db.HasTable(todaySubTableName))
		assert.True(suite.T(), suite.service.db.HasTable(tomorrowSubTableName))
	}
}

func (suite *AlarmsTestSuite) TestRotateSubTables() {
	var (
		now            = time.Now()
		from, to       time.Time
		resultSubTable *ResultSubTable
		err            error
		count          int
	)

	// Create test result sub tables and records
	for i := 1; i >= 0; i-- {
		from = now.Add(-time.Duration(i*rotateAfterHours) * time.Hour)
		to = from.Add(time.Duration(rotateAfterHours) * time.Hour)

		// Create a new sub table
		resultSubTable, err = suite.service.createSubTable(
			ResultParentTableName,
			getSubtableName(ResultParentTableName, from),
			from,
			to,
		)
		assert.NoError(suite.T(), err, "Creating sub table failed")

		// Update the created_at field
		resultSubTable.CreatedAt = from
		err := suite.service.db.Save(resultSubTable).Error
		assert.NoError(suite.T(), err, "Updating created_at failed")
	}

	// 2 sub tables
	suite.service.db.Model(new(ResultSubTable)).Count(&count)
	assert.Equal(suite.T(), 2, count)

	// Let's rotate the sub tables
	err = suite.service.RotateSubTables()

	// Error should be nil
	assert.Nil(suite.T(), err)

	// 1 sub table
	suite.service.db.Model(new(ResultSubTable)).Count(&count)
	assert.Equal(suite.T(), 1, count)

	// Let's rotate the sub tables again
	err = suite.service.RotateSubTables()

	// Error should be nil
	assert.Nil(suite.T(), err)

	// 1 sub tables
	suite.service.db.Model(new(ResultSubTable)).Count(&count)
	assert.Equal(suite.T(), 1, count)
}

func TestGetSubTableName(t *testing.T) {
	var (
		expected string
		actual   string
	)

	expected = "parent_table_2016_12_27"
	actual = getSubtableName(
		"parent_table",
		time.Date(2016, time.December, 27, 15, 0, 0, 0, time.UTC),
	)
	assert.Equal(t, expected, actual)

	expected = "parent_table_2016_01_02"
	actual = getSubtableName(
		"parent_table",
		time.Date(2016, time.January, 2, 15, 0, 0, 0, time.UTC),
	)
	assert.Equal(t, expected, actual)

	expected = "parent_table_2016_02_01"
	actual = getSubtableName(
		"parent_table",
		time.Date(2016, time.February, 1, 15, 0, 0, 0, time.UTC),
	)
	assert.Equal(t, expected, actual)
}

func (suite *AlarmsTestSuite) TestCreateSubTable() {
	var (
		today                = time.Date(2016, time.February, 9, 0, 0, 0, 0, time.UTC)
		tomorrow             = today.Add(24 * time.Hour)
		dayAfterTomorrow     = tomorrow.Add(24 * time.Hour)
		todaySubTableName    = "alarm_results_2016_02_09"
		tomorrowSubTableName = "alarm_results_2016_02_10"
		resultSubTable       *ResultSubTable
		err                  error
		resultSubTables      []*ResultSubTable
		results              []*Result
	)

	// Sub tables should not exist yet
	assert.False(suite.T(), suite.service.db.HasTable(todaySubTableName))
	assert.False(suite.T(), suite.service.db.HasTable(tomorrowSubTableName))

	// No sub table records should exists
	err = suite.db.Find(&resultSubTables).Error
	assert.NoError(suite.T(), err, "Fetching data failed")
	assert.Equal(suite.T(), 0, len(resultSubTables))

	resultSubTable, err = suite.service.createSubTable(
		ResultParentTableName,
		todaySubTableName,
		today,
		tomorrow,
	)

	// Error should be nil and a today's sub table should have been created
	if assert.Nil(suite.T(), err) {
		assert.True(suite.T(), suite.service.db.HasTable(todaySubTableName))
		assert.False(suite.T(), suite.service.db.HasTable(tomorrowSubTableName))
	}

	// Correct result sub table object should be returned
	if assert.NotNil(suite.T(), resultSubTable) {
		assert.Equal(suite.T(), todaySubTableName, resultSubTable.Name)
	}

	// The today's sub table record should have been created
	err = suite.db.Find(&resultSubTables).Error
	assert.NoError(suite.T(), err, "Fetching data failed")
	assert.Equal(suite.T(), 1, len(resultSubTables))
	assert.Equal(suite.T(), todaySubTableName, resultSubTables[0].Name)

	resultSubTable, err = suite.service.createSubTable(
		ResultParentTableName,
		tomorrowSubTableName,
		tomorrow,
		dayAfterTomorrow,
	)

	// Error should be nil and a tomorrow's sub table should have been created
	if assert.Nil(suite.T(), err) {
		assert.True(suite.T(), suite.service.db.HasTable(todaySubTableName))
		assert.True(suite.T(), suite.service.db.HasTable(tomorrowSubTableName))
	}

	// Correct result sub table object should be returned
	if assert.NotNil(suite.T(), resultSubTable) {
		assert.Equal(suite.T(), tomorrowSubTableName, resultSubTable.Name)
	}

	// The tomorrow's sub table record should have been created
	err = suite.db.Order("id").Find(&resultSubTables).Error
	assert.NoError(suite.T(), err, "Fetching data failed")
	assert.Equal(suite.T(), 2, len(resultSubTables))
	assert.Equal(suite.T(), tomorrowSubTableName, resultSubTables[1].Name)

	// These results should be inserted in the today's table
	results = []*Result{
		newResult(todaySubTableName, suite.alarms[0], today, 123),
		newResult(todaySubTableName, suite.alarms[0], today.Add(1*time.Hour), 234),
	}
	for _, result := range results {
		err := suite.db.Create(result).Error
		assert.NoError(suite.T(), err, "Inserting test data failed")
	}

	// These results should be inserted in the tomorrow's table
	results = []*Result{
		newResult(tomorrowSubTableName, suite.alarms[0], tomorrow, 321),
		newResult(tomorrowSubTableName, suite.alarms[0], tomorrow.Add(1*time.Hour), 432),
	}
	for _, result := range results {
		err := suite.db.Create(result).Error
		assert.NoError(suite.T(), err, "Inserting test data failed")
	}

	// Check data is being aggregated from all sub tables
	results = make([]*Result, 0)
	err = suite.db.Preload("Alarm").Find(&results).Error
	assert.NoError(suite.T(), err, "Fetching data failed")
	assert.Equal(suite.T(), 4, len(results))
	assert.Equal(suite.T(), int64(123), results[0].RequestTime)
	assert.Equal(suite.T(), int64(234), results[1].RequestTime)
	assert.Equal(suite.T(), int64(321), results[2].RequestTime)
	assert.Equal(suite.T(), int64(432), results[3].RequestTime)

	// Check data is correctly distributed to the today's table
	results = make([]*Result, 0)
	err = suite.service.db.Table(todaySubTableName).Find(&results).Error
	assert.NoError(suite.T(), err, "Fetching data failed")
	assert.Equal(suite.T(), 2, len(results))
	assert.Equal(suite.T(), int64(123), results[0].RequestTime)
	assert.Equal(suite.T(), int64(234), results[1].RequestTime)

	// Check data is correctly distributed to the tomorrow's sub table
	results = make([]*Result, 0)
	err = suite.service.db.Table(tomorrowSubTableName).Find(&results).Error
	assert.NoError(suite.T(), err, "Fetching data failed")
	assert.Equal(suite.T(), 2, len(results))
	assert.Equal(suite.T(), int64(321), results[0].RequestTime)
	assert.Equal(suite.T(), int64(432), results[1].RequestTime)
}
