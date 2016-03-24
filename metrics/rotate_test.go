package metrics

import (
	"time"

	"github.com/stretchr/testify/assert"
)

func (suite *MetricsTestSuite) TestRotateSubTables() {
	var (
		now      = time.Now()
		from, to time.Time
		subTable *SubTable
		err      error
		count    int
	)

	// Create test result sub tables and records
	for i := 1; i >= 0; i-- {
		from = now.Add(-time.Duration(i*RotateAfterHours) * time.Hour)
		to = from.Add(time.Duration(RotateAfterHours) * time.Hour)

		// Create a new sub table
		subTable, err = suite.service.createRequestTimeSubTable(
			RequestTimeParentTableName,
			getSubTableName(RequestTimeParentTableName, from),
			from,
			to,
		)
		assert.NoError(suite.T(), err, "Creating sub table failed")

		// Update the created_at field
		subTable.CreatedAt = from
		err := suite.service.db.Save(subTable).Error
		assert.NoError(suite.T(), err, "Updating created_at failed")
	}

	// 2 sub tables
	suite.service.db.Model(new(SubTable)).Count(&count)
	assert.Equal(suite.T(), 2, count)

	// Let's rotate the sub tables
	err = suite.service.RotateSubTables()

	// Error should be nil
	assert.Nil(suite.T(), err)

	// 1 sub table
	suite.service.db.Model(new(SubTable)).Count(&count)
	assert.Equal(suite.T(), 1, count)

	// Let's rotate the sub tables again
	err = suite.service.RotateSubTables()

	// Error should be nil
	assert.Nil(suite.T(), err)

	// 1 sub table
	suite.service.db.Model(new(SubTable)).Count(&count)
	assert.Equal(suite.T(), 1, count)
}
