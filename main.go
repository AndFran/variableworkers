package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func worker(workChannel chan int, result chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		i, ok := <-workChannel
		// fake long processes
		time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
		if !ok {
			return
		}
		i *= i
		result <- i
	}
}

func producer(workChannel chan int) {
	// can vary the amount of work
	for i := 0; i < 400; i++ {
		workChannel <- i
	}
	close(workChannel)
}

func main() {
	workChannel := make(chan int)
	resultChannel := make(chan int)
	done := make(chan struct{})

	var wg sync.WaitGroup

	go producer(workChannel)

	// can vary the number of workers
	for i := 0; i < 200; i++ {
		wg.Add(1)
		go worker(workChannel, resultChannel, &wg)
	}

	go func() {
		wg.Wait()
		close(done)
	}()

	counter := 0

	for {
		select {
		case i := <-resultChannel:
			counter++
			fmt.Println(i)
		case <-done:
			fmt.Println("Done")
			goto loop
		}
	}
loop:
	fmt.Println("Processed ", counter, " items")
}
