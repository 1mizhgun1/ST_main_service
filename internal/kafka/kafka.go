package kafka

import (
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"

	"github.com/1mizhgun1/ST_main_service/internal/consts"
	"github.com/1mizhgun1/ST_main_service/internal/storage"
	"github.com/1mizhgun1/ST_main_service/internal/utils"
)

func PutSegmentToKafka(segment utils.CodeRequest) error {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer([]string{consts.KafkaAddr}, config)
	if err != nil {
		return fmt.Errorf("error creating producer: %w", err)
	}
	defer producer.Close()

	segmentString, _ := json.Marshal(segment)
	message := &sarama.ProducerMessage{
		Topic: consts.KafkaTopic,
		Value: sarama.StringEncoder(segmentString),
	}

	_, _, err = producer.SendMessage(message)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

func ReadFromKafka() error {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	consumer, err := sarama.NewConsumer([]string{consts.KafkaAddr}, config)
	if err != nil {
		return fmt.Errorf("error creating consumer: %w", err)
	}
	defer consumer.Close()

	partitionConsumer, err := consumer.ConsumePartition(consts.KafkaTopic, 0, sarama.OffsetNewest)
	if err != nil {
		return fmt.Errorf("error opening topic: %w", err)
	}
	defer partitionConsumer.Close()

	for {
		select {
		case message := <-partitionConsumer.Messages():
			messageData := utils.CodeRequest{}
			err := json.Unmarshal(message.Value, &messageData)
			if err != nil {
				fmt.Printf("Error reading from kafka: %v", err)
			}
			storage.AddSegment(messageData)
		case err := <-partitionConsumer.Errors():
			fmt.Printf("Error: %s\n", err.Error())
		}
	}
}
