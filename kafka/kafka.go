package kafka

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
	"redirect-service/config"
)

func PublishMessage(id string, url string) {
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{config.KafkaBroker},
		Topic:    config.KafkaTopic,
		Balancer: &kafka.LeastBytes{},
	})

	msg := kafka.Message{
		Key:   []byte(id),
		Value: []byte(url),
	}

	err := w.WriteMessages(context.Background(), msg)
	if err != nil {
		log.Printf("Failed to send message to Kafka: %v\n", err)
	}
	log.Printf("Sent Message on kafka")
	w.Close()
}
