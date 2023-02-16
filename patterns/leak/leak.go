package main

import (
	"context"
	"github.com/sirupsen/logrus"
	"math/rand"
	"time"
)

/*
If a goroutine is responsible for creating a goroutine, it is responsible for
ensuring that goroutine is properly terminated when the parent goroutine exits.
*/

func doWork(strings <-chan string) <-chan any {
	completed := make(chan any)
	go func() {
		defer logrus.Info("doWork() exited")
		defer close(completed)
		for s := range strings {
			logrus.Info(s)
		}
	}()
	return completed
}

func doWorkFixed(ctx context.Context, strings <-chan string) <-chan any {
	completed := make(chan any)
	go func() {
		defer logrus.Info("doWork() exited")
		defer close(completed)
		for {
			select {
			case s, ok := <-strings:
				if !ok {
					return
				}
				logrus.Info(s)
			case <-ctx.Done():
				return
			}
		}
	}()
	return completed
}

func newRandStream() <-chan int {
	randStream := make(chan int)
	go func() {
		defer logrus.Info("newRandStream() exited")
		defer close(randStream)
		for {
			randStream <- rand.Int()
		}
	}()
	return randStream
}

func newRandStreamFixed(ctx context.Context) <-chan int {
	randStream := make(chan int)
	go func() {
		defer logrus.Info("newRandStream() exited")
		defer close(randStream)
		for {
			select {
			case randStream <- rand.Int():
			case <-ctx.Done():
				return
			}

		}
	}()
	return randStream
}

//func main() {
//	ctx, cancelFn := context.WithCancel(context.Background())
//	doWorkFixed(ctx, nil)
//	<-time.After(3 * time.Second)
//	cancelFn()
//	<-time.After(3 * time.Second)
//	logrus.Info("Done!")
//}

func main() {
	ctx := context.Background()
	ctx, cancelFn := context.WithCancel(ctx)

	randStream := newRandStreamFixed(ctx)
	logrus.Info("3 random ints:")
	for i := 0; i < 3; i++ {
		logrus.Info(<-randStream)
	}
	cancelFn()
	<-time.After(1 * time.Second)
	logrus.Info("Done!")
}
