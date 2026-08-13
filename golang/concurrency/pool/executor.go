package main

import (
	"fmt"
	"sync"
)

func main() {

	cp := NewConnectionPool()
	wg := &sync.WaitGroup{}
	wg.Add(20)

	for i := 0; i < 20; i++ {
		go func(i int) {
			defer wg.Done()
			if i == 6 {
				cp.Close()
			}
			conn, err := cp.Acquire()
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(conn.Display)
			fmt.Printf("executing task %d\n", i)

		}(i)
	}

	wg.Wait()
}
