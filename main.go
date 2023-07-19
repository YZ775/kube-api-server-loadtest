package main

import (
	"context"
	"fmt"
	"io"
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
	// kubeSummaryExporterAPIEndpoint := fmt.Sprintf("http://127.0.0.1:9779/node/%s", nodeName)

	// ApiPath := fmt.Sprintf("/api/v1/nodes/%s/proxy/stats/summary", nodeName)

	fmt.Printf("sending %d requests/sec to %s\n...\n", NUM_REQUESTS, summaryAPIEndpoint)
	// fmt.Printf("executing %d ops/sec to %s\n...\n", NUM_REQUESTS, ApiPath)

	// req.Header.Set("X-Prometheus-Scrape-Timeout-Seconds", "20")

	for {
		wg := sync.WaitGroup{}
		for i := 0; i < NUM_REQUESTS; i++ {
			wg.Add(1)
			go func() {
				// err := exec.Command("kubectl", "get", "--raw", ApiPath).Run()
				// if err != nil {
				// 	fmt.Println(err)
				// }

				// res, err := http.Get(kubeSummaryExporterAPIEndpoint)

				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				req, _ := http.NewRequestWithContext(ctx, "GET", summaryAPIEndpoint, nil)
				client := new(http.Client)
				res, err := client.Do(req)
				if err != nil {
					fmt.Println(err)
					wg.Done()
					return
				}
				if res.StatusCode != 200 {
					b, _ := io.ReadAll(res.Body)
					fmt.Printf("statusCode:%d, %v\n", res.StatusCode, string(b))
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
