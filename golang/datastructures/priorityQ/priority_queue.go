package main

import (
	"fmt"
	"math"
	"math/rand"
)

type PriorityQueue struct {
	messages []message
}

type message struct {
	priority int
	content  string
}

func (pq *PriorityQueue) Publish(priority int, content string) {
	msg := message{priority, content}
	pq.messages = append(pq.messages, msg)
	heapify(len(pq.messages)-1, pq.messages)
}

func (pq *PriorityQueue) Subscribe() string {
	if len(pq.messages) == 0 {
		return ""
	}
	msg := pq.messages[0]
	pq.messages = pq.messages[1:]
	heapify(len(pq.messages)-1, pq.messages)
	return msg.content
}

func heapify(index int, messages []message) {
	if index <= 0 {
		return
	}

	parent := (index - 1) / 2

	if messages[parent].priority < messages[index].priority {
		parentMsg := messages[parent]
		childMsg := messages[index]

		messages[index] = parentMsg
		messages[parent] = childMsg
		heapify(parent, messages)

	} else {
		return
	}
}

func NewPriorityQueue(size int) PriorityQueue {
	return PriorityQueue{
		messages: make([]message, 0, size),
	}
}

func main() {
	pq := NewPriorityQueue(10)
	for i := 0; i < 10; i++ {
		pq.Publish(rand.Int(), fmt.Sprintf("messgeContent for %d", i))
	}

	pq.Publish(math.MaxInt, "messge Content for max")
	for i := 0; i < len(pq.messages); i++ {
		fmt.Printf("%d::%s\n", pq.messages[i].priority, pq.messages[i].content)
	}

}
