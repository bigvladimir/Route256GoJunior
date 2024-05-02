package cacheupdater

// Message is a struct that is marshalled to json and stored in the Kafka topic
type Message struct {
	ID int64 `json:"id"`
}
