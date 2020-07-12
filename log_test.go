package logwrap

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestLogger_Log(t *testing.T) {
	t.Run("log sends a message to the implementation with a message and default level of info", func(t *testing.T) {
		mockImpl := MockImpl{}
		mockImpl.On("Impl", mock.Anything, mock.Anything).Once()

		expectedMessage := "message"
		expectedLevel := Info

		logger := New(mockImpl.Impl)
		logger.Log(context.Background(), expectedMessage)

		assert.True(t, mockImpl.AssertExpectations(t))

		capturedMessage := mockImpl.Calls[0].Arguments.Get(1).(Message)
		assert.Equal(t, expectedLevel, capturedMessage.Level)
		assert.Equal(t, expectedMessage, capturedMessage.Message)
	})

	t.Run("log sends a message to the implementation with an increasing sequence number", func(t *testing.T) {
		mockImpl := MockImpl{}
		mockImpl.On("Impl", mock.Anything, mock.Anything).Twice()

		logger := New(mockImpl.Impl)
		logger.Log(context.Background(), "anything")
		logger.Log(context.Background(), "anything")

		assert.True(t, mockImpl.AssertExpectations(t))

		var capturedMessage [2]Message
		capturedMessage[0] = mockImpl.Calls[0].Arguments.Get(1).(Message)
		capturedMessage[1] = mockImpl.Calls[1].Arguments.Get(1).(Message)

		assert.Equal(t, uint64(1), capturedMessage[0].Sequence)
		assert.Equal(t, uint64(2), capturedMessage[1].Sequence)
	})

	t.Run("log sends a message to the implementation with an inserted timestamp", func(t *testing.T) {
		mockImpl := MockImpl{}
		mockImpl.On("Impl", mock.Anything, mock.Anything).Once()

		beforeTime := time.Now()

		logger := New(mockImpl.Impl)
		logger.Log(context.Background(), "anything")

		afterTime := time.Now()

		assert.True(t, mockImpl.AssertExpectations(t))

		capturedMessage := mockImpl.Calls[0].Arguments.Get(1).(Message)

		assert.True(t, beforeTime.Before(capturedMessage.Timestamp) || beforeTime.Equal(capturedMessage.Timestamp))
		assert.True(t, afterTime.After(capturedMessage.Timestamp) || beforeTime.Equal(capturedMessage.Timestamp))
	})
}

func TestLogger_LogFatal(t *testing.T) {
	t.Run("log sends a message to the implementation with level of fatal", func(t *testing.T) {
		mockImpl := MockImpl{}
		mockImpl.On("Impl", mock.Anything, mock.Anything).Once()

		expectedLevel := Fatal

		logger := New(mockImpl.Impl)
		logger.LogFatal(context.Background(), "message")

		assert.True(t, mockImpl.AssertExpectations(t))

		capturedMessage := mockImpl.Calls[0].Arguments.Get(1).(Message)
		assert.Equal(t, expectedLevel, capturedMessage.Level)
	})
}

func TestLogger_LogError(t *testing.T) {
	t.Run("log sends a message to the implementation with level of error", func(t *testing.T) {
		mockImpl := MockImpl{}
		mockImpl.On("Impl", mock.Anything, mock.Anything).Once()

		expectedLevel := Error

		logger := New(mockImpl.Impl)
		logger.LogError(context.Background(), "message")

		assert.True(t, mockImpl.AssertExpectations(t))

		capturedMessage := mockImpl.Calls[0].Arguments.Get(1).(Message)
		assert.Equal(t, expectedLevel, capturedMessage.Level)
	})
}

func TestLogger_LogWarn(t *testing.T) {
	t.Run("log sends a message to the implementation with level of warn", func(t *testing.T) {
		mockImpl := MockImpl{}
		mockImpl.On("Impl", mock.Anything, mock.Anything).Once()

		expectedLevel := Warn

		logger := New(mockImpl.Impl)
		logger.LogWarn(context.Background(), "message")

		assert.True(t, mockImpl.AssertExpectations(t))

		capturedMessage := mockImpl.Calls[0].Arguments.Get(1).(Message)
		assert.Equal(t, expectedLevel, capturedMessage.Level)
	})
}

func TestLogger_LogInfo(t *testing.T) {
	t.Run("log sends a message to the implementation with level of info", func(t *testing.T) {
		mockImpl := MockImpl{}
		mockImpl.On("Impl", mock.Anything, mock.Anything).Once()

		expectedLevel := Info

		logger := New(mockImpl.Impl)
		ctx := logger.AddOptionsToContext(context.Background(), Level(Fatal))
		logger.LogInfo(ctx, "message")

		assert.True(t, mockImpl.AssertExpectations(t))

		capturedMessage := mockImpl.Calls[0].Arguments.Get(1).(Message)
		assert.Equal(t, expectedLevel, capturedMessage.Level)
	})
}

func TestLogger_LogDebug(t *testing.T) {
	t.Run("log sends a message to the implementation with level of debug", func(t *testing.T) {
		mockImpl := MockImpl{}
		mockImpl.On("Impl", mock.Anything, mock.Anything).Once()

		expectedLevel := Debug

		logger := New(mockImpl.Impl)
		logger.LogDebug(context.Background(), "message")

		assert.True(t, mockImpl.AssertExpectations(t))

		capturedMessage := mockImpl.Calls[0].Arguments.Get(1).(Message)
		assert.Equal(t, expectedLevel, capturedMessage.Level)
	})
}

func TestLogger_LogTrace(t *testing.T) {
	t.Run("log sends a message to the implementation with level of trace", func(t *testing.T) {
		mockImpl := MockImpl{}
		mockImpl.On("Impl", mock.Anything, mock.Anything).Once()

		expectedLevel := Trace

		logger := New(mockImpl.Impl)
		logger.LogTrace(context.Background(), "message")

		assert.True(t, mockImpl.AssertExpectations(t))

		capturedMessage := mockImpl.Calls[0].Arguments.Get(1).(Message)
		assert.Equal(t, expectedLevel, capturedMessage.Level)
	})
}
