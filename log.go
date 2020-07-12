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

// LogPanic calls Log while appending Level(Panic) as an option.
func (l Logger) LogPanic(ctx context.Context, message string, options ...Option) {
	options = append(options, Level(Panic))
	l.Log(ctx, message, options...)
}

// LogFatal calls Log while appending Level(Fatal) as an option.
func (l Logger) LogFatal(ctx context.Context, message string, options ...Option) {
	options = append(options, Level(Fatal))
	l.Log(ctx, message, options...)
}

// LogError calls Log while appending Level(Error) as an option.
func (l Logger) LogError(ctx context.Context, message string, options ...Option) {
	options = append(options, Level(Error))
	l.Log(ctx, message, options...)
}

// LogWarn calls Log while appending Level(Warn) as an option.
func (l Logger) LogWarn(ctx context.Context, message string, options ...Option) {
	options = append(options, Level(Warn))
	l.Log(ctx, message, options...)
}

// LogInfo calls Log while appending Level(Info) as an option.
func (l Logger) LogInfo(ctx context.Context, message string, options ...Option) {
	options = append(options, Level(Info))
	l.Log(ctx, message, options...)
}

// LogDebug calls Log while appending Level(Debug) as an option.
func (l Logger) LogDebug(ctx context.Context, message string, options ...Option) {
	options = append(options, Level(Debug))
	l.Log(ctx, message, options...)
}

// LogTrace calls Log while appending Level(Trace) as an option.
func (l Logger) LogTrace(ctx context.Context, message string, options ...Option) {
	options = append(options, Level(Trace))
	l.Log(ctx, message, options...)
}
