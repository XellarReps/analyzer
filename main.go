package main

import (
	"analyzer/logparser"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	var path, csvPath, calcMode string
	flag.StringVar(&path, "path", "", "***.log file path")
	flag.StringVar(&csvPath, "csv_path", "", "result csv file (if calc_mode = profile)")
	flag.StringVar(&calcMode, "calc_mode", "", "choose the time counting method (all, mean_epoch"+
		"_time, profile)")
	flag.Parse()

	if path == "" {
		err := errors.New("log file path not specified")
		fmt.Println(err)
		return
	}

	if calcMode == "" {
		err := errors.New("counting method is not selected")
		fmt.Println(err)
		return
	}

	if calcMode != "all" && calcMode != "mean_epoch_time" && calcMode != "profile" {
		err := errors.New("non-existent method (run ./analyzer --help)")
		fmt.Println(err)
		return
	}

	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	if calcMode == "all" {
		result, err := logparser.CalculateAllTime(file)
		if err != nil {
			fmt.Println(err)
			return
		}
		resultDur := time.Duration(result * (int64(time.Millisecond) / int64(time.Nanosecond)))
		fmt.Printf("%f\n", resultDur.Minutes())
	} else if calcMode == "mean_epoch_time" {
		result, err := logparser.CalculateMeanEpochTime(file)
		if err != nil {
			fmt.Println(err)
			return
		}
		resultDur := time.Duration(result * (int64(time.Millisecond) / int64(time.Nanosecond)))
		fmt.Printf("%f\n", resultDur.Minutes())
	} else if calcMode == "profile" {
		if csvPath == "" {
			err := errors.New("csv file path not specified")
			fmt.Println(err)
			return
		}
		outFile, err := os.OpenFile(csvPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
		defer outFile.Close()
		result, err := logparser.CalculateGraphOps(file)
		if err != nil {
			fmt.Println(err)
			return
		}
		newResult := make(map[string]float64)
		for key, value := range result {
			newResult[key] = (time.Duration(value * (int64(time.Millisecond) / int64(time.Nanosecond)))).Minutes()
		}

		records := [][]string{
			{"Op_type", "Time"},
		}

		for node, tm := range newResult {
			timeStr := fmt.Sprintf("%f", tm)
			records = append(records, []string{node, timeStr})
		}

		w := csv.NewWriter(outFile)
		err = w.WriteAll(records)
		if err != nil {
			fmt.Println(err)
			return
		}
		if err := w.Error(); err != nil {
			fmt.Println(err)
			return
		}
	}
}
