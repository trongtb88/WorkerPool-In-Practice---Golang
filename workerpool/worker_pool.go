package workerpool

import (
	"context"
	"sync"
)

type WorkerPool struct {
	ctx          context.Context
	maxGoRoutine int
	job          chan func()
	wg           *sync.WaitGroup // using it to wait all tasks finish
}

func NewWorkerPool(ctx context.Context, maxGoRoutine int) *WorkerPool {
	return &WorkerPool{
		ctx:          ctx,
		maxGoRoutine: maxGoRoutine,
		job:          make(chan func()),
		wg:           &sync.WaitGroup{},
	}
}

func (wp *WorkerPool) AddJob(JobWrapper func()) {
	wp.wg.Add(1)
	wp.job <- JobWrapper
}

func (wp *WorkerPool) CloseJob() {
	close(wp.job)
}

func (wp *WorkerPool) Run() {
	for i := 0; i < wp.maxGoRoutine; i++ {
		go func() {
			for job := range wp.job {
				// skip the rest job, if got ctx error
				if wp.ctx.Err() != nil {
					wp.wg.Done()
					continue
				}

				job()
				wp.wg.Done()
			}
		}()
	}
}

func (wp *WorkerPool) Wait() {
	wp.wg.Wait()
}
