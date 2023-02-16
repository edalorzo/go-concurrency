package common

import "context"

func Producer(ctx context.Context, begin, end int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for i := begin; i < end; i++ {
			select {
			case out <- i:
			case <-ctx.Done():
				return
			}
		}
	}()
	return out
}
