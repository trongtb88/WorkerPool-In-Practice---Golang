package main

import (
	"../workerpool"
	"context"
	"log"
)

func main() {

	// Start Worker Pool.
	totalWorker := 5
	wp := workerpool.NewWorkerPool(context.Background(), totalWorker)
	wp.Run()

	type result struct {
		id    int
		value int
	}

	totalTask := 10000
	resultC := make(chan result, totalTask)

	for i := 0; i < totalTask; i++ {
		id := i + 1
		wp.AddJob(func() {
			log.Printf("[main] Starting task %d", id)
			resultC <- result{id, id * 2}
		})
	}
	wp.Wait()

	for i := 0; i < totalTask; i++ {
		res := <-resultC
		log.Printf("[main] Task %d has been finished with result %d", res.id, res.value)
	}

}
