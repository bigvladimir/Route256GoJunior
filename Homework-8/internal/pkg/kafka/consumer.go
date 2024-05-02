package kafka

//go:generate mockgen -source ./consumer.go -destination=./mocks/consumer.go -package=mock_kafka

import (
	"fmt"
	"time"

	"github.com/IBM/sarama"
)

// Consumer is a receiver part of kafka
type Consumer struct {
	brokers        []string
	singleConsumer sarama.Consumer
}

// NewConsumer creates Consumer
func NewConsumer(brokers []string) (*Consumer, error) {
	config := sarama.NewConfig()

	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.AutoCommit.Enable = true
	config.Consumer.Offsets.AutoCommit.Interval = 5 * time.Second

	config.Consumer.Offsets.Initial = sarama.OffsetNewest

	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("sarama.NewConsumer: %w", err)
	}

	return &Consumer{
		brokers:        brokers,
		singleConsumer: consumer,
	}, err
}

// Partitions returns topic partitions
func (k *Consumer) Partitions(topic string) ([]int32, error) {
	return k.singleConsumer.Partitions(topic)
}

// ConsumePartition creates consumer for specific partition
func (k *Consumer) ConsumePartition(topic string, partition int32, offset int64) (sarama.PartitionConsumer, error) {
	return k.singleConsumer.ConsumePartition(topic, partition, offset)
}

// Close closes Consumer session
func (k *Consumer) Close() error {
	err := k.singleConsumer.Close()
	if err != nil {
		return fmt.Errorf("singleConsumer.Close: %w", err)
	}

	return nil
}
