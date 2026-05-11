package main

import (
	"fmt"
	"github.com/rabbitmq/amqp091-go"
	"time"
)

func main() {
	conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		fmt.Println(err)
		return
	}

	receiveMessage(ch, "test-queue")

}

func receiveMessage(ch *amqp091.Channel, queueName string) {
	msgs, err := ch.Consume(queueName, "", true, false, false, false, nil)
	if err != nil {
		return
	}

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			fmt.Println(string(d.Body))
		}
	}()

	go func() {
		time.Sleep(5 * time.Minute)
		forever <- true
	}()

	<-forever

}
