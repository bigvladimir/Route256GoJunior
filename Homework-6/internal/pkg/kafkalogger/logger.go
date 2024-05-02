//go:generate mockgen -source ./logger.go -destination=./mocks/logger.go -package=mock_logger
package kafkalogger

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/IBM/sarama"
)

type producerOps interface {
	SendSyncMessage(message *sarama.ProducerMessage) (partition int32, offset int64, err error)
	SendSyncMessages(messages []*sarama.ProducerMessage) error
	Close() error
}

// KafkaLogger stores one Kafka topic name and let to send messages to it
type KafkaLogger struct {
	producer producerOps
	topic    string
}

// NewKafkaLogger creates KafkaLogger struct
func NewKafkaLogger(producer producerOps, topic string) *KafkaLogger {
	return &KafkaLogger{
		producer,
		topic,
	}
}

// LogMessage sends a message to Kafka
func (s *KafkaLogger) LogMessage(message Message) error {
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

// LogMessage sends several messages to Kafka
func (s *KafkaLogger) LogMessages(messages []Message) error {
	var kafkaMsg []*sarama.ProducerMessage
	var ms *sarama.ProducerMessage
	var err error

	for _, m := range messages {
		if err = validateMessage(m); err != nil {
			return fmt.Errorf("kafkalogger.validateMessage: %w", err)
		}
		ms, err = s.buildMessage(m)
		if err != nil {
			return fmt.Errorf("kafkalogger.buildMessage: %w", err)
		}

		kafkaMsg = append(kafkaMsg, ms)
	}

	if len(messages) > 0 {
		err = s.producer.SendSyncMessages(kafkaMsg)
		if err != nil {
			return fmt.Errorf("producer.SendSyncMessages: %w", err)
		}
	}

	return nil
}

func (s *KafkaLogger) buildMessage(message Message) (*sarama.ProducerMessage, error) {
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
	if message.RequestTime.After(time.Now()) {
		return errors.New("Message from the future")
	}
	if message.HTTPMethod == "" {
		return errors.New("Empty HTTP method field")
	}
	if message.RawRequest == "" {
		return errors.New("Empty raw request field")
	}
	return nil
}
