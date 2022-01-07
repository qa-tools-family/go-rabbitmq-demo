

package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/wagslane/go-rabbitmq"
)

var consumerName = "example"

func processor (d rabbitmq.Delivery) rabbitmq.Action {
	log.Printf("received: %v", string(d.Body))
	time.Sleep(10 * time.Second)
	log.Printf("consumed: %v", string(d.Body))
	// rabbitmq.Ack, rabbitmq.NackDiscard, rabbitmq.NackRequeue
	return rabbitmq.Ack
}

func main() {
	consumer, err := rabbitmq.NewConsumer(
		"amqp://user:password@localhost:5672/", amqp.Config{},
		rabbitmq.WithConsumerOptionsLogging,
	)
	if err != nil {
		log.Fatal(err)
	}
	err = consumer.StartConsuming(
		processor,
		"my_queue",
		[]string{"routing_key", "routing_key_2"},
		rabbitmq.WithConsumeOptionsConcurrency(10),
		//rabbitmq.WithConsumeOptionsQueueDurable,
		//rabbitmq.WithConsumeOptionsQuorum,
		rabbitmq.WithConsumeOptionsBindingExchangeName("events"),
		rabbitmq.WithConsumeOptionsBindingExchangeKind("topic"),
		rabbitmq.WithConsumeOptionsBindingExchangeDurable,
		rabbitmq.WithConsumeOptionsConsumerName(consumerName),
	)
	if err != nil {
		log.Fatal(err)
	}

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
	<-done
	fmt.Println("stopping consumer")

	// wait for server to acknowledge the cancel
	consumer.StopConsuming(consumerName, false)
	// todo: consumer 需要支持一个 Wait Task Finished 方法
	time.Sleep(time.Second * 20)
	consumer.Disconnect()
}