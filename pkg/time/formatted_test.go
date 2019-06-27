package time

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func ExampleFormatted() {
	someTime, err := time.Parse("2006-01-02T15:04:05Z07:00", "2018-12-24T08:41:53Z")
	if err != nil {
		fmt.Println(err.Error())
	}
	customTime := Formatted(someTime)

	fmt.Println(customTime)
	// Output:
	// 2018-12-24T08:41:53Z
}

func ExampleFormatted_withTimeZone() {
	someTime, err := time.Parse("2006-01-02T15:04:05Z07:00", "2018-12-24T08:41:53+07:00")
	if err != nil {
		fmt.Println(err.Error())
	}
	customTime := Formatted(someTime)

	fmt.Println(customTime)
	// Output:
	// 2018-12-24T08:41:53+07:00
}

func TestNewUTCTime(t *testing.T) {
	utcTime := time.Time(NewUTCTime())

	assert.Equal(t, time.UTC, utcTime.Location())
}

func TestFormattedTimeMarshalJSON(t *testing.T) {
	someTime, err := time.Parse("2006-01-02T15:04:05Z07:00", "2018-12-24T08:41:53Z")
	require.NoError(t, err)

	data, err := Formatted(someTime).MarshalJSON()

	require.NoError(t, err)
	assert.Equal(t, []byte("\"2018-12-24T08:41:53Z\""), data)
}

func TestFormattedTimeUnmarshalJSON(t *testing.T) {
	formatted := &Formatted{}
	err := formatted.UnmarshalJSON([]byte("\"2018-12-24T08:41:53Z\""))
	require.NoError(t, err)

	assert.Equal(t, "2018-12-24T08:41:53Z", formatted.String())
}

func TestFormattedTimeString(t *testing.T) {
	someTime, err := time.Parse("2006-01-02T15:04:05Z07:00", "2018-12-24T08:41:53Z")
	require.NoError(t, err)

	assert.Equal(t, "2018-12-24T08:41:53Z", Formatted(someTime).String())
}

func TestFormattedTimeJsonExample(t *testing.T) {
	justNow := time.Now()
	value := &struct{ Infinity Formatted }{}
	jsonPayload := []byte("{\"infinity\": \"" + justNow.Format(time.RFC3339) + "\"}")

	err := json.Unmarshal(jsonPayload, value)

	require.NoError(t, err)
	assert.Equal(t, Formatted(justNow).String(), value.Infinity.String())
}
