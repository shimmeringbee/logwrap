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
		Data:      map[string]interface{}{},
		Timestamp: time.Now(),
		Sequence:  atomic.AddUint64(l.sequence, 1),
	}

	for _, option := range l.options {
		option(&outgoingMessage)
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

// Panic is an alias for LogPanic.
func (l Logger) Panic(ctx context.Context, message string, options ...Option) {
	l.LogPanic(ctx, message, options...)
}

// LogFatal calls Log while appending Level(Fatal) as an option.
func (l Logger) LogFatal(ctx context.Context, message string, options ...Option) {
	options = append(options, Level(Fatal))
	l.Log(ctx, message, options...)
}

// Fatal is an alias for LogFatal.
func (l Logger) Fatal(ctx context.Context, message string, options ...Option) {
	l.LogFatal(ctx, message, options...)
}

// LogError calls Log while appending Level(Error) as an option.
func (l Logger) LogError(ctx context.Context, message string, options ...Option) {
	options = append(options, Level(Error))
	l.Log(ctx, message, options...)
}

// Error is an alias for LogError.
func (l Logger) Error(ctx context.Context, message string, options ...Option) {
	l.LogError(ctx, message, options...)
}

// LogWarn calls Log while appending Level(Warn) as an option.
func (l Logger) LogWarn(ctx context.Context, message string, options ...Option) {
	options = append(options, Level(Warn))
	l.Log(ctx, message, options...)
}

// Warn is an alias for LogWarn.
func (l Logger) Warn(ctx context.Context, message string, options ...Option) {
	l.LogWarn(ctx, message, options...)
}

// LogInfo calls Log while appending Level(Info) as an option.
func (l Logger) LogInfo(ctx context.Context, message string, options ...Option) {
	options = append(options, Level(Info))
	l.Log(ctx, message, options...)
}

// Info is an alias for LogInfo.
func (l Logger) Info(ctx context.Context, message string, options ...Option) {
	l.LogInfo(ctx, message, options...)
}

// LogDebug calls Log while appending Level(Debug) as an option.
func (l Logger) LogDebug(ctx context.Context, message string, options ...Option) {
	options = append(options, Level(Debug))
	l.Log(ctx, message, options...)
}

// Debug is an alias for LogDebug.
func (l Logger) Debug(ctx context.Context, message string, options ...Option) {
	l.LogDebug(ctx, message, options...)
}

// LogTrace calls Log while appending Level(Trace) as an option.
func (l Logger) LogTrace(ctx context.Context, message string, options ...Option) {
	options = append(options, Level(Trace))
	l.Log(ctx, message, options...)
}

// Trace is an alias for LogTrace.
func (l Logger) Trace(ctx context.Context, message string, options ...Option) {
	l.LogTrace(ctx, message, options...)
}
