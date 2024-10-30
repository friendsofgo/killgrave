package templating

import (
	"encoding/json"
	"fmt"
	"time"
)

func JsonMarshal(v interface{}) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func TimeNow() string {
	return time.Now().Format(time.RFC3339)
}

func TimeUTC(t string) (string, error) {
	parsedTime, err := time.Parse(time.RFC3339, t)
	if err != nil {
		return "", fmt.Errorf("error parsing time: %v", err)
	}
	return parsedTime.UTC().Format(time.RFC3339), nil
}

func TimeAdd(t string, d string) (string, error) {
	parsedTime, err := time.Parse(time.RFC3339, t)
	if err != nil {
		return "", fmt.Errorf("error parsing time: %v", err)
	}
	duration, err := time.ParseDuration(d)
	if err != nil {
		return "", fmt.Errorf("error parsing duration: %v", err)
	}
	return parsedTime.Add(duration).Format(time.RFC3339), nil
}

func TimeFormat(t string, layout string) (string, error) {
	parsedTime, err := time.Parse(time.RFC3339, t)
	if err != nil {
		return "", fmt.Errorf("error parsing time: %v", err)
	}
	return parsedTime.Format(layout), nil
}
