package core

import (
	"homework/internal/pkg/kafkalogger"
)

type loggerOps interface {
	LogMessage(message kafkalogger.Message) error
	LogMessages(message []kafkalogger.Message) error
}

// LogMessage custom log one message
func (s *Service) LogMessage(message kafkalogger.Message) error {
	return s.logger.LogMessage(message)
}

// LogMessages custom log a group of messages
func (s *Service) LogMessages(messages []kafkalogger.Message) error {
	return s.logger.LogMessages(messages)
}
