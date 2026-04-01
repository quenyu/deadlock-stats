package main

import (
	"fmt"
	"sync"
)

func main() {
	ch1 := make(chan int)
	ch2 := make(chan int)

	go func() {
		for i := 0; i < 5; i++ {
			ch1 <- i
		}
		close(ch1)
	}()

	go func() {
		for i := 5; i < 10; i++ {
			ch2 <- i
		}
		close(ch2)
	}()

	for v := range mergeTwoChannels(ch1, ch2) {
		fmt.Println(v)
	}
}

func mergeTwoChannels(ch1 <-chan int, ch2 <-chan int) <-chan int {
	res := make(chan int)
	wg := sync.WaitGroup{}

	output := func(c <-chan int) {
		for val := range c {
			res <- val
		}
		wg.Done()
	}

	wg.Add(2)
	go output(ch1)
	go output(ch2)

	go func() {
		wg.Wait()
		close(res)
	}()

	return res
}
