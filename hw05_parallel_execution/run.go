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

func worker(wg *sync.WaitGroup, input <-chan Task, output chan<- error, quit chan bool) {
	defer wg.Done()

	select {
	case <-quit:
		close(output)
		return
	default:
		for task := range input {
			if res := task(); res != nil {
				output <- res
				fmt.Println(len(output))
			}
		}
	}
}

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, workersCount, maxErrorsCount int) error {
	/// канал с тасками
	inputCh := make(chan Task, len(tasks))
	// канал для сбора результатов
	outputCh := make(chan error, maxErrorsCount)
	quit := make(chan bool)
	// вейтгруп для ожидания рутинок
	wg := &sync.WaitGroup{}
	// m := &sync.Mutex{}

	for i := range tasks {
		inputCh <- tasks[i]
	}
	close(inputCh)

	for i := 0; i < workersCount; i++ {
		wg.Add(1)
		go worker(wg, inputCh, outputCh, quit)
	}

	output := 0
	for _ = range outputCh {
		output++
		if output >= maxErrorsCount {
			quit <- true
		}
	}
	fmt.Println(output)

	wg.Wait()

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

	workersCount := 10
	maxErrorsCount := 23
	Run(tasks, workersCount, maxErrorsCount)

}
