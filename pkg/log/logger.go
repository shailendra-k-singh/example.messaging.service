package log

import (
	log "github.com/sirupsen/logrus"
)

func NewLogger(level string, format string) *log.Logger {
	logger := log.StandardLogger()

	l, err := log.ParseLevel(level)
	if err != nil {
		logger.Errorf("Error while setting input log level:%s, setting default INFO level", err)
		logger.SetLevel(log.InfoLevel)
		return logger
	}
	logger.SetLevel(l)
	if format == "json" {
		logger.SetFormatter(&log.JSONFormatter{TimestampFormat: "2020/08/23 19:00:00.000"})
	}
	return logger
}
