package main

import (
	"fmt"
	"sync"
	"time"
)

type result struct {
	url string
	name string
}

func worker(n int, done <-chan struct{}, jobs <-chan int, c chan<- string) {
	for j := range jobs {
		select {
		case c <- fmt.Sprintf("Worker: %d, Job: %d", n, j):
		case <-done:
			return
		}
	}
}

func main() {
	jobs := make(chan int)
	done := make(chan struct{})
	c := make(chan string)

	var wg sync.WaitGroup
	const numWorkers = 5
	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go func(n int) {
			worker(n, done, jobs, c)
			wg.Done()
		}(i)
	}

	go func() {
		wg.Wait()
		close(c)
	}()

	go func() {
		for i := 0; i < 10; i++ {
			jobs <- i
			time.Sleep(10 * time.Millisecond)
		}

		close(done)
	
	}()

	for r := range c {
		fmt.Println(r)
	}
}
