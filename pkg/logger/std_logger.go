package logger

import (
	"log"
)

// StdLogger адаптирует наш логгер к стандартному интерфейсу log.Logger
type StdLogger struct {
	logger *Logger
}

// NewStdLogger создает новый адаптер для стандартного логгера
func NewStdLogger(logger *Logger) *log.Logger {
	return log.New(&StdLogger{logger: logger}, "", 0)
}

// Write реализует интерфейс io.Writer для log.Logger
func (l *StdLogger) Write(p []byte) (n int, err error) {
	l.logger.Info(nil, string(p), nil)
	return len(p), nil
}
