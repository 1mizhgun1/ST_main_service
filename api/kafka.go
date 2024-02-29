package api

import (
	"fmt"
	"github.com/IBM/sarama"
)

const (
	kafkaAddr  = "localhost:9092"
	kafkaTopic = "segments"
)

func putSegmentToKafka(segment string) error {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer([]string{kafkaAddr}, config)
	if err != nil {
		return fmt.Errorf("error creating producer: %w", err)
	}
	defer producer.Close()

	message := &sarama.ProducerMessage{
		Topic: kafkaTopic,
		Value: sarama.StringEncoder(segment),
	}

	_, _, err = producer.SendMessage(message)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}
