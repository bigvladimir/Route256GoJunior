//go:generate mockgen -source ./producer.go -destination=./mocks/producer.go -package=mock_kafka
package kafka

import (
	"fmt"

	"github.com/IBM/sarama"
)

type Producer struct {
	brokers      []string
	syncProducer sarama.SyncProducer
}

func NewProducer(brokers []string) (*Producer, error) {
	syncProducer, err := newSyncProducer(brokers)
	if err != nil {
		return nil, fmt.Errorf("newSyncProducer: %w", err)
	}

	producer := &Producer{
		brokers:      brokers,
		syncProducer: syncProducer,
	}

	return producer, nil
}

func newSyncProducer(brokers []string) (sarama.SyncProducer, error) {
	syncProducerConfig := sarama.NewConfig()

	syncProducerConfig.Producer.Partitioner = sarama.NewRoundRobinPartitioner

	syncProducerConfig.Producer.RequiredAcks = sarama.WaitForAll

	syncProducerConfig.Producer.Idempotent = true
	syncProducerConfig.Net.MaxOpenRequests = 1

	syncProducerConfig.Producer.Return.Successes = true
	syncProducerConfig.Producer.Return.Errors = true

	syncProducer, err := sarama.NewSyncProducer(brokers, syncProducerConfig)
	if err != nil {
		return nil, fmt.Errorf("sarama.NewSyncProducer: %w", err)
	}

	return syncProducer, nil
}

func (k *Producer) SendSyncMessage(message *sarama.ProducerMessage) (partition int32, offset int64, err error) {
	return k.syncProducer.SendMessage(message)
}

func (k *Producer) SendSyncMessages(messages []*sarama.ProducerMessage) error {
	err := k.syncProducer.SendMessages(messages)
	if err != nil {
		fmt.Println("syncProducer.SendMessages:", err)
	}

	return err
}

func (k *Producer) Close() error {
	err := k.syncProducer.Close()
	if err != nil {
		return fmt.Errorf("syncProducer.Close: %w", err)
	}

	return nil
}
