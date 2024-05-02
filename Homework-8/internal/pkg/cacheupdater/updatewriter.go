// Package cacheupdater во многом копипаста логгера, если и дальше это масштабировать то надо писать что-то более универсальное
package cacheupdater

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/IBM/sarama"
)

type producerOps interface {
	SendSyncMessage(message *sarama.ProducerMessage) (partition int32, offset int64, err error)
	Close() error
}

// CacheUpdateWriter stores one Kafka topic name and let to send messages to it
type CacheUpdateWriter struct {
	producer producerOps
	topic    string
}

// NewCacheUpdateWriter creates new CacheUpdateWriter
func NewCacheUpdateWriter(producer producerOps, topic string) *CacheUpdateWriter {
	return &CacheUpdateWriter{
		producer: producer,
		topic:    topic,
	}
}

// SendMessage sends a message to Kafka
func (s *CacheUpdateWriter) SendMessage(message Message) error {
	if err := validateMessage(message); err != nil {
		return fmt.Errorf("kafkalogger.validateMessage: %w", err)
	}
	kafkaMsg, err := s.buildMessage(message)
	if err != nil {
		return fmt.Errorf("kafkalogger.buildMessage: %w", err)
	}

	_, _, err = s.producer.SendSyncMessage(kafkaMsg)
	if err != nil {
		return fmt.Errorf("producer.SendSyncMessage: %w", err)
	}

	return nil
}

func (s *CacheUpdateWriter) buildMessage(message Message) (*sarama.ProducerMessage, error) {
	msg, err := json.Marshal(message)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal: %w", err)
	}

	return &sarama.ProducerMessage{
		Topic:     s.topic,
		Value:     sarama.ByteEncoder(msg),
		Partition: -1,
	}, nil
}

func validateMessage(message Message) error {
	if message.ID <= 0 {
		return errors.New("Incorrect ID")
	}

	return nil
}
