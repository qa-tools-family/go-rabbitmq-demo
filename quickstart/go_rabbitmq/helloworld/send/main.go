package main

import (
	amqp "github.com/rabbitmq/amqp091-go"
	logger "github.com/sirupsen/logrus"
	"github.com/wagslane/go-rabbitmq"
	"log"
)

type customLogger struct{}

// Printf is the only method needed in the Logger interface to function properly.
func (c *customLogger) Printf(fmt string, args ...interface{}) {
	log.Printf("mylogger: " + fmt, args...)
}

func main() {
	mylogger := &customLogger{}
	publisher, err := rabbitmq.NewPublisher("amqp://user:password@localhost:5672/", amqp.Config{},
		rabbitmq.WithPublisherOptionsLogger(mylogger))

	if err != nil {
		logger.Fatal(err)
	}
	err = publisher.Publish([]byte("hello, world"), []string{"hello1"})
	if err != nil {
		logger.Fatal(err)
	}
}
