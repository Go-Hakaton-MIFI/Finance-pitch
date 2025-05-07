package logger

import (
	"context"
	"finance-backend/pkg/utils"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	instance *logrus.Logger
}

func NewLogger() *Logger {
	log := logrus.New()

	log.SetOutput(os.Stdout)

	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})

	log.SetLevel(logrus.InfoLevel)

	return &Logger{instance: log}
}

func (l *Logger) Info(ctx context.Context, message string, fields map[string]interface{}) {
	if fields == nil {
		fields = make(map[string]interface{})
	}

	if requestID, ok := utils.GetRequestIDFromContext(ctx); ok {
		fields["request_id"] = requestID
	}

	l.instance.WithFields(fields).Info(message)
}

func (l *Logger) Warn(ctx context.Context, message string, fields map[string]interface{}) {
	if fields == nil {
		fields = make(map[string]interface{})
	}

	if requestID, ok := utils.GetRequestIDFromContext(ctx); ok {
		fields["request_id"] = requestID
	}

	l.instance.WithFields(fields).Warn(message)
}

func (l *Logger) Error(ctx context.Context, message string, fields map[string]interface{}) {
	if fields == nil {
		fields = make(map[string]interface{})
	}

	if requestID, ok := utils.GetRequestIDFromContext(ctx); ok {
		fields["request_id"] = requestID
	}

	l.instance.WithFields(fields).Error(message)
}

func (l *Logger) Debug(ctx context.Context, message string, fields map[string]interface{}) {
	if fields == nil {
		fields = make(map[string]interface{})
	}

	if requestID, ok := utils.GetRequestIDFromContext(ctx); ok {
		fields["request_id"] = requestID
	}

	l.instance.WithFields(fields).Debug(message)
}

func (l *Logger) Fatal(ctx context.Context, message string, fields map[string]interface{}) {
	if fields == nil {
		fields = make(map[string]interface{})
	}

	if requestID, ok := utils.GetRequestIDFromContext(ctx); ok {
		fields["request_id"] = requestID
	}

	l.instance.WithFields(fields).Fatal(message)
}

func (l *Logger) GetInstance() *logrus.Logger {
	return l.instance
}
