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

func main() {
	fmt.Fprintf(os.Stdout, "success/s\tfail/s\n")
	var success, failure int
	start := time.Now()
	logged := start
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns: 1,
		},
	}
	for time.Since(start) < 10*time.Second {
		resp, err := client.Get("https://example.com/")
		if err != nil {
			log.Print(err)
			failure++
			continue
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Print(err)
			failure++
			continue
		}
		hash := fnv.New32a()
		_, err = hash.Write(body)
		if err != nil {
			log.Print(err)
			failure++
			continue
		}
		if hash.Sum32() != 3712873988 {
			log.Print(hash.Sum32())
			failure++
			continue
		}
		success++
		if time.Since(logged) > time.Second {
			runTime := time.Since(logged)
			successRate := int(float64(success) / runTime.Seconds())
			failureRate := int(float64(failure) / runTime.Seconds())
			fmt.Fprintf(os.Stdout, "%d\t%d\n", successRate, failureRate)
			success = 0
			failure = 0
			logged = time.Now()
		}
		client.CloseIdleConnections()
	}
}
