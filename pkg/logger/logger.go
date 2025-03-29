package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

func Init() {
	// Настраиваем формат вывода логов
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	log.SetOutput(os.Stdout)
	log.SetLevel(logrus.InfoLevel)
}

func Info(msg string, fields logrus.Fields) {
	log.WithFields(fields).Info(msg)
}

func Error(msg string, err error, fields logrus.Fields) {
	log.WithFields(fields).WithError(err).Error(msg)
}

func Debug(msg string, fields logrus.Fields) {
	log.WithFields(fields).Debug(msg)
}
