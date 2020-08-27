package main

import (
	"fmt"
	"net/http"
)

type callResponse struct {
	threadNr        int
	resp            *http.Response
	err             error
	continueChannel chan bool
}

func main() {

	totalCallsToMake, maxCallsInParallel, url := parseArgs()

	if totalCallsToMake == 0 {
		fmt.Printf("Calling url [%s] with [%d] calls simultaneous, until stopped with CTRL+C...\n",
			url, maxCallsInParallel)
	} else {
		fmt.Printf("Calling url [%s] with [%d] calls simultaneous, for a total of [%d] calls...\n",
			url, maxCallsInParallel, totalCallsToMake)
	}

	tr := &http.Transport{
		DisableKeepAlives:      true,
	}
	client := &http.Client{
		Transport: tr,
	}

	reportBackChannel := make(chan callResponse)

	for threadNr := 0; threadNr < maxCallsInParallel; threadNr++ {
		continueChannel := make(chan bool)
		go makeCall(threadNr, client, url, reportBackChannel, continueChannel)
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

func makeCall(threadNr int, client *http.Client, url string, reportBackChannel chan callResponse, continueChannel chan bool) {
	for {
		resp, err := client.Get(url)
		if resp != nil {
			_ = resp.Body.Close()
		}
		callResp := callResponse{threadNr, resp, err, continueChannel}

		reportBackChannel <- callResp

		startNextCall := <-continueChannel
		if !startNextCall {
			return
		}
	}
}
