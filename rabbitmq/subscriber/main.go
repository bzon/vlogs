package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/peterbourgon/ff"
	"github.com/streadway/amqp"
)

func main() {
	fs := flag.NewFlagSet("publisher", flag.ExitOnError)
	interval := fs.Duration("interval", 1*time.Second, "publishing interval")
	user := fs.String("user", "guest", "RabbitMQ user")
	pass := fs.String("pass", "guest", "RabbitMQ password")
	addr := fs.String("addr", "localhost:5672", "RabbitMQ address")
	node := fs.String("node", "/", "RabbitMQ Node")
	if err := ff.Parse(fs, os.Args[1:], ff.WithEnvVarPrefix("ENV")); err != nil {
		log.Fatal(err)
	}

	amqpAddr := fmt.Sprintf("amqp://%s:%s@%s%s", *user, *pass, *addr, *node)
	conn, err := amqp.Dial(amqpAddr)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when usused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			time.Sleep(*interval)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
