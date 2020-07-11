package logwrap

import (
	"context"
	"sync/atomic"
	"time"
)

// Impl is the interface for the actual implementation of logging to implement.
type Impl func(context.Context, Message)

// LogLevel is the log level type.
type LogLevel uint

// Possible log levels that messages can be made against.
const (
	// Fatal level, a severe enough issue has occurred the the application can no longer continue.
	Fatal LogLevel = iota
	// Error level, a severe issue has been encountered, but the application has recovered.
	Error
	// Warn level, an issue has occurred which has not caused a operational issue, but should not have happened.
	Warn
	// Info level, general information about the applications progress, decisions or checkpoints reached.
	Info
	// Debug level, verbose logging usually only needed by a operator when fault finding.
	Debug
	// Trace level, extreme diagnostics reporting very fine details, usually only needed by developers.
	Trace
)

const contextKeyOptions = "_ShimmeringBeeLogOptions"
const defaultLevel = Info

var loggerSequence *uint64

func init() {
	var initialSequence uint64
	loggerSequence = &initialSequence
}

// Logger is the representation of a stream of logs, it should always be instantiated with `New`.
type Logger struct {
	impl     Impl
	sequence *uint64
	unique   uint64
}

// Option is an interface for a option a Log call can take, adding or modifying data on a Message.
type Option func(*Message)

// Message structure is the struct sent to a logging implementation, it includes all fields.
type Message struct {
	// Level of log message.
	Level LogLevel
	// Message is the human readable version of the message.
	Message string
	// Fields are a free form map of data to log, usually for structured logging.
	Fields map[string]interface{}
	// Timestamp at which the log was made.
	Timestamp time.Time
	// Sequence is a monotonic sequence number, used to determine log order with high frequency/low interval logs.
	Sequence uint64
}

// New constructs a new logger, taking the backend implement which will actually log.
func New(i Impl) Logger {
	var initialSequence uint64
	return Logger{impl: i, sequence: &initialSequence, unique: atomic.AddUint64(loggerSequence, 1)}
}

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

type contextKey struct {
	base   string
	unique uint64
}

func (l Logger) contextKey(base string) contextKey {
	return contextKey{
		base:   base,
		unique: l.unique,
	}
}

// AddOptionsToContext add default Option's to a context for this specific logger (i.e. two loggers will have different
// options on the same context). These are always processed first, before any Option's provided during Log.
func (l Logger) AddOptionsToContext(ctx context.Context, options ...Option) context.Context {
	optionsToAdd := append(l.getOptionsFromContext(ctx), options...)
	return context.WithValue(ctx, l.contextKey(contextKeyOptions), optionsToAdd)
}

func (l Logger) getOptionsFromContext(ctx context.Context) []Option {
	uncast := ctx.Value(l.contextKey(contextKeyOptions))
	if uncast == nil {
		return []Option{}
	}

	if options, ok := uncast.([]Option); ok {
		return options
	}

	return []Option{}
}
