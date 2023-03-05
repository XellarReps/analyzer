package main

import (
	"analyzer/logparser"
	"errors"
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	var path, calcMode string
	flag.StringVar(&path, "path", "", "***.log file path")
	flag.StringVar(&calcMode, "calc_mode", "", "choose the time counting method (all, mean_epoch_time)")
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

	if calcMode != "all" && calcMode != "mean_epoch_time" {
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
	} else {
		result, err := logparser.CalculateMeanEpochTime(file)
		if err != nil {
			fmt.Println(err)
			return
		}
		resultDur := time.Duration(result * (int64(time.Millisecond) / int64(time.Nanosecond)))
		fmt.Printf("%f\n", resultDur.Minutes())
	}
}
