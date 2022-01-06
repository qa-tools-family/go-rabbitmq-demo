package main

import (
	amqp "github.com/rabbitmq/amqp091-go"
	logger "github.com/sirupsen/logrus"
	"github.com/wagslane/go-rabbitmq"
)

func main() {
	publisher, err := rabbitmq.NewPublisher("amqp://user:password@localhost:5672/", amqp.Config{})
	if err != nil {
		logger.Fatal(err)
	}
	err = publisher.Publish([]byte("hello, world"), []string{"hello1"})
	if err != nil {
		logger.Fatal(err)
	}
}
