package main

import (
	"fmt"
	"sync"
	"time"
)

func generator(out chan<- int) {
	for i := 1; i <= 10; i++ {
		out <- i
	}
	close(out)
}

func squarer(in <-chan int, workers int) <-chan int {
	out := make(chan int)
	var wg sync.WaitGroup

	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for v := range in {
				time.Sleep(1 * time.Second)
				out <- v * v
			}
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func main() {
	numbersChan := make(chan int)
	go generator(numbersChan)

	squaresChan := squarer(numbersChan, 3)

	for result := range squaresChan {
		fmt.Println("get square:", result)
	}
}
