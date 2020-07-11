package logwrap

import (
	"context"
	"sync/atomic"
	"time"
)

// Log processes and logs the provided message, applying any options which have been stored in the context first and
// then those passed into Log.
func (l Logger) Log(ctx context.Context, message string, options ...Option) {
	outgoingMessage := Message{
		Level:     defaultLevel,
		Message:   message,
		Fields:    map[string]interface{}{},
		Timestamp: time.Now(),
		Sequence:  atomic.AddUint64(l.sequence, 1),
	}

	for _, option := range l.getOptionsFromContext(ctx) {
		option(&outgoingMessage)
	}

	for _, option := range options {
		option(&outgoingMessage)
	}

	l.impl(ctx, outgoingMessage)
}
