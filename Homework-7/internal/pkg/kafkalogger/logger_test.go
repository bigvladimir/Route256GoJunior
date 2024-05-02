package kafkalogger

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_buildMessage(t *testing.T) {
	t.Parallel()
	t.Run("smoke", func(t *testing.T) {
		t.Parallel()

		kl := setUpKafkaLogger(t)
		defer kl.tearDown()
		inputMess := Message{
			RequestTime: time.Now(),
			HTTPMethod:  "POST",
			RawRequest:  "raw request text",
		}
		JSONmess, err := json.Marshal(inputMess)
		require.NoError(t, err)
		wantedSaramaMess := sarama.ProducerMessage{
			Topic:     testTopic,
			Value:     sarama.ByteEncoder(JSONmess),
			Partition: -1,
		}

		actualMess, err := kl.kafkaLogger.buildMessage(inputMess)

		require.NoError(t, err)
		assert.Equal(t, wantedSaramaMess, *actualMess)
	})
}

func Test_validateMessage(t *testing.T) {
	t.Parallel()
	testTime := time.Now()

	tt := []struct {
		name         string
		inputMessage Message
		wantedErr    error
	}{
		{
			name: "ok",
			inputMessage: Message{
				RequestTime: testTime,
				HTTPMethod:  "POST",
				RawRequest:  "raw request text",
			},
			wantedErr: nil,
		},
		{
			name: "wrong time",
			inputMessage: Message{
				RequestTime: testTime.Add(12 * time.Hour),
				HTTPMethod:  "POST",
				RawRequest:  "raw request text",
			},
			wantedErr: errors.New("Message from the future"),
		},
		{
			name: "empty method",
			inputMessage: Message{
				RequestTime: testTime,
				HTTPMethod:  "",
				RawRequest:  "raw request text",
			},
			wantedErr: errors.New("Empty HTTP method field"),
		},
		{
			name: "empty raw request",
			inputMessage: Message{
				RequestTime: testTime,
				HTTPMethod:  "POST",
				RawRequest:  "",
			},
			wantedErr: errors.New("Empty raw request field"),
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			actualErr := validateMessage(tc.inputMessage)

			assert.Equal(t, tc.wantedErr, actualErr)
		})
	}
}
