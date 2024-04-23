package cmd

import (
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/spf13/cobra"
)

var loadCmd = &cobra.Command{
	Use:   "load",
	Short: "load test a url by sending HTTP requests to it",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// get url flag
		url, err := cmd.Flags().GetString("url")
		if err != nil {
			fmt.Println(err)
			return
		}
		err = validateUrl(url)
		if err != nil {
			fmt.Println(err)
			return
		}

		// get number of requests flag
		n, err := cmd.Flags().GetInt("number")
		if err != nil {
			fmt.Println(err)
			return
		}
		err = validateNumber(n)
		if err != nil {
			fmt.Println(err)
			return
		}

		// get concurrency flag
		c, err := cmd.Flags().GetInt("concurrency")
		if err != nil {
			fmt.Println(err)
			return
		}
		err = validateConcurrency(c, n)
		if err != nil {
			fmt.Println(err)
			return
		}

		// measure response times
		response_times := make([]int64, n)

		// send requests
		var n_success atomic.Uint64
		concurrency_counter := c
		wg := sync.WaitGroup{}
		for i := 0; i < n; i++ {
			wg.Add(1)
			go func(count int) {
				defer wg.Done()
				status_code, request_time, err := makeRequest(url)
				response_times[count] = int64(request_time)
				if err != nil {
					fmt.Println(err)
					return
				}
				if status_code <= 299 && status_code >= 200 {
					n_success.Add(1)
				}
			}(i)
			concurrency_counter -= 1
			if concurrency_counter == 0 {
				wg.Wait()
				concurrency_counter = c
			}
		}
		wg.Wait()
		var total_request_time int64 = 0
		for _, val := range response_times {
			total_request_time += val
		}
		average_request_time := total_request_time / int64(n_success.Load())
		fmt.Printf("success rate: %d/%d\n", n_success.Load(), n)
		fmt.Printf("average request time: %f", float32(average_request_time)/float32(1_000_000_000))
	},
}

func init() {
	rootCmd.AddCommand(loadCmd)
	loadCmd.Flags().StringP("url", "u", "", "the url to send requests to")
	loadCmd.Flags().IntP("number", "n", 1, "the number of requests to send to that url")
	loadCmd.Flags().IntP("concurrency", "c", 1, "the number of requests to send concurrently")
}

func validateUrl(url string) error {
	first_4 := url[0:4]
	if first_4 != "http" {
		return fmt.Errorf("invalid url, expected url to start with http, got %s", first_4)
	}
	return nil
}

func validateNumber(n int) error {
	if n <= 0 {
		return fmt.Errorf("please enter a number greater than 0")
	}
	if n > 10_000 {
		return fmt.Errorf("number too large")
	}
	return nil
}

func validateConcurrency(c int, n int) error {
	if c > n {
		return fmt.Errorf("cannot concurrently send more requests than number specified. %d > %d", c, n)
	}
	return nil
}

type StatusCode int

func makeRequest(url string) (StatusCode, time.Duration, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, 0, err
	}
	client := http.Client{}
	start := time.Now()
	resp, err := client.Do(req)
	request_time := time.Since(start)
	if err != nil {
		return 0, request_time, err
	}
	return StatusCode(resp.StatusCode), request_time, nil
}
