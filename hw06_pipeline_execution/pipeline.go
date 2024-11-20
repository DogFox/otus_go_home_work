package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

// воркер получает первое значение и далее прогоняет все стейджы, вычитывая и записывая обратно в канал
func worker(value interface{}, done In, stages ...Stage) Out {
	//небуферезированнный блокирует сразу, так как нет читателя, поэтому тут буфер 1
	workerCh := make(Bi, 1)
	workerCh <- value
	defer close(workerCh)

	for _, stage := range stages {
		for {
			select {
			case value, ok := <-stage(workerCh):
				if !ok {
					return workerCh
				}
				workerCh <- value
				break
			case <-done:
				return workerCh
			}
		}
		// select {
		// case <-done:
		// 	return workerCh
		// case res := <-stage(workerCh):
		// 	fmt.Println(res)
		// 	workerCh <- res
		// }
	}

	return workerCh
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	wg := &sync.WaitGroup{}
	multiplexedStream := make(chan interface{})
	// multiplex := func(c <-chan interface{}) {
	// 	defer wg.Done()
	// 	for i := range c {
	// 		select {
	// 		case <-done:
	// 			return
	// 		case multiplexedStream <- i:
	// 		}
	// 	}
	// }

	for value := range in {
		wg.Add(1)
		// go multiplex(worker(value, done, stages...))
		pipeline := worker(value, done, stages...)

		for s := range pipeline {
			fmt.Println(s.(string))
		}
	}

	go func() {
		wg.Wait()
		close(multiplexedStream)
	}()

	return multiplexedStream
}

func main() {
	g := func(_ string, f func(v interface{}) interface{}) Stage {
		return func(in In) Out {
			out := make(Bi)
			go func() {
				defer close(out)
				for v := range in {
					time.Sleep(time.Millisecond * 100)
					out <- f(v)
				}
			}()
			return out
		}
	}

	stages := []Stage{
		g("Dummy", func(v interface{}) interface{} { return v }),
		g("Multiplier (* 2)", func(v interface{}) interface{} { return v.(int) * 2 }),
		g("Adder (+ 100)", func(v interface{}) interface{} { return v.(int) + 100 }),
		g("Stringifier", func(v interface{}) interface{} { return strconv.Itoa(v.(int)) }),
	}

	in := make(Bi)
	// data := []int{1, 2, 3, 4, 5}
	data := []int{1}

	go func() {
		for _, v := range data {
			in <- v
		}
		close(in)
	}()

	result := make([]string, 0, 10)
	start := time.Now()
	for s := range ExecutePipeline(in, nil, stages...) {
		result = append(result, s.(string))
	}
	elapsed := time.Since(start)

	fmt.Println(elapsed)
	fmt.Println(result)

}
