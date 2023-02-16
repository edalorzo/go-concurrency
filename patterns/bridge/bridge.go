package main

import (
	"context"
	"fmt"
	"github.com/edalorzo/go-concurrency/patterns/common"
)

func bridge[T any](ctx context.Context, chanStream <-chan <-chan T) <-chan T {

	outStream := make(chan T)

	go func() {
		defer close(outStream)
		for {
			var stream <-chan T

			select {
			case maybeStream, ok := <-chanStream:
				if !ok {
					return
				}
				stream = maybeStream
			case <-ctx.Done():
				return
			}

			for val := range stream {
				select {
				case outStream <- val:
				case <-ctx.Done():
				}
			}
		}
	}()
	return outStream
}

func main() {

	ctx := context.Background()

	chanStream := make(chan (<-chan int))
	outStream := bridge(ctx, chanStream)

	go func() {
		defer close(chanStream)
		chanStream <- common.Producer(ctx, 0, 25)
		chanStream <- common.Producer(ctx, 25, 50)
		chanStream <- common.Producer(ctx, 50, 75)
		chanStream <- common.Producer(ctx, 50, 75)
		chanStream <- common.Producer(ctx, 75, 100)
	}()

	for val := range outStream {
		fmt.Println(val)
	}

}
