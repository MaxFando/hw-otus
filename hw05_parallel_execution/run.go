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

	var errorsCount int32
	errorsCh := make(chan error, maxErrorsCount)
	var wg sync.WaitGroup

	taskCh := tasksChannelGenerate(tasks)
	for i := 0; i < workersCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for task := range taskCh {
				if err := task(); err != nil {
					if atomic.LoadInt32(&errorsCount) >= int32(maxErrorsCount) {
						errorsCh <- ErrErrorsLimitExceeded
						return
					}

					atomic.AddInt32(&errorsCount, 1)
				}
			}
		}()
	}

	// Ждем завершения всех ворекров
	wg.Wait()
	close(errorsCh)
	err := <-errorsCh
	if err != nil {
		return err
	}
	return nil
}

func tasksChannelGenerate(tasks []Task) <-chan Task {
	tasksCh := make(chan Task, len(tasks))

	for _, task := range tasks {
		tasksCh <- task
	}

	close(tasksCh)

	return tasksCh
}
