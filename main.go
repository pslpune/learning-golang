package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

const (
	MAX_THREAD_POOL = 10
)

// Parallel - what do you mean by that ? Core 1 (CPU1) - task A || Core 2 (CPU2) - task B
// Core 2 Duo Core 1 (4 threads) (clocks, Clockd speed 2 Ghz - in a single sec it is capable of taking 2x10^12 instructions , FSB) = Thread 1 -task A || Thread 2 - TaskB || Thread 3 || Thread4
// Being AGNOSTIC of whether it is scheduled on single core or multicore when a language lets you program uniformly its called *asynchrnoy*
func main() {
	fmt.Println("now learning asynchrony in Go")
	// comic, err := DownloadComic(100)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(comic["safe_title"])
	jobs := []int{}
	for i := 100; i < 120; i++ {
		jobs = append(jobs, i)
	}
	// chnResults := make(chan map[string]interface{}, len(jobs))
	// defer close(chnResults)
	// task - 1,2,3,4,5  .. 20 Asynchronous
	// task 4 is picked up for execution
	var wg sync.WaitGroup
	chnResults := make(chan map[string]interface{}, len(jobs))
	defer close(chnResults)
	// a pool of go-routines (referred to as thread pools in other platforms)
	// a fixed set of go routines then cyclically takes up all the jobs one after other.
	// this saves penalty for creation, destroy and somewhat from scheduling each go-routine.

	for _, j := range jobs {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			result, err := DownloadComic(idx)
			if err != nil {
				fmt.Println("We have an error downloading the comic")
				return
			}
			// fmt.Println(result["safe_title"])
			chnResults <- result
		}(j)
	}
	go func(chnRead chan map[string]interface{}) {
		for r := range chnRead {
			fmt.Println(r["safe_title"])
		}
	}(chnResults)
	wg.Wait()
	// <- time.After(1*time.Second)
	fmt.Println("End of the program")
}

func DownloadComic(comicIndex int) (map[string]interface{}, error) {
	url := fmt.Sprintf("https://xkcd.com/%d/info.0.json", comicIndex)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create new request %s", err)
	}
	cl := &http.Client{
		Timeout: 4 * time.Second,
	}
	resp, err := cl.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request over http %s", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get favorable response from server %d", resp.StatusCode)
	}
	// this is wehre we read the json response payload
	defer resp.Body.Close()
	byt, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("invalid json payload from server %s", err)
	}
	// defer // this is the first time i have used it ! .. will come to this a bit later
	result := map[string]interface{}{}
	err = json.Unmarshal(byt, &result)
	if err != nil {
		return nil, fmt.Errorf("fauled to read payload from server %s", err)
	}
	return result, nil
}
