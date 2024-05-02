package cacheupdater

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

// CacheUpdateReader stores Kafka topic names and let subscribe ones
type CacheUpdateReader struct {
	consumer       consumerOps
	validTopicsSet map[string]string
}

// NewCacheUpdateReader creates new CacheUpdateReader
func NewCacheUpdateReader(consumer consumerOps, validTopicsSet map[string]string) *CacheUpdateReader {
	return &CacheUpdateReader{
		consumer:       consumer,
		validTopicsSet: validTopicsSet,
	}
}

// Subscribe starts reading Kafka topic and writing it in simple int channel
func (r *CacheUpdateReader) Subscribe(ctx context.Context, topic string, output chan<- int64) error {
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
						log.Println("json.Unmarshal:", err)
						continue
					}
					output <- m.ID
				case <-ctx.Done():
					close(output)
					return
				}
			}
		}(pc)
	}

	return nil
}
