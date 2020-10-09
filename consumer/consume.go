package main

import (
	"bytes"
	"log"
	"time"

	"github.com/streadway/amqp"
)

var (
	amqpURI   = "amqp://guest:guest@localhost:5672/"
	conn      *amqp.Connection
	channel   *amqp.Channel
	taskQueue amqp.Queue
	err       error
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

	taskQueue, err = channel.QueueDeclare(
		"task", //name
		true,   //durable
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")

	err = channel.Qos(
		1, //fetch count
		0, //preFetch size
		false,
	)
	failOnError(err, "Failed to set Qos!")

	messages, err := channel.Consume(
		taskQueue.Name, //queue
		"",             // consumer
		false,          // auto-ack
		false,          // exclusive
		false,          // no-local
		false,          // no-wait
		nil,            // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range messages {
			log.Printf(" [x] Received a message: %s", d.Body)
			dotCount := bytes.Count(d.Body, []byte("."))
			t := time.Duration(dotCount)
			time.Sleep(t * time.Second)
			log.Printf("Done")
			d.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
