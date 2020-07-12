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

		switch message.Level {
		case logwrap.Panic:
			logger.Panicf("[%s] \"%s\" %s", message.Level.String(), message.Message, fieldData)
		case logwrap.Fatal:
			logger.Fatalf("[%s] \"%s\" %s", message.Level.String(), message.Message, fieldData)
		default:
			logger.Printf("[%s] \"%s\" %s", message.Level.String(), message.Message, fieldData)
		}
	}
}
