package main

import (
	"fmt"
	"time"
)

func main() {

	ch := make(chan int)

	go func() {
		for {
			select {
			case ch <- 1:
				fmt.Println("sent 1")
			}
		}
	}()

	<-time.After(5 * time.Second)
	fmt.Println(<-ch)
	<-time.After(1 * time.Second)

}
