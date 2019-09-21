## Rabbitmq with Go tutorial

Install RabbitMQ. 

Access the admin page from http://localhost:15672.

User is `guest` and password is `guest`.

```console
$ docker-compose up -d
```

Run the publisher.

```console
$ go run ./publisher/main.go -interval 100ms
```

Run the subscriber.

```console
$ go run ./subscriber/main.go -interval 1s
```
