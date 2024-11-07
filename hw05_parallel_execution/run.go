package main

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func worker(wg *sync.WaitGroup, input <-chan Task, output chan<- error, quit chan bool) {
	defer wg.Done()

	for {
		select {
		case <-quit:
			return
		case task, ok := <-input:
			if !ok {
				return
			}
			if res := task(); res != nil {
				output <- res
			}
		}
	}
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, workersCount, maxErrorsCount int) error {
	/// канал с тасками
	inputCh := make(chan Task)
	defer close(inputCh)
	// канал для сбора результатов
	outputCh := make(chan error, len(tasks))
	defer close(outputCh)
	// канал сигнал
	quit := make(chan bool, workersCount)
	defer close(quit)

	// вейтгруп для ожидания рутинок
	wg := &sync.WaitGroup{}

	// сначала подготавливаем пул рабочих
	for i := 0; i < workersCount; i++ {
		wg.Add(1)
		go worker(wg, inputCh, outputCh, quit)
	}

	// скармливаем задачи, проще на берегу не запускать чем потом чтото проверять
	maxErrorReached := false
	for i := range tasks {
		if len(outputCh) >= maxErrorsCount && !maxErrorReached {
			maxErrorReached = true
			break
		}
		inputCh <- tasks[i]
	}

	// хочется чтобы сигнальный канал был как флаг
	for i := 0; i < workersCount; i++ {
		quit <- true
	}

	wg.Wait()

	// fmt.Println(len(outputCh))
	if maxErrorReached {
		return ErrErrorsLimitExceeded
	}

	return nil
}
