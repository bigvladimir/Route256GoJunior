package kafkalogger

import (
	"testing"

	"github.com/golang/mock/gomock"

	mock_logger "homework/internal/pkg/kafkalogger/mocks"
)

var testTopic = "logs"

type kafkaLoggerFixtures struct {
	ctrl            *gomock.Controller
	kafkaLogger     KafkaLogger
	mockProducerOps *mock_logger.MockproducerOps
}

func setUpKafkaLogger(t *testing.T) kafkaLoggerFixtures {
	ctrl := gomock.NewController(t)
	mockProducerOps := mock_logger.NewMockproducerOps(ctrl)
	kafkaLogger := KafkaLogger{
		producer: mockProducerOps,
		topic:    testTopic,
	}
	return kafkaLoggerFixtures{
		ctrl:            ctrl,
		kafkaLogger:     kafkaLogger,
		mockProducerOps: mockProducerOps,
	}
}

func (a *kafkaLoggerFixtures) tearDown() {
	a.ctrl.Finish()
}
