package main

import (
	"context"
	"fmt"
	"time"
)

func doWork(ctx context.Context, pulseInterval time.Duration) (<-chan any, <-chan time.Time) {

	heartbeat := make(chan any)
	results := make(chan time.Time)

	go func() {

		// comment to simulate unhealthy goroutine
		//defer close(heartbeat)
		//defer close(results)

		pulse := time.Tick(pulseInterval) // emits a time tick every pulseInterval
		sendPulse := func() {
			select {
			case heartbeat <- struct{}{}:
			default:
			}
		}

		sendResult := func(result time.Time) {
			for {
				select {
				case <-pulse:
					sendPulse()
				case results <- result:
					return
				case <-ctx.Done():
					return
				}
			}
		}

		workGen := time.Tick(2 * pulseInterval)
		//uncomment to demonstrate unhealthy goroutine
		for i := 0; i < 2; i++ {
			select {
			case <-pulse:
				sendPulse()
			case r := <-workGen:
				sendResult(r)
			case <-ctx.Done():
				return
			}
		}
	}()

	return heartbeat, results
}

func main() {
	ctx, cancelFn := context.WithCancel(context.Background())
	time.AfterFunc(10*time.Second, cancelFn)

	const timeout = 2 * time.Second

	heartbeat, results := doWork(ctx, timeout/2)
	for {
		select {
		case _, ok := <-heartbeat:
			if !ok {
				return
			}
			fmt.Println("pulse")
		case r, ok := <-results:
			if !ok {
				return
			}
			fmt.Printf("results %v\n", r.Second())
		case <-time.After(timeout):
			fmt.Println("worker goroutine is not healthy!")
			// dismantle the unhealthy goroutine here!
			// cancelFn() and restart!
			return
		}
	}
}
