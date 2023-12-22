package logger

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	logger         *logrus.Logger
	terminalLogger *logrus.Logger
}

func NewLogger(output io.Writer, logLevel logrus.Level) *Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(output)
	logger.SetLevel(logLevel)

	// separate logger to log debug level notifications to stdout
	terminalLogger := logrus.New()
	terminalLogger.SetOutput(os.Stdout)
	terminalLogger.SetFormatter(&logrus.TextFormatter{})
	terminalLogger.SetLevel(logrus.DebugLevel)
	return &Logger{
		logger:         logger,
		terminalLogger: terminalLogger,
	}
}

func (l *Logger) Error(msg string, fields map[string]interface{}) {
	l.terminalLogger.WithFields(fields).Error(msg)
	l.logger.WithFields(fields).Error(msg)
}

func (l *Logger) Warning(msg string, fields map[string]interface{}) {
	l.terminalLogger.WithFields(fields).Warning(msg)
	l.logger.WithFields(fields).Warning(msg)
}

func (l *Logger) Info(msg string, fields map[string]interface{}) {
	l.terminalLogger.WithFields(fields).Info(msg)
	l.logger.WithFields(fields).Info(msg)
}
