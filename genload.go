package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
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

type callCounter struct {
	callCount int
	mutex     sync.Mutex
}

func main() {
	totalCallsToMake, maxCallsInParallel, url := parseArgs()

	if maxCallsInParallel == 0 {
		maxCallsInParallel = 1
	}

	if totalCallsToMake == 0 {
		fmt.Printf("Calling url [%s] with [%d] calls simultaneous, until stopped with CTRL+C...\n",
			url, maxCallsInParallel)
	} else {
		fmt.Printf("Calling url [%s] with [%d] calls simultaneous, for a total of [%d] calls...\n",
			url, maxCallsInParallel, totalCallsToMake)
	}

	counter := callCounter{}

	callsMade := make(chan int)

	for i := 0; i < maxCallsInParallel; i++ {
		go makeCalls(i, url, &counter, callsMade)
	}

	for {
		calls := <-callsMade
		if calls >= totalCallsToMake && totalCallsToMake != 0 {
			fmt.Println("All done!")
			return
		}
	}
}

func makeCalls(i int, url string, counter *callCounter, callsMade chan int) {
	for {
		resp, err := http.Get(url)
		if err != nil {
			fmt.Printf("%d: ERROR: %v\n", i, err)
		} else {
			fmt.Printf("%3d: %s\n", i, resp.Status)
		}

		counter.mutex.Lock()
		counter.callCount++
		callsMade <- counter.callCount
		counter.mutex.Unlock()
	}
}
