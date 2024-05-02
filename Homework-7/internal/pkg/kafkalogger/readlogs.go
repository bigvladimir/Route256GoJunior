package kafkalogger

//go:generate mockgen -source ./readlogs.go -destination=./mocks/readlogs.go -package=mock_logger

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/IBM/sarama"
)

type consumerOps interface {
	Partitions(topic string) ([]int32, error)
	ConsumePartition(topic string, partition int32, offset int64) (sarama.PartitionConsumer, error)
}

// LogWatcher stores Kafka topic names and let subscribe ones
type LogWatcher struct {
	consumer       consumerOps
	validTopicsSet map[string]string
}

// NewLogWatcher creates LogWatcher struct
func NewLogWatcher(consumer consumerOps, validTopicsSet map[string]string) *LogWatcher {
	return &LogWatcher{
		consumer:       consumer,
		validTopicsSet: validTopicsSet,
	}
}

// Subscribe starts reading Kafka topic and writing it in default log
func (r *LogWatcher) Subscribe(ctx context.Context, topic string) error {
	_, ok := r.validTopicsSet[topic]
	if !ok {
		return errors.New("Can not find topic")
	}

	partitionList, err := r.consumer.Partitions(topic)
	if err != nil {
		return fmt.Errorf("consumer.Partitions: %w", err)
	}
	initialOffset := sarama.OffsetNewest

	for _, partition := range partitionList {
		pc, err := r.consumer.ConsumePartition(topic, partition, initialOffset)
		if err != nil {
			return fmt.Errorf("consumer.ConsumePartition: %w", err)
		}

		go func(pc sarama.PartitionConsumer) {
			for {
				select {
				case message := <-pc.Messages():
					m := Message{}
					err := json.Unmarshal(message.Value, &m)
					if err != nil {
						log.Println("Can not unmarshal log:", err, "\n Row log:", message.Value)
						continue
					}
					log.Printf(
						"Request time: %v\nMethod: %s\nRaw request:\n%s", m.RequestTime, m.HTTPMethod, m.RawRequest,
					)
				case <-ctx.Done():
					return
				}
			}
		}(pc)
	}

	return nil
}
