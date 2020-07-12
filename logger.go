package logwrap

import (
	"context"
	"sync/atomic"
	"time"
)

// Impl is the interface for the actual implementation of logging to implement.
//
// The default recommendation is that implementations should not block during execution. The implementation should make
// use of go concurrency techniques to remove blocking code from calling functions go routine.
//
// Should an implementation block by design (such as assured delivery of logs), this should be made explicitly clear in
// any documentation.
//
// Implementations should obey the semantics of Panic and Fatal levels, panic()ing and os.Exit(-1) respectively after
// the log has been made.
type Impl func(context.Context, Message)

// LogLevel is the log level type.
type LogLevel uint

// String provides a text description of the level.
func (l LogLevel) String() string {
	switch l {
	case Panic:
		return "PANIC"
	case Fatal:
		return "FATAL"
	case Error:
		return "ERROR"
	case Warn:
		return "WARN"
	case Info:
		return "INFO"
	case Debug:
		return "DEBUG"
	case Trace:
		return "TRACE"
	default:
		return "UNKNOWN"
	}
}

// Possible log levels that messages can be made against.
const (
	// Panic level, the error encountered immediately panics the application.
	Panic LogLevel = iota
	// Fatal level, a severe enough issue has occurred the the application can no longer continue.
	Fatal
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
	impl      Impl
	sequence  *uint64
	unique    uint64
	segmentID *uint64
	options   []Option
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
	var initialSegmentID uint64
	return Logger{
		impl:      i,
		sequence:  &initialSequence,
		unique:    atomic.AddUint64(loggerSequence, 1),
		segmentID: &initialSegmentID,
		options:   []Option{},
	}
}

// AddOptionsToLogger adds default options to the logger which do not vary by implementation, and are applied first
// before any context or log specific messages.
func (l *Logger) AddOptionsToLogger(options ...Option) {
	l.options = append(l.options, options...)
}
