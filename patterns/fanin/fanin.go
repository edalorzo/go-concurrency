package main

import (
	"context"
	"fmt"
	"github.com/edalorzo/go-concurrency/patterns/common"
	"sync"
)

// this is a fan-in pattern implementation where order of the elements does not matter
func fanIn[T any](ctx context.Context, channels ...<-chan T) <-chan T {

	var wg sync.WaitGroup
	multiplexedStream := make(chan T)

	multiplex := func(c <-chan T) {
		defer wg.Done()
		for val := range c {
			select {
			case multiplexedStream <- val:
			case <-ctx.Done():
			}
		}
	}

	wg.Add(len(channels))
	for _, c := range channels {
		go multiplex(c)
	}

	go func() {
		wg.Wait()
		close(multiplexedStream)
	}()

	return multiplexedStream
}

func main() {
	ctx := context.Background()

	producers := []<-chan int{
		common.Producer(ctx, 0, 25),
		common.Producer(ctx, 25, 50),
		common.Producer(ctx, 50, 75),
		common.Producer(ctx, 75, 100),
	}

	combined := fanIn(ctx, producers...)

	for val := range combined {
		fmt.Println(val)
	}

}
