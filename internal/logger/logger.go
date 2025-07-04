package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

type Logger interface {
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	WithField(key string, value interface{}) Logger
	WithFields(fields map[string]interface{}) Logger
}

type JSONLogger struct {
	logger *logrus.Logger
}

func NewLogger() Logger {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
	})
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.InfoLevel)

	return &JSONLogger{logger: log}
}

func (l *JSONLogger) Info(args ...interface{}) {
	l.logger.Info(args...)
}

func (l *JSONLogger) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

func (l *JSONLogger) Error(args ...interface{}) {
	l.logger.Error(args...)
}

func (l *JSONLogger) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}

func (l *JSONLogger) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

func (l *JSONLogger) Warnf(format string, args ...interface{}) {
	l.logger.Warnf(format, args...)
}

func (l *JSONLogger) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}

func (l *JSONLogger) Debugf(format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}

func (l *JSONLogger) Fatalf(format string, args ...interface{}) {
	l.logger.Fatalf(format, args...)
}

func (l *JSONLogger) WithField(key string, value interface{}) Logger {
	return &JSONLogger{logger: l.logger.WithField(key, value).Logger}
}

func (l *JSONLogger) WithFields(fields map[string]interface{}) Logger {
	return &JSONLogger{logger: l.logger.WithFields(logrus.Fields(fields)).Logger}
}
