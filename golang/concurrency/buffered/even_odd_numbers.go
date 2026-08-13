package main

import (
	"fmt"
	"sync"
)

const maxNumbers = 100

type ExecuteFunc = func() (int, bool)
type PushFunc = func(x int)

var (
	even = make(chan int)
	odd  = make(chan int)
)

func main() {

	wg := &sync.WaitGroup{}
	wg.Add(3)

	go startExecution(wg)
	go printNumbers(getExecuteFunc(odd), getPushFunc(even), wg, "odd func")
	go printNumbers(getExecuteFunc(even), getPushFunc(odd), wg, "even func")

	wg.Wait()

}

func getExecuteFunc(num chan int) func() (int, bool) {
	return func() (int, bool) {
		x, ok := <-num
		return x, ok
	}
}

func getPushFunc(num chan int) func(x int) {
	return func(x int) {
		num <- x
	}
}

func startExecution(wg *sync.WaitGroup) {
	odd <- 1
	wg.Done()
}

func printNumbers(executeFunc ExecuteFunc, pushFunc PushFunc, wg *sync.WaitGroup, funcType string) {
	for {
		x, ok := executeFunc()
		if !ok {
			fmt.Println("channel closed")
			wg.Done()
			return
		}
		fmt.Printf("%s :: %d\n", funcType, x)
		next := x + 1
		if next > maxNumbers {
			close(even)
			close(odd)
			wg.Done()
			return
		}
		pushFunc(next)
	}
}
