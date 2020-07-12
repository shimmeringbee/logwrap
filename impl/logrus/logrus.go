package logrus

import (
	"context"
	"github.com/shimmeringbee/logwrap"
	realLogrus "github.com/sirupsen/logrus"
)

func Wrap(dest *realLogrus.Logger) logwrap.Impl {
	return func(ctx context.Context, message logwrap.Message) {
		dest.WithFields(message.Fields).WithTime(message.Timestamp).Log(mapLogLevels(message.Level), message.Message)
	}
}

func mapLogLevels(level logwrap.LogLevel) realLogrus.Level {
	switch level {
	case logwrap.Fatal:
		return realLogrus.FatalLevel
	case logwrap.Error:
		return realLogrus.ErrorLevel
	case logwrap.Warn:
		return realLogrus.WarnLevel
	case logwrap.Info:
		return realLogrus.InfoLevel
	case logwrap.Debug:
		return realLogrus.DebugLevel
	case logwrap.Trace:
		return realLogrus.TraceLevel
	default:
		return realLogrus.InfoLevel
	}
}
