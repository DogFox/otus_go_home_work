package main

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func worker(wg *sync.WaitGroup, input <-chan Task, output chan<- error, quit chan bool, maxErrorsCount int) {
	defer wg.Done()

	for {
		select {
		case <-quit:
			return
		case task, ok := <-input:
			if !ok {
				return
			}
			// если в канале уже полна коробочка - можно не выполнять
			if len(output) == maxErrorsCount {
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
	inputCh := make(chan Task, len(tasks))
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
		go worker(wg, inputCh, outputCh, quit, maxErrorsCount)
	}

	// скармливаем задачи
	for i := range tasks {
		inputCh <- tasks[i]
	}

	maxErrorReached := false
	for {
		if len(outputCh) >= maxErrorsCount && !maxErrorReached {
			maxErrorReached = true
			break
		}
	}

	// хочется чтобы сигнальный канал был как флаг
	for i := 0; i < workersCount; i++ {
		quit <- true
	}

	wg.Wait()

	fmt.Println(len(outputCh))
	if maxErrorReached {
		return ErrErrorsLimitExceeded
	}

	return nil
}

func main() {
	tasksCount := 50
	tasks := make([]Task, 0, tasksCount)

	var runTasksCount int32

	for i := 0; i < tasksCount; i++ {
		err := fmt.Errorf("error from task %d", i)
		tasks = append(tasks, func() error {
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
			atomic.AddInt32(&runTasksCount, 1)
			return err
		})
	}

	workersCount := 2
	maxErrorsCount := 30
	Run(tasks, workersCount, maxErrorsCount)

}
