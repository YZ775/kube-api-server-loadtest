package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"
)

func main() {
	nodeName := os.Args[1]
	NUM_REQUESTS, err := strconv.Atoi(os.Args[2])
	if err != nil {
		panic(err)
	}
	summaryAPIEndpoint := fmt.Sprintf("http://127.0.0.1:8001/api/v1/nodes/%s/proxy/stats/summary", nodeName)

	fmt.Printf("sending %d requests/sec to %s\n...\n", NUM_REQUESTS, summaryAPIEndpoint)
	for {
		wg := sync.WaitGroup{}
		for i := 0; i < NUM_REQUESTS; i++ {
			wg.Add(1)
			go func() {
				res, err := http.Get(summaryAPIEndpoint)
				if err != nil {
					fmt.Println(err)
				}
				if res.StatusCode != 200 {
					fmt.Printf("statusCode:%s\n", res.Status)
				}
				defer res.Body.Close()
				wg.Done()
			}()
		}
		wg.Wait()
		fmt.Printf("%d requests sent\n", NUM_REQUESTS)
		time.Sleep(1 * time.Second)
	}

}
