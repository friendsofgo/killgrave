package templating

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJsonMarshal(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name:     "simple map",
			input:    map[string]string{"key": "value"},
			expected: `{"key":"value"}`,
		},
		{
			name:     "nested map",
			input:    map[string]interface{}{"key": map[string]string{"nestedKey": "nestedValue"}},
			expected: `{"key":{"nestedKey":"nestedValue"}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := JsonMarshal(tt.input)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestJsonMarshal_Invalid(t *testing.T) {
	input := make(chan int)
	_, err := JsonMarshal(input)
	assert.Error(t, err)
}

func TestTimeNow(t *testing.T) {
	before := time.Now().Truncate(time.Second)
	got := TimeNow()
	after := time.Now().Truncate(time.Second)

	parsedTime, err := time.Parse(time.RFC3339, got)
	assert.NoError(t, err, "TimeNow() returned invalid time format")

	assert.True(t, parsedTime.After(before) || parsedTime.Equal(before))
	assert.True(t, parsedTime.Before(after) || parsedTime.Equal(after))
}

func TestTimeUTC(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "valid RFC3339 time",
			input:    "2023-10-15T13:34:02Z",
			expected: "2023-10-15T13:34:02Z",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TimeUTC(tt.input)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestTimeUTC_Invalid(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		errMsg string
	}{
		{
			name:   "invalid time format",
			input:  "invalid-time",
			errMsg: "parsing time",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := TimeUTC(tt.input)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.errMsg)
		})
	}
}

func TestTimeAdd(t *testing.T) {
	tests := []struct {
		name     string
		time     string
		duration string
		expected string
	}{
		{
			name:     "add 1 hour",
			time:     "2023-10-15T13:34:02Z",
			duration: "1h",
			expected: "2023-10-15T14:34:02Z",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TimeAdd(tt.time, tt.duration)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestTimeAdd_Invalid(t *testing.T) {
	tests := []struct {
		name     string
		time     string
		duration string
		errMsg   string
	}{
		{
			name:     "invalid time format",
			time:     "invalid-time",
			duration: "1h",
			errMsg:   "parsing time",
		},
		{
			name:     "invalid duration format",
			time:     "2023-10-15T13:34:02Z",
			duration: "invalid-duration",
			errMsg:   "time: invalid duration",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := TimeAdd(tt.time, tt.duration)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.errMsg)
		})
	}
}

func TestTimeFormat(t *testing.T) {
	tests := []struct {
		name     string
		time     string
		layout   string
		expected string
	}{
		{
			name:     "valid time and layout",
			time:     "2023-10-15T13:34:02Z",
			layout:   "2006-01-02 15:04:05",
			expected: "2023-10-15 13:34:02",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := TimeFormat(tt.time, tt.layout)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func TestTimeFormat_Invalid(t *testing.T) {
	_, err := TimeFormat("invalid-time", "2006-01-02 15:04:05")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "parsing time")
}
