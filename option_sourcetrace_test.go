package logwrap

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"strings"
	"testing"
)

// Used for testing to keep log line in a static place, any change above the log line will require expectations changing.
func callLog(ctx context.Context, logger Logger, message string) {
	logger.Log(ctx, message, SourceTrace)
}

func TestSourceTrace(t *testing.T) {
	expectedFileSuffix := "option_sourcetrace_test.go"
	expectedLine := 13
	expectedFunction := "github.com/shimmeringbee/logwrap.callLog"

	t.Run("logs the first before anything in logwrap is called", func(t *testing.T) {
		mockImpl := MockImpl{}
		mockImpl.On("Impl", mock.Anything, mock.Anything).Once()

		logger := New(mockImpl.Impl)
		callLog(context.Background(), logger, "anything")

		assert.True(t, mockImpl.AssertExpectations(t))
		capturedMessage := mockImpl.Calls[0].Arguments.Get(1).(Message)

		actualSourceLocation := capturedMessage.Fields[SourceTraceField].(SourceLocation)
		assert.Equal(t, expectedLine, actualSourceLocation.Line)
		assert.Equal(t, expectedFunction, actualSourceLocation.Function)
		assert.True(t, strings.HasSuffix(actualSourceLocation.File, expectedFileSuffix))
	})
}
