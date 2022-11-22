package killgrave

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"
	"time"
)

func TestResponseDelayMarshal(t *testing.T) {
	for _, tc := range []string{`"1s"`, `"2s:7s"`} {
		t.Run(tc, func(t *testing.T) {
			var rd ResponseDelay
			err := json.Unmarshal([]byte(tc), &rd)
			if err != nil {
				t.Fatal("json.Unmarshal fail: ", err)
			}

			raw, err := json.Marshal(&rd)
			if err != nil {
				t.Fatal("json.Marshal fail: ", err)
			}

			if string(raw) != tc {
				t.Fatalf("expected: %v, got: %v", tc, string(raw))
			}
		})
	}
}

func TestResponseDelayUnmarshal(t *testing.T) {
	testCases := map[string]struct {
		input string
		delay ResponseDelay
		err   error
	}{
		"Invalid type": {
			input: `23`,
			err:   errors.New("error"),
		},
		"Valid empty delay": {
			input: `""`,
			delay: ResponseDelay{delay: 0, offset: 0},
		},
		"Valid fixed delay": {
			input: `"1s"`,
			delay: getDelay(t, "1s", "0s"),
		},
		"Fixed delay without unit suffix": {
			input: `"13"`,
			err:   errors.New("error"),
		},
		"Valid range delay": {
			input: `"2s:7s"`,
			delay: getDelay(t, "2s", "5s"),
		},
		"Range delay with incorrect delimiter": {
			input: `"1m-3s"`,
			err:   errors.New("error"),
		},
		"Range delay with extra field": {
			input: `"1m:3s:5s"`,
			err:   errors.New("error"),
		},
		"Range delay where second point is before first": {
			input: `"5s:1s"`,
			err:   errors.New("error"),
		},
		"Range delay where second point invalid": {
			input: `"5s:1"`,
			err:   errors.New("error"),
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			var delay ResponseDelay
			err := json.Unmarshal([]byte(tc.input), &delay)
			if err != nil && tc.err == nil {
				t.Fatalf("not expected any error and got: %v", err)
			}

			if err == nil && tc.err != nil {
				t.Fatalf("expected an error and got nil")
			}
			if !reflect.DeepEqual(tc.delay, delay) {
				t.Fatalf("expected: %v, got: %v", tc.delay, delay)
			}

		})
	}
}

func TestResponseDelay(t *testing.T) {
	testCases := map[string]struct {
		delay ResponseDelay
	}{
		"Empty delay": {
			delay: ResponseDelay{delay: 0, offset: 0},
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
				if delay < min || delay > max {
					t.Errorf("delay should be in interval [%d, %d] but actual value is: %d", min, max, delay)
				}
			}

		})
	}
}

func getDelay(t *testing.T, min string, offset string) ResponseDelay {
	minDuration, err := time.ParseDuration(min)
	if err != nil {
		t.Fatal("ParseDuration min fail: ", err)
	}
	offsetDuration, err := time.ParseDuration(offset)
	if err != nil {
		t.Fatal("ParseDuration max fail: ", err)
	}
	return ResponseDelay{delay: int64(minDuration), offset: int64(offsetDuration)}
}
