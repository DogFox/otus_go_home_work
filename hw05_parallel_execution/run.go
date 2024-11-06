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

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, workersCount, maxErrorsCount int) error {
	/// канал с тасками
	inputCh := make(chan Task)
	// канал для сбора результатов
	outputCh := make(chan error)

	// вейтгруп для ожидания рутинок
	wg := &sync.WaitGroup{}

	//слайс результатов
	output := make([]error, 0, len(tasks))

	///
	go func() {
		defer close(inputCh)

		for i := range tasks {
			inputCh <- tasks[i]
		}
	}()

	go func() {
		for i := 0; i < workersCount; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for task := range inputCh {
					res := task()
					if (res) != nil {
						outputCh <- task()
					}
				}
			}()

		}
		wg.Wait()
		close(outputCh)
	}()

	// собираем результаты
	for res := range outputCh {
		output = append(output, res)
	}

	fmt.Println(output)

	// Place your code here.
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
