package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/qa-tools-family/go-rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	publisher, err := rabbitmq.NewPublisher(
		"amqp://user:password@localhost:5672/", amqp.Config{},
		rabbitmq.WithPublisherOptionsLogging,
	)
	if err != nil {
		log.Fatal(err)
	}

	// 接收消息返回信息
	returns := publisher.NotifyReturn()
	go func() {
		for r := range returns {
			log.Printf("message returned from server: %s", string(r.Body))
		}
	}()

	// 接收消息发送完成信息
	confirmations := publisher.NotifyPublish()
	go func() {
		for c := range confirmations {
			log.Printf("message confirmed from server. tag: %v, ack: %v", c.DeliveryTag, c.Ack)
		}
	}()

	// block main thread - wait for shutdown signal
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()

	fmt.Println("awaiting signal")

	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ticker.C:
			for i:=1; i<5; i++ {
				err = publisher.Publish(
					[]byte("hello, world"),
					[]string{"routing_key"},
					rabbitmq.WithPublishOptionsContentType("application/json"),
					rabbitmq.WithPublishOptionsMandatory,
					rabbitmq.WithPublishOptionsPersistentDelivery,
					rabbitmq.WithPublishOptionsExchange("events"),
				)
				if err != nil {
					log.Println(err)
				}
			}
		case <-done:
			fmt.Println("stopping publisher")
			err := publisher.StopPublishing()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("publisher stopped")
			return
		}
	}
}