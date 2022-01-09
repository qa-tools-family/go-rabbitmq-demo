package main

import (
	amqp "github.com/rabbitmq/amqp091-go"
	logger "github.com/sirupsen/logrus"
	"github.com/qa-tools-family/go-rabbitmq"
)

func main()  {
	consumer, err := rabbitmq.NewConsumer("amqp://user:password@localhost:5672/", amqp.Config{})
	if err != nil {
		logger.Fatal(err)
	}
	// 声明队列 + 创建 Consume + 从 Consume 中读取数据
	err = consumer.StartConsuming(
		func(d rabbitmq.Delivery) rabbitmq.Action {
			logger.Printf("consumed: %v", string(d.Body))
			// rabbitmq.Ack, rabbitmq.NackDiscard, rabbitmq.NackRequeue
			return rabbitmq.Ack
		},
		"hello1",
		[]string{},
	)
	if err != nil {
		logger.Fatal(err)
	}

	// 阻塞等待
	forever := make(chan bool)
	<-forever
}