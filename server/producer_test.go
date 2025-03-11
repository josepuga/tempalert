package main

import (
	"testing"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// Integration Test: Produce Kafka Message

func TestKafkaProducer(t *testing.T) {

	// Create Producer
	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})
	if err != nil {
		t.Fatalf("Failed to initialize Kafka Producer: %v\n", err)
	}
	defer producer.Close()

	// Produce Message
	alertTopic := "sensor_alert"
	
    message := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &alertTopic,
			Partition: kafka.PartitionAny},
		Value: []byte("Test Alert from testKafkaProducer()"),
	}
    err = producer.Produce(message, nil)
    if err != nil {
        t.Errorf("Failed to produce message: %v", err)
    }
    time.Sleep(1 * time.Second) // Wait for Kafka to process the message
}
