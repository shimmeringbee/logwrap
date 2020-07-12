package logwrap

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"math"
	"testing"
)

type MockImpl struct {
	mock.Mock
}

func (l *MockImpl) Impl(ctx context.Context, msg Message) {
	l.Called(ctx, msg)
}

func TestLogLevel_String(t *testing.T) {
	t.Run("log levels output to their correct string representation", func(t *testing.T) {
		assert.Equal(t, "PANIC", Panic.String())
		assert.Equal(t, "FATAL", Fatal.String())
		assert.Equal(t, "ERROR", Error.String())
		assert.Equal(t, "WARN", Warn.String())
		assert.Equal(t, "INFO", Info.String())
		assert.Equal(t, "DEBUG", Debug.String())
		assert.Equal(t, "TRACE", Trace.String())

		// Unknown logwrap level.
		assert.Equal(t, "UNKNOWN", LogLevel(math.MaxUint64).String())
	})
}
