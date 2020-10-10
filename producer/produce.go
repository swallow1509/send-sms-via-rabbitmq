package main

import (
	"log"
	"os"
	"strings"

	"github.com/streadway/amqp"
)

var (
	amqpURI   = "amqp://guest:guest@localhost:5672/"
	conn      *amqp.Connection
	channel   *amqp.Channel
	taskQueue amqp.Queue
	err       error
)

//Operation usually fails if RabbitMQ server is not started
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

//Function for getting string from Terminal
func fromBody(str []string) string {
	var s string
	if (len(str) < 2) || os.Args[1] == "" {
		s = "Hello"
	} else {
		s = strings.Join(str[1:], " ")
	}
	return s
}

func main() {
	//Creating connection
	conn, err = amqp.Dial(amqpURI)
	failOnError(err, "Failed to connect to RabbitMQ!")
	defer conn.Close()

	//Opening channel on top of connection. Connections are expensive.
	//Channels are cheap. Multiple channels can be created on one channel.
	channel, err = conn.Channel()
	failOnError(err, "Failed to create channel")
	defer channel.Close()

	//Declaring a "taskQueue" with "task" name property
	taskQueue, err = channel.QueueDeclare(
		"task", //name
		true,   //durable
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")

	//Getting messages from Terminal and  print n times
	body := fromBody(os.Args)

	//Publish properties
	err = channel.Publish(
		"",             //exchange
		taskQueue.Name, //routing key
		false,          //mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(body),
		},
	)
	failOnError(err, "Failed to publish message")
}
