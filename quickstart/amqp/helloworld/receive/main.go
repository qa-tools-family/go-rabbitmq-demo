package main

import (
	amqp "github.com/rabbitmq/amqp091-go"
	logger "github.com/sirupsen/logrus"
	"log"
)

func main() {
	conn, err := amqp.Dial("amqp://user:password@localhost:5672/")
	if err != nil {
		logger.Fatal("connection to rabbitmq failed: ", err)
	}
	defer conn.Close()

	channel, err := conn.Channel()
	if err != nil {
		logger.Fatal("get rabbitmq channel failed: ", err)
	}
	defer channel.Close()

	// 确保 queue 存在
	queue, err := channel.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		logger.Fatal("declare queue failed: ", err)
	}

	msgs, err := channel.Consume(
		queue.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		logger.Fatal("Failed to register a consumer: ", err)
	}

	// 关于 channel 的使用需要学习一下
	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	// 尝试不断从 channel 中读取数据
	<-forever
}
