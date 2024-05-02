package core

import (
	"homework/internal/pkg/kafkalogger"
)

type loggerOps interface {
	LogMessage(message kafkalogger.Message) error
	LogMessages(message []kafkalogger.Message) error
}

func (s *Service) LogMessage(message kafkalogger.Message) error {
	return s.logger.LogMessage(message)
}

func (s *Service) LogMessages(messages []kafkalogger.Message) error {
	return s.logger.LogMessages(messages)
}
