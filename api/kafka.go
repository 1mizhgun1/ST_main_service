package api

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/IBM/sarama"
)

const (
	kafkaAddr       = "localhost:9092"
	kafkaTopic      = "segments"
	kafkaReadPeriod = 2 * time.Second
)

func putSegmentToKafka(segment codeRequest) error {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer([]string{kafkaAddr}, config)
	if err != nil {
		return fmt.Errorf("error creating producer: %w", err)
	}
	defer producer.Close()

	segmentString, _ := json.Marshal(segment)
	message := &sarama.ProducerMessage{
		Topic: kafkaTopic,
		Value: sarama.StringEncoder(segmentString),
	}

	_, _, err = producer.SendMessage(message)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

func readFromKafka() error {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer([]string{kafkaAddr}, config)
	if err != nil {
		return fmt.Errorf("error creating consumer: %w", err)
	}
	defer consumer.Close()

	partitionConsumer, err := consumer.ConsumePartition(kafkaTopic, 0, sarama.OffsetNewest)
	if err != nil {
		return fmt.Errorf("error opening topic: %w", err)
	}
	defer partitionConsumer.Close()

	for {
		select {
		case message := <-partitionConsumer.Messages():
			messageData := codeRequest{}
			err := json.Unmarshal(message.Value, &messageData)
			if err != nil {
				fmt.Printf("Error reading from kafka: %v", err)
			}
			addSegment(messageData)
		case err := <-partitionConsumer.Errors():
			fmt.Printf("Error: %s\n", err.Error())
		}
	}
}
