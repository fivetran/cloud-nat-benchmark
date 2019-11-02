package main

import (
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

const runTime = 2 * time.Second

type counts struct {
	success, failure int
}

var tryThreads = []int{
	1,
	10,
	25,
	50,
	100,
	200,
	300,
	400,
	500,
}

func main() {
	fmt.Fprintf(os.Stdout, "threads,success/s,fail/s\n")
	for _, threads := range tryThreads {
		reply := make(chan counts, threads)
		for i := 0; i < threads; i++ {
			go createRequests(reply)
		}
		var total counts
		for n := 0; n < threads; n++ {
			partial := <-reply
			total.success += partial.success
			total.failure += partial.failure
		}
		successRate := int(float64(total.success) / runTime.Seconds())
		failureRate := int(float64(total.failure) / runTime.Seconds())
		fmt.Fprintf(os.Stdout, "%d,%d,%d\n", threads, successRate, failureRate)
	}
}

func createRequests(reply chan counts) {
	var result counts
	start := time.Now()
	for time.Since(start) < runTime {
		client := &http.Client{
			Transport: &http.Transport{
				MaxIdleConns: 1,
			},
		}
		resp, err := client.Get("https://example.com/")
		if err != nil {
			log.Print(err)
			result.failure++
			continue
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Print(err)
			result.failure++
			continue
		}
		hash := fnv.New32a()
		_, err = hash.Write(body)
		if err != nil {
			log.Print(err)
			result.failure++
			continue
		}
		if hash.Sum32() != 3712873988 {
			log.Print(hash.Sum32())
			result.failure++
			continue
		}
		result.success++
	}
	reply <- result
}
