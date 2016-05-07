package util

import (
	"database/sql"
	"database/sql/driver"
	"net/http"
	"testing"
	"time"

	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestParseTimestamp(t *testing.T) {
	var (
		parsedTimestamp time.Time
		err             error
	)

	parsedTimestamp, err = ParseTimestamp("bogus")
	assert.NotNil(t, err)

	parsedTimestamp, err = ParseTimestamp("2016-05-04T12:08:35Z")
	assert.Nil(t, err)
	assert.Equal(t, 2016, parsedTimestamp.UTC().Year())
	assert.Equal(t, time.May, parsedTimestamp.UTC().Month())
	assert.Equal(t, 4, parsedTimestamp.UTC().Day())
	assert.Equal(t, 12, parsedTimestamp.UTC().Hour())
	assert.Equal(t, 8, parsedTimestamp.UTC().Minute())
	assert.Equal(t, 35, parsedTimestamp.UTC().Second())

	parsedTimestamp, err = ParseTimestamp("2016-05-04T12:08:35+07:00")
	assert.Nil(t, err)
	assert.Equal(t, 2016, parsedTimestamp.UTC().Year())
	assert.Equal(t, time.May, parsedTimestamp.UTC().Month())
	assert.Equal(t, 4, parsedTimestamp.UTC().Day())
	assert.Equal(t, 5, parsedTimestamp.UTC().Hour())
	assert.Equal(t, 8, parsedTimestamp.UTC().Minute())
	assert.Equal(t, 35, parsedTimestamp.UTC().Second())
}

func TestIntOrNull(t *testing.T) {
	nullInt := PositiveIntOrNull(1)
	assert.True(t, nullInt.Valid)

	value, err := nullInt.Value()
	assert.Nil(t, err)
	assert.Equal(t, int64(1), value)
}

func TestPositiveIntOrNull(t *testing.T) {
	var (
		nullInt sql.NullInt64
		value   driver.Value
		err     error
	)

	// When the number is negative
	nullInt = PositiveIntOrNull(-1)

	// nullInt.Valid should be false
	assert.False(t, nullInt.Valid)

	// nullInt.Value() should return nil
	value, err = nullInt.Value()
	assert.Nil(t, err)
	assert.Nil(t, value)

	// When the number is greater than zero
	nullInt = PositiveIntOrNull(1)

	// nullInt.Valid should be true
	assert.True(t, nullInt.Valid)

	// nullInt.Value() should return the integer
	value, err = nullInt.Value()
	assert.Nil(t, err)
	assert.Equal(t, int64(1), value)
}

func TestFloatOrNull(t *testing.T) {
	nullFloat := FloatOrNull(1.5)
	assert.True(t, nullFloat.Valid)

	value, err := nullFloat.Value()
	assert.Nil(t, err)
	assert.Equal(t, 1.5, value)
}

func TestPositiveFloatOrNull(t *testing.T) {
	var (
		nullFloat sql.NullFloat64
		value     driver.Value
		err       error
	)

	// When the number is negative
	nullFloat = PositiveFloatOrNull(-0.5)

	// nullFloat.Valid should be false
	assert.False(t, nullFloat.Valid)

	// nullFloat.Value() should return nil
	value, err = nullFloat.Value()
	assert.Nil(t, err)
	assert.Nil(t, value)

	// When the number is greater than zero
	nullFloat = PositiveFloatOrNull(1.5)

	// nullFloat.Valid should be true
	assert.True(t, nullFloat.Valid)

	// nullFloat.Value() should return the integer
	value, err = nullFloat.Value()
	assert.Nil(t, err)
	assert.Equal(t, 1.5, value)
}

func TestStringOrNull(t *testing.T) {
	var (
		nullString sql.NullString
		value      driver.Value
		err        error
	)

	// When the string is empty
	nullString = StringOrNull("")

	// nullString.Valid should be false
	assert.False(t, nullString.Valid)

	// nullString.Value() should return nil
	value, err = nullString.Value()
	assert.Nil(t, err)
	assert.Nil(t, value)

	// When the string is not empty
	nullString = StringOrNull("foo")

	// nullString.Valid should be true
	assert.True(t, nullString.Valid)

	// nullString.Value() should return the string
	value, err = nullString.Value()
	assert.Nil(t, err)
	assert.Equal(t, "foo", value)
}

func TestTimeOrNull(t *testing.T) {
	var (
		nullTime pq.NullTime
		value    driver.Value
		err      error
	)

	// When the time is nil
	nullTime = TimeOrNull(nil)

	// nullTime.Valid should be false
	assert.False(t, nullTime.Valid)

	// nullInt.Value() should return nil
	value, err = nullTime.Value()
	assert.Nil(t, err)
	assert.Nil(t, value)

	// When the time is time.Time instance
	now := time.Now()
	nullTime = TimeOrNull(&now)

	// nullTime.Valid should be true
	assert.True(t, nullTime.Valid)

	// nullTime.Value() should return the time.Time
	value, err = nullTime.Value()
	assert.Nil(t, err)
	assert.Equal(t, now, value)
}

func TestStringInSlice(t *testing.T) {
	assert.True(t, StringInSlice("a", []string{"a", "b", "c"}))

	assert.False(t, StringInSlice("d", []string{"a", "b", "c"}))
}

func TestSpaceDelimitedStringNotGreater(t *testing.T) {
	assert.True(t, SpaceDelimitedStringNotGreater("", "bar foo qux"))

	assert.True(t, SpaceDelimitedStringNotGreater("foo", "bar foo qux"))

	assert.True(t, SpaceDelimitedStringNotGreater("bar foo qux", "foo bar qux"))

	assert.False(t, SpaceDelimitedStringNotGreater("foo bar qux bogus", "bar foo qux"))
}

func TestParseBearerTokenNotFound(t *testing.T) {
	r, err := http.NewRequest("GET", "http://1.2.3.4/something", nil)
	assert.NoError(t, err, "Request setup should not get an error")
	r.Header.Add("Authorization", "bogus bogus")

	token, err := ParseBearerToken(r)

	// Token should be nil
	assert.Nil(t, token)

	// Correct error should be returned
	if assert.NotNil(t, err) {
		assert.Equal(t, "Bearer token not found", err.Error())
	}
}

func TestParseBearerToken(t *testing.T) {
	r, err := http.NewRequest("GET", "http://1.2.3.4/something", nil)
	assert.NoError(t, err, "Request setup should not get an error")
	r.Header.Add("Authorization", "Bearer test_token")

	token, err := ParseBearerToken(r)

	// Error should be nil
	assert.Nil(t, err)

	// Correct token should be returned
	if assert.NotNil(t, token) {
		assert.Equal(t, []byte("test_token"), token)
	}
}
