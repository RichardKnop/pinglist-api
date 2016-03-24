package timeseries

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetSubTableName(t *testing.T) {
	var (
		expected string
		actual   string
	)

	expected = "parent_table_2016_12_27"
	actual = getSubTableName(
		"parent_table",
		time.Date(2016, time.December, 27, 15, 0, 0, 0, time.UTC),
	)
	assert.Equal(t, expected, actual)

	expected = "parent_table_2016_01_02"
	actual = getSubTableName(
		"parent_table",
		time.Date(2016, time.January, 2, 15, 0, 0, 0, time.UTC),
	)
	assert.Equal(t, expected, actual)

	expected = "parent_table_2016_02_01"
	actual = getSubTableName(
		"parent_table",
		time.Date(2016, time.February, 1, 15, 0, 0, 0, time.UTC),
	)
	assert.Equal(t, expected, actual)
}
