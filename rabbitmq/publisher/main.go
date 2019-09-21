package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/docker/docker/pkg/namesgenerator"
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
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	for {
		// body := "Hello World!"
		body := namesgenerator.GetRandomName(1)
		err = ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			})
		failOnError(err, "Failed to publish a message")

		log.Printf("Successfully published %s", body)

		// loop every second.
		time.Sleep(*interval)
	}

}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
