package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	var ch0 chan int

	//curiously, although the select does not use the context
	//it does not break until the context is cancelled/timeout
	select {
	case ch0 <- 1:
		//case <-ctx.Done():
		//	fmt.Println("yay!")
	}

	fmt.Println(ctx)
}
