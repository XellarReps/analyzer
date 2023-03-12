package logparser

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

const (
	epochStart  = "epoch_start"
	epochStop   = "epoch_stop"
	startPrefix = "start_"
	stopPrefix  = "stop_"
)

type LogEvent struct {
	TimeMs int64  `json:"time_ms"`
	Key    string `json:"key"`
	Value  any    `json:"value"`
}

type GraphOp struct {
	NameOp string
	Number any
	TimeMs int64
}

type Stack struct {
	items []GraphOp
}

func (s *Stack) Push(data GraphOp) {
	s.items = append(s.items, data)
}

func (s *Stack) Pop() {
	if s.IsEmpty() {
		return
	}
	s.items = s.items[:len(s.items)-1]
}

func (s *Stack) Top() (GraphOp, error) {
	if s.IsEmpty() {
		return GraphOp{}, fmt.Errorf("stack is empty")
	}
	return s.items[len(s.items)-1], nil
}

func (s *Stack) IsEmpty() bool {
	if len(s.items) == 0 {
		return true
	}
	return false
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
	for i := 0; i < len(logEvents)-1; i += 2 {
		start := logEvents[i].TimeMs
		stop := logEvents[i+1].TimeMs
		sumTimeEpoch += stop - start
	}
	return sumTimeEpoch, nil
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

func CalculateGraphOps(file *os.File) (map[string]int64, error) {
	scanner := bufio.NewScanner(file)
	timeOps := make(map[string]int64)
	var st Stack
	for scanner.Scan() {
		event := scanner.Text()
		event = event[9:] // delete prefix ":::MLLOG"
		var logEvent LogEvent
		err := json.Unmarshal([]byte(event), &logEvent)
		if err != nil {
			return nil, err
		}
		if len(logEvent.Key) >= len(startPrefix) && logEvent.Key[:len(startPrefix)] == startPrefix {
			opName := logEvent.Key[len(startPrefix):]
			tm := logEvent.TimeMs
			val := logEvent.Value
			st.Push(GraphOp{
				NameOp: opName,
				TimeMs: tm,
				Number: val,
			})
		} else if len(logEvent.Key) >= len(stopPrefix) && logEvent.Key[:len(stopPrefix)] == stopPrefix {
			opName := logEvent.Key[len(stopPrefix):]
			tm := logEvent.TimeMs
			val := logEvent.Value
			operation, err := st.Top()
			if err != nil {
				return nil, err
			}
			st.Pop()
			if operation.NameOp != opName && operation.Number != val {
				return nil, fmt.Errorf("incorrect sequence")
			}
			timeOps[opName] += tm - operation.TimeMs
		}
	}
	return timeOps, nil
}
