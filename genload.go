package main

import (
	"fmt"
	"log"
	"net/http"
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
	return
}

type callResponse struct {
	threadNr        int
	resp            *http.Response
	err             error
	continueChannel chan bool
}

func main() {
	totalCallsToMake, maxCallsInParallel, url := parseArgs()

	if maxCallsInParallel == 0 {
		maxCallsInParallel = 1
	}
	if maxCallsInParallel > totalCallsToMake && totalCallsToMake != 0 {
		maxCallsInParallel = totalCallsToMake
	}

	if totalCallsToMake == 0 {
		fmt.Printf("Calling url [%s] with [%d] calls simultaneous, until stopped with CTRL+C...\n",
			url, maxCallsInParallel)
	} else {
		fmt.Printf("Calling url [%s] with [%d] calls simultaneous, for a total of [%d] calls...\n",
			url, maxCallsInParallel, totalCallsToMake)
	}

	reportBackChannel := make(chan callResponse)

	for threadNr := 0; threadNr < maxCallsInParallel; threadNr++ {
		continueChannel := make(chan bool)
		go makeCall(threadNr, url, reportBackChannel, continueChannel)
	}

	callCounter := maxCallsInParallel
	runningThreads := maxCallsInParallel
	counter := 1

	for {
		callReport := <-reportBackChannel

		if callCounter < totalCallsToMake || totalCallsToMake == 0 {
			callCounter++
			callReport.continueChannel <- true
		} else {
			callReport.continueChannel <- false
			runningThreads--
		}

		if callReport.err != nil {
			fmt.Printf("%5d: [%3d]: ERROR: %v\n", counter, callReport.threadNr, callReport.err)
		} else {
			fmt.Printf("%5d: [%3d]: %s\n", counter, callReport.threadNr, callReport.resp.Status)
		}
		counter++

		if runningThreads == 0 {
			return
		}
	}
}

func makeCall(threadNr int, url string, reportBackChannel chan callResponse, continueChannel chan bool) {
	for {
		resp, err := http.Get(url)
		callResp := callResponse{threadNr, resp, err, continueChannel}

		reportBackChannel <- callResp

		startNextCall := <-continueChannel
		if !startNextCall {
			return
		}
	}
}
