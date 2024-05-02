package kafkalogger

import (
	"time"
)

// Message is a struct that is marshalled to json and stored in the Kafka log topic
type Message struct {
	RequestTime time.Time `json:"creation_time"`
	HTTPMethod  string    `json:"http_method"`
	RawRequest  string    `json:"raw_request"`
}
