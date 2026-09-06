package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func main() {
	sq := sharedQueue{Message: make([]string, 0)}
	sq.AddMessage()
	http.Handle("/read", &sq)

	_ = http.ListenAndServe(":8080", nil)
}

type sharedQueue struct {
	Message         []string
	currentPosition int
}

func (q *sharedQueue) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	msg := q.Message[0]
	q.Message = q.Message[1:]
	w.WriteHeader(200)
	w.Write([]byte(msg))
}

func (q *sharedQueue) AddMessage() {

	for i := 0; i < 10000; i++ {
		q.Message = append(q.Message, fmt.Sprintf("message: %s", strconv.Itoa(i)))
	}

}
