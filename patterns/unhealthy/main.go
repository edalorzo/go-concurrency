package main

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"time"
)

type startGoroutineFn func(ctx context.Context, pulseInterval time.Duration) (heartbeat <-chan any)

func newSteward(timeout time.Duration, startGoroutine startGoroutineFn) startGoroutineFn {

	return func(ctx context.Context, pulseInterval time.Duration) <-chan any {

		heartbeat := make(chan any)

		go func() {
			defer close(heartbeat)

			var (
				wardContext   context.Context
				wardCancelFn  context.CancelFunc
				wardHeartbeat <-chan any
			)
			startWard := func() {
				wardContext, wardCancelFn = context.WithCancel(ctx)
				wardHeartbeat = startGoroutine(wardContext, pulseInterval/2)
			}
			startWard()
			pulse := time.Tick(pulseInterval)

		monitorLoop:
			for {
				// notice the timeout is only restarted when we receive a heartbeat
				timeoutSignal := time.After(timeout)
				for {
					select {
					case <-pulse:
						select {
						case heartbeat <- struct{}{}:
						default:
						}
					case <-wardHeartbeat:
						continue monitorLoop
					case <-timeoutSignal:
						fmt.Println("steward: ward unhealthy; restarting")
						wardCancelFn()
						startWard()
						continue monitorLoop
					case <-ctx.Done():
						wardCancelFn()
						return
					}
				}
			}
		}()

		return heartbeat
	}
}

func doWork(ctx context.Context, pulseInterval time.Duration) <-chan any {
	logrus.Info("ward: Hello, I'm responsible")
	go func() {
		<-ctx.Done()
		logrus.Info("ward: I'm halting")
	}()
	return nil
}

func main() {

	ctx, cancelFunc := context.WithCancel(context.Background())
	doWithSteward := newSteward(4*time.Second, doWork)

	// let our code run for 10 seconds
	time.AfterFunc(20*time.Second, func() {
		logrus.Info("main: halting steward and ward")
		cancelFunc()
	})

	for range doWithSteward(ctx, 4*time.Second) {
	}

	logrus.Info("Done")

}
