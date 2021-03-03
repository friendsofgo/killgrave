package http

import (
	"encoding/json"
	"errors"
	"math/rand"
	"strings"
	"time"
)

// ResponseDelay represent time delay before server responds.
type ResponseDelay struct {
	delay  int64
	offset int64
}

// Delay return random time.Duration with respect to specified time range.
func (d *ResponseDelay) Delay() time.Duration {
	offset := d.offset
	if offset > 0 {
		offset = rand.Int63n(d.offset)
	}
	return time.Duration(d.delay + offset)
}

// UnmarshalYAML of yaml.Unmarshaler interface.
// Input should be string, consisting of substring that can be parsed by time.ParseDuration,
// or two similar substrings seperated by ":".
func (d *ResponseDelay) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var input string
	if err := unmarshal(&input); err != nil {
		return err
	}

	return d.parseDelay(input)
}

// UnmarshalJSON of json.Unmarshaler interface.
// Input should be string, consisting of substring that can be parsed by time.ParseDuration,
// or two similar substrings seperated by ":".
func (d *ResponseDelay) UnmarshalJSON(data []byte) error {
	var input string
	if err := json.Unmarshal(data, &input); err != nil {
		return err
	}

	return d.parseDelay(input)
}

func (d *ResponseDelay) parseDelay(input string) error {
	const delimiter = ":"

	if input == "" {
		return nil
	}
	inputParts := strings.Split(input, delimiter)
	if len(inputParts) > 2 {
		return errors.New("expected one delimiter at most")
	}
	minDelay, err := time.ParseDuration(inputParts[0])
	if err != nil {
		return err
	}
	var offset int64
	if len(inputParts) == 2 {
		maxDelay, err := time.ParseDuration(inputParts[1])
		if err != nil {
			return err
		}
		offset = int64(maxDelay) - int64(minDelay)
	}
	if offset < 0 {
		return errors.New("second value should be greater than the first one")
	}
	d.delay = int64(minDelay)
	d.offset = offset
	return nil
}
