package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {
	args := os.Args[1:]
	times, err := strconv.Atoi(args[0])
	url := args[1]

	if err != nil {
		log.Fatal("Usage: genload <times-to-call> <url>")
	}

	if times == 0 {
		fmt.Printf("Calling url [%s] until stopped...\n", url)
	} else {
		fmt.Printf("Calling url [%s] [%d] times...\n", url, times)
	}

	ch := make(chan int, times)

	for i := 0; i < times; i++ {
		go func(i int) {
			resp, err := http.Get(url)
			if err != nil {
				fmt.Printf("%d: ERROR: %v\n", i, err)
			} else {
				fmt.Printf("%3d: %s\n", i, resp.Status)
			}

			ch <- i
		}(i)
	}

	for i := 0; i < times; i++ {
		<-ch
	}
}
