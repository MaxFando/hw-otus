package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, workersCount, maxErrorsCount int) error {
	if maxErrorsCount <= 0 {
		return ErrErrorsLimitExceeded
	}

	tasksCh := make(chan Task)
	var errorsCount int64
	var wg sync.WaitGroup

	for i := 0; i < workersCount; i++ {
		wg.Add(1)
		go worker(&wg, tasksCh, &errorsCount, maxErrorsCount)
	}

	wg.Add(1)
	go generator(&wg, tasks, tasksCh, &errorsCount, maxErrorsCount)

	wg.Wait()

	if errorsCount >= int64(maxErrorsCount) {
		return ErrErrorsLimitExceeded
	}

	return nil
}

func worker(wg *sync.WaitGroup, tasksCh chan Task, errorsCount *int64, maxErrorsCount int) {
	defer wg.Done()

	for task := range tasksCh {
		if err := task(); err != nil {
			if atomic.LoadInt64(errorsCount) >= int64(maxErrorsCount) {
				return
			}

			atomic.AddInt64(errorsCount, 1)
		}
	}
}

func generator(wg *sync.WaitGroup, tasks []Task, tasksCh chan Task, errorsCount *int64, maxErrorsCount int) {
	defer wg.Done()
	defer close(tasksCh)

	for _, task := range tasks {
		if atomic.LoadInt64(errorsCount) >= int64(maxErrorsCount) {
			return
		}
		tasksCh <- task
	}
}
