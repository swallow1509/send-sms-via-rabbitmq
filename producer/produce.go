package main

import (
	"log"

	"github.com/streadway/amqp"
)

var (
	amqpURI  = "amqp://guest:guest@localhost:5672/"
	conn     *amqp.Connection
	channel  *amqp.Channel
	newQueue amqp.Queue
	err      error
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	conn, err = amqp.Dial(amqpURI)
	failOnError(err, "Failed to connect to RabbitMQ!")
	defer conn.Close()

	channel, err = conn.Channel()
	failOnError(err, "Failed to create channel")
	defer channel.Close()

	newQueue, err = channel.QueueDeclare(
		"newQueue",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")

	body := "Hello, AMQP!"
	err = channel.Publish(
		"",
		newQueue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		},
	)
	failOnError(err, "Failed to publish message")
}
