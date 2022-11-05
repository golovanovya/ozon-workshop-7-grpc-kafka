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

	producer, err := newSyncProducer(BrokersList)
	if err != nil {
		log.Fatal(err)
	}

	rand.Seed(time.Now().Unix())

	// Inject info into message
	msg := sarama.ProducerMessage{
		Topic: KafkaTopic,
		Key:   sarama.StringEncoder(fmt.Sprintf("%d", rand.Intn(1000000))),
		Value: sarama.StringEncoder(fmt.Sprintf("%d", rand.Intn(1000))),
	}

	p, o, err := producer.SendMessage(&msg)
	if err != nil {
		log.Fatalln("Failed to send message:, err")
	}

	log.Printf("Successful to write message, topic %s, offset:%d, partition: %d\n", KafkaTopic, o, p)

	err = producer.Close()
	if err != nil {
		log.Fatalln("Failed to close producer:", err)
	}
}

func newSyncProducer(brokerList []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Version = sarama.V2_8_0_0
	// Waits for all in-sync replicas to commit before responding.
	config.Producer.RequiredAcks = sarama.WaitForAll
	// The total number of times to retry sending a message (default 3).
	config.Producer.Retry.Max = 3
	// How long to wait for the cluster to settle between retries (default 100ms).
	config.Producer.Retry.Backoff = time.Millisecond * 250
	// idempotent producer has a unique producer ID and uses sequence IDs for each message,
	// allowing the broker to ensure, on a per-partition basis, that it is committing ordered messages with no duplication.
	//config.Producer.Idempotent = true
	if config.Producer.Idempotent == true {
		config.Producer.Retry.Max = 1
		config.Net.MaxOpenRequests = 1
	}
	//  Successfully delivered messages will be returned on the Successes channe
	config.Producer.Return.Successes = true
	// Generates partitioners for choosing the partition to send messages to (defaults to hashing the message key)
	_ = config.Producer.Partitioner

	producer, err := sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		return nil, fmt.Errorf("starting Sarama producer: %w", err)
	}

	return producer, nil
}
