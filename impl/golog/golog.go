package golog

import (
	"context"
	"encoding/json"
	"github.com/shimmeringbee/logwrap"
	"log"
)

// Wrap implements a log/Logger wrapper, allowing the output from logwrap to be sent to the standard Go log package.
//
// log/Logger does not support any mechanism for overriding the timestamp, as such the message timestamp will be
// ignored.
func Wrap(logger *log.Logger) logwrap.Impl {
	return func(ctx context.Context, message logwrap.Message) {
		fieldData, err := json.Marshal(message.Fields)
		if err != nil {
			fieldData = []byte("{}")
		}

		var logIt func(format string, v ...interface{})

		switch message.Level {
		case logwrap.Panic:
			logIt = logger.Panicf
		case logwrap.Fatal:
			logIt = logger.Fatalf
		default:
			logIt = logger.Printf
		}

		logIt("[%s] \"%s\" %s", message.Level.String(), message.Message, fieldData)
	}
}
