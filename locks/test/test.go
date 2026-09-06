package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"
)

func main() {
	tq := testSharedQueue{}
	http.Handle("/test", &tq)

	_ = http.ListenAndServe(":8089", nil)
}

type testSharedQueue struct{}

func (q *testSharedQueue) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	wg := &sync.WaitGroup{}
	wg.Add(3)

	for i := 0; i < 3; i++ {
		go worker(r.Context(), i, wg, w)
	}

	wg.Wait()
}

func worker(pctx context.Context, i int, wg *sync.WaitGroup, w http.ResponseWriter) {
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:    2,
			IdleConnTimeout: 5 * time.Second,
		},
	}

	portToInt, _ := strconv.Atoi("8080")
	port := portToInt + i + 1
	for j := 0; j < 100; j++ {
		// ctx, _ := context.WithTimeout(pctx, 50*time.Second)
		req, _ := http.NewRequestWithContext(pctx, http.MethodGet, fmt.Sprintf("http://localhost:%d/get", port), nil)
		resp, err := client.Do(req)

		if resp != nil {
			respBytes, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			w.Write([]byte(fmt.Sprintf("good worker: %d, response: %v for server: %d\n", i, string(respBytes), port)))
		} else {
			w.Write([]byte(fmt.Sprintf("worker: %d, response: %v for server: %d\n", i, err, port)))
		}

	}
	wg.Done()
}
