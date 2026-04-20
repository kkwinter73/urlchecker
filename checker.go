package main

import (
	"fmt"
	"net/http"
	"time"
)

type Result struct {
	URL        string
	StatusCode int
	Duration   time.Duration
	Err        error
}

// String は結果を人間が読める形に整形する
func (r Result) String() string {
	if r.Err != nil {
		return fmt.Sprintf("[NG] %s -> error: %v (%.2fs)", r.URL, r.Err, r.Duration.Seconds())
	}
	return fmt.Sprintf("[OK] %s -> %d (%.2fs)", r.URL, r.StatusCode, r.Duration.Seconds())
}

func Check(url string) Result {
	start := time.Now()

	resp, err := http.Get(url)
	duration := time.Since(start)

	if err != nil {
		return Result{URL: url, Duration: duration, Err: err}
	}

	defer resp.Body.Close()

	return Result{
		URL:        url,
		StatusCode: resp.StatusCode,
		Duration:   duration,
	}
}

func CheckAll(urls []string) []Result {
	results := make([]Result, len(urls))
	ch := make(chan struct {
		index  int
		result Result
	}, len(urls))

	for i, url := range urls {
		go func(index int, u string) {
			ch <- struct {
				index  int
				result Result
			}{index, Check(u)}
		}(i, url)
	}

	for range urls {
		r := <-ch
		results[r.index] = r.result
	}

	return results
}
