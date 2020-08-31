package main

import (
	"fmt"
	"net/http"
	"time"
)

type callResponse struct {
	threadNr        int
	resp            *http.Response
	err             error
	continueChannel chan bool
	elapsedTime     time.Duration
}

func (cr *callResponse) ElapsedTimeInMillis() float32 {

	return float32(cr.elapsedTime.Microseconds()) / 1000
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
		ForceAttemptHTTP2: true,
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
			fmt.Printf("%5d: [%3d]: <%5.1f>: ERROR: %v\n",
				counter, callReport.threadNr, callReport.ElapsedTimeInMillis(), callReport.err)
		} else {
			fmt.Printf("%5d: [%3d]: <%5.1f>:  %s\n",
				counter, callReport.threadNr, callReport.ElapsedTimeInMillis(), callReport.resp.Status)
		}
		counter++

		if runningThreads == 0 {
			return
		}
	}
}

func makeCall(threadNr int, client *http.Client, url string, reportBackChannel chan callResponse, continueChannel chan bool) {
	for {
		startTime := time.Now()
		resp, err := client.Get(url)
		if resp != nil {
			_ = resp.Body.Close()
		}
		endTime := time.Now()
		elapsedTime := endTime.Sub(startTime)

		callResp := callResponse{threadNr, resp, err, continueChannel, elapsedTime}

		reportBackChannel <- callResp

		startNextCall := <-continueChannel
		if !startNextCall {
			return
		}
	}
}
