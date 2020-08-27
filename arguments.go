package main

import (
	"log"
	"os"
	"strconv"
)

func showUsageAndQuit() {
	log.Fatal("Usage: genload <total-calls> <max-parallel-calls> <url>")
}

func parseArgs() (totalCallsToMake int, maxCallsInParallel int, url string) {

	var err error

	switch len(os.Args) {
	case 3:
		totalCallsToMake, err = strconv.Atoi(os.Args[1])
		if err != nil {
			showUsageAndQuit()
		}
		maxCallsInParallel = 0
		url = os.Args[2]
	case 4:
		totalCallsToMake, err = strconv.Atoi(os.Args[1])
		if err != nil {
			showUsageAndQuit()
		}
		maxCallsInParallel, err = strconv.Atoi(os.Args[2])
		if err != nil {
			showUsageAndQuit()
		}
		url = os.Args[3]
	default:
		showUsageAndQuit()
	}

	if totalCallsToMake < 0 {
		totalCallsToMake = 0
	}

	if maxCallsInParallel < 1 {
		maxCallsInParallel = 1
	}

	if maxCallsInParallel > totalCallsToMake && totalCallsToMake != 0 {
		maxCallsInParallel = totalCallsToMake
	}

	return
}
