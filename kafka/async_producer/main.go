package main

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"github.com/Shopify/sarama"
)

var (
	KafkaTopic  = "example-topic"
	BrokersList = []string{"localhost:9092"}
)

func main() {
	log.Printf("Kafka brokers: %s", strings.Join(BrokersList, ", "))

	producer, err := newProducer(BrokersList)
	if err != nil {
		log.Fatal(err)
	}

	rand.Seed(time.Now().Unix())

	// Inject info into message
	msg := sarama.ProducerMessage{
		Topic:   KafkaTopic,
		Key:     sarama.StringEncoder("random_number"),
		Value:   sarama.StringEncoder(fmt.Sprintf("%d", rand.Intn(1000))),
		Headers: []sarama.RecordHeader{{Key: []byte("SomeKey"), Value: []byte("SomeValue")}},
	}

	producer.Input() <- &msg
	successMsg := <-producer.Successes()
	log.Println("Successful to write message, offset:", successMsg.Offset)

	err = producer.Close()
	if err != nil {
		log.Fatalln("Failed to close producer:", err)
	}
}

func newProducer(brokerList []string) (sarama.AsyncProducer, error) {
	config := sarama.NewConfig()
	config.Version = sarama.V2_5_0_0
	// So we can know the partition and offset of messages.
	config.Producer.Return.Successes = true

	producer, err := sarama.NewAsyncProducer(brokerList, config)
	if err != nil {
		return nil, fmt.Errorf("starting Sarama producer: %w", err)
	}

	// We will log to STDOUT if we're not able to produce messages.
	go func() {
		for err := range producer.Errors() {
			log.Println("Failed to write message:", err)
		}
	}()

	return producer, nil
}
