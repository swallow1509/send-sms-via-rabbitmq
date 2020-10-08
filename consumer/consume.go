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
		log.Fatalf("%s : %s", msg, err)
	}
}

func main() {
	conn, err = amqp.Dial(amqpURI)
	failOnError(err, "Failed to connect to RabbitMQ(consumer)")
	defer conn.Close()

	channel, err = conn.Channel()
	failOnError(err, "Failed to open a channel(consumer)")
	defer channel.Close()

	newQueue, err = channel.QueueDeclare(
		"newQueue",
		false,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue(consumer)")

	messages, err := channel.Consume(
		newQueue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range messages {
			log.Printf(" [x] Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
