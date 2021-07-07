package logger

import (
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	once   sync.Once
	logger *logrus.Logger
)

func Get(level logrus.Level) *logrus.Logger {
	once.Do(func() {
		log := logrus.New()
		log.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: time.RFC822,
		})

		log.SetLevel(level)
		logger = log
	})

	return logger
}
