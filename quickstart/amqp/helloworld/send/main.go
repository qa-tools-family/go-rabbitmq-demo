package main


import (
	amqp "github.com/rabbitmq/amqp091-go"
	logger "github.com/sirupsen/logrus"
)

func main()  {
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

	body := "Hello World!"
	for true {
		err = channel.Publish(
			"",         // exchange
			queue.Name, // routing key
			false,      // mandatory
			false,      // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			})
		if err != nil {
			logger.Fatal("Failed to publish a message: ", err)
		}
	}

}

