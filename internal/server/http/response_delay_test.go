package http

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestResponseDelayUnmarshal(t *testing.T) {
	testCases := map[string]struct {
		input   string
		delay   ResponseDelay
		wantErr bool
	}{
		"Invalid type": {
			input:   `23`,
			wantErr: true,
		},
		"Valid empty delay": {
			input:   `""`,
			delay:   ResponseDelay{0, 0},
			wantErr: false,
		},
		"Valid fixed delay": {
			input:   `"1s"`,
			delay:   getDelay(t, "1s", "0s"),
			wantErr: false,
		},
		"Fixed delay without unit suffix": {
			input:   `"13"`,
			wantErr: true,
		},
		"Valid range delay": {
			input:   `"2s:7s"`,
			delay:   getDelay(t, "2s", "5s"),
			wantErr: false,
		},
		"Range delay with incorrect delimiter": {
			input:   `"1m-3s"`,
			wantErr: true,
		},
		"Range delay with extra field": {
			input:   `"1m:3s:5s"`,
			wantErr: true,
		},
		"Range delay where second point is before first": {
			input:   `"5s:1s"`,
			wantErr: true,
		},
		"Range delay where second point invalid": {
			input:   `"5s:1"`,
			wantErr: true,
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var delay ResponseDelay
			err := json.Unmarshal([]byte(tc.input), &delay)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.delay, delay)

		})
	}
}

func TestResponseDelay(t *testing.T) {
	testCases := map[string]struct {
		delay ResponseDelay
	}{
		"Empty delay": {
			delay: ResponseDelay{0, 0},
		},
		"Fixed delay": {
			delay: getDelay(t, "2s", "0s"),
		},
		"Range delay": {
			delay: getDelay(t, "2s", "5s"),
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			min := tc.delay.delay
			max := min + tc.delay.offset
			for i := 0; i < 10; i++ {
				delay := int64(tc.delay.Delay())
				assert.GreaterOrEqual(t, max, delay)
				assert.GreaterOrEqual(t, delay, min)
			}

		})
	}
}

func getDelay(t *testing.T, min string, offset string) ResponseDelay {
	minDuration, err := time.ParseDuration(min)
	assert.Nil(t, err)
	offsetDuration, err := time.ParseDuration(offset)
	assert.Nil(t, err)
	return ResponseDelay{int64(minDuration), int64(offsetDuration)}
}
