package main

import (
	"fmt"
	"github.com/rabbitmq/amqp091-go"
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

	q, err := ch.QueueDeclare("test-queue", true, false, false, false, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = publishMessage(ch, q.Name, []byte("hello queue testing, this is persistant!!"))
	if err != nil {
		return
	}

}

func publishMessage(ch *amqp091.Channel, queueName string, body []byte) error {
	err := ch.Publish("", queueName, false, false, amqp091.Publishing{
		DeliveryMode: amqp091.Persistent,
		ContentType:  "text/plain",
		Body:         body,
	})
	if err != nil {
		return err
	}
	return nil
}
