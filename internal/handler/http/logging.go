package http

import "github.com/sirupsen/logrus"

func logField(handler string) logrus.Fields {
	return logrus.Fields{
		"handler": handler,
	}
}

func logError(handler string, err error) {
	logrus.WithFields(logField(handler)).Error(err)
}
