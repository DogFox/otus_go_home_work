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
	outputCh := make(chan error, maxErrorsCount)
	quit := make(chan bool)
	// вейтгруп для ожидания рутинок
	wg := &sync.WaitGroup{}
	// m := &sync.Mutex{}
	wg.Add(1)
	go func() {
		defer wg.Done()
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
				select {
				case <-quit:
					return
				default:
					for task := range inputCh {
						if res := task(); res != nil {
							outputCh <- res
							fmt.Println(len(outputCh))
						}
					}
				}
			}()
		}
		wg.Wait()
		close(outputCh)
	}()

	// собираем результаты
	//слайс результатов
	// output := make([]error, 0, len(tasks))
	// for res := range outputCh {
	// 	output = append(output, res)
	// 	// if len(output) >= maxErrorsCount {
	// 	// 	quit <- true
	// 	// }
	// }

	// fmt.Println(output)

	// Place your code here.
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
