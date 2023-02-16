package main

import (
	"context"
	"fmt"
	"sync"
)

func tee[T any](ctx context.Context, in <-chan T) (<-chan T, <-chan T) {

	out1 := make(chan T)
	out2 := make(chan T)

	go func() {

		defer close(out1)
		defer close(out2)

		for val := range in {
			var out1, out2 = out1, out2
			for i := 0; i < 2; i++ {
				select {
				case out1 <- val:
					out1 = nil
				case out2 <- val:
					out2 = nil
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	return out1, out2
}

func main() {
	ctx := context.Background()

	ch := make(chan int)
	out1, out2 := tee(ctx, ch)

	go func() {
		defer close(ch)
		for i := 0; i < 10; i++ {
			ch <- i
		}
	}()

	wg := sync.WaitGroup{}
	wg.Add(2)

	consumer := func(out <-chan int) {
		defer wg.Done()
		for val := range out {
			fmt.Printf("%d ", val)
		}
	}

	go consumer(out1)
	go consumer(out2)

	wg.Wait()

	fmt.Printf("\nDone\n")

}
