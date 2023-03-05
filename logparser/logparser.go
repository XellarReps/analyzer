package logparser

import (
	"bufio"
	"encoding/json"
	"os"
)

const (
	epochStart = "epoch_start"
	epochStop  = "epoch_stop"
)

type LogEvent struct {
	Namespace string `json:"namespace"`
	TimeMs    int64  `json:"time_ms"`
	EventType string `json:"event_type"`
	Key       string `json:"key"`
	Value     any    `json:"value"`
}

func CalculateAllTime(file *os.File) (int64, error) {
	scanner := bufio.NewScanner(file)
	var logEvents []LogEvent
	for scanner.Scan() {
		event := scanner.Text()
		event = event[9:] // delete prefix ":::MLLOG"
		var logEvent LogEvent
		err := json.Unmarshal([]byte(event), &logEvent)
		if err != nil {
			return 0, err
		}
		logEvents = append(logEvents, logEvent)
	}
	if err := scanner.Err(); err != nil {
		return 0, err
	}
	if len(logEvents) == 0 {
		return 0, nil
	}
	result := logEvents[len(logEvents)-1].TimeMs - logEvents[0].TimeMs
	return result, nil
}

func CalculateMeanEpochTime(file *os.File) (int64, error) {
	scanner := bufio.NewScanner(file)
	var logEvents []LogEvent
	for scanner.Scan() {
		event := scanner.Text()
		event = event[9:] // delete prefix ":::MLLOG"
		var logEvent LogEvent
		err := json.Unmarshal([]byte(event), &logEvent)
		if err != nil {
			return 0, err
		}
		if logEvent.Key == epochStart || logEvent.Key == epochStop {
			logEvents = append(logEvents, logEvent)
		}
	}
	if err := scanner.Err(); err != nil {
		return 0, err
	}
	if len(logEvents) == 0 || len(logEvents)%2 != 0 {
		return 0, nil
	}
	var sumTimeEpoch int64
	var cnt int64
	for i := 0; i < len(logEvents)-1; i += 2 {
		start := logEvents[i].TimeMs
		stop := logEvents[i+1].TimeMs
		sumTimeEpoch += stop - start
		cnt++
	}
	return sumTimeEpoch / cnt, nil
}
