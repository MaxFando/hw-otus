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
	var wg sync.WaitGroup
	sem := make(chan struct{}, workersCount) // Семофор для контроля кол-ва воркеров, работающих одновременно

	for _, task := range tasks {
		sem <- struct{}{} // Записваем в семофор
		wg.Add(1)

		go func(task Task) {
			defer wg.Done()
			defer func() { <-sem }() // Освобождаем семафор

			if err := task(); err != nil {
				if atomic.LoadInt32(&errorsCount) == int32(maxErrorsCount) {
					return
				}

				atomic.AddInt32(&errorsCount, 1)
			}
		}(task)

		// Если кол-во ошибок в канале равно максимальному кол-ву ошибок, то перестаем обрабатывать задачи
		if atomic.LoadInt32(&errorsCount) == int32(maxErrorsCount) {
			return ErrErrorsLimitExceeded
		}
	}

	// Ждем завершения всех ворекров
	wg.Wait()

	return nil
}
