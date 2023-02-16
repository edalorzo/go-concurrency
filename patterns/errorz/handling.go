package main

import (
	"context"
	"github.com/sirupsen/logrus"
	"net/http"
)

/*
Here we see a gouroutine that is doing its best to signal that there's an error.
What else can it do? It can't pass it back. How many errors is too many?
Should it continue making requests?
*/

func checkStatus(ctx context.Context, urls ...string) <-chan *http.Response {
	responses := make(chan *http.Response)
	go func() {
		defer close(responses)
		for _, url := range urls {
			resp, err := http.Get(url)
			if err != nil {
				logrus.Error("Oh uh", err)
				continue
			}
			select {
			case responses <- resp:
			case <-ctx.Done():
			}
		}
	}()
	return responses
}

type Result struct {
	Error    error
	Response *http.Response
}

// This is desirable because the goroutine that spawned the producer goroutine
// has more context about the running program and can make more intelligent decisions
// about what to do with the errors.
func checkStatusFixed(ctx context.Context, urls ...string) <-chan Result {
	results := make(chan Result)
	go func() {
		defer close(results)
		for _, url := range urls {
			resp, err := http.Get(url)
			result := Result{Error: err, Response: resp}
			select {
			case results <- result:
			case <-ctx.Done():
			}
		}
	}()
	return results
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	urls := []string{
		"https://www.google.com",
		"https://www.badhost.com",
		//"https://www.badhost.com",
		"https://www.google.com",
		"https://www.google.com",
	}

	errCount := 0
	for result := range checkStatusFixed(ctx, urls...) {
		if result.Error != nil {
			errCount++
			if errCount >= 2 {
				cancel()
			}
			logrus.Error("Oh uh", result.Error)
			continue
		}
		logrus.Info("Response:", result.Response.Status)
	}
}
