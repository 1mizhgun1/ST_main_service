package api

import (
	"fmt"
	"log"

	"github.com/IBM/sarama"
)

const (
	kafkaAddr  = "localhost:9092"
	kafkaTopic = "segments"
)

func putSegmentToKafka(segment string) {
	config := sarama.NewConfig()

	producer, err := sarama.NewSyncProducer([]string{kafkaAddr}, config)
	if err != nil {
		log.Fatalf("Error creating producer: %s", err)
	}
	defer func() {
		if err := producer.Close(); err != nil {
			log.Fatalf("Error closing producer: %s", err)
		}
	}()

	message := &sarama.ProducerMessage{
		Topic: kafkaTopic,
		Value: sarama.StringEncoder("Your message data here"),
	}

	partition, offset, err := producer.SendMessage(message)
	if err != nil {
		log.Fatalf("Failed to send message: %s", err)
	}
	fmt.Printf("Message sent to partition %d at offset %d\n", partition, offset)
}
