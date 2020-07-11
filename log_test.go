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

func TestLogger_AddOptionsToContext(t *testing.T) {
	t.Run("adds an option to the context which is used in following logs", func(t *testing.T) {
		mockImpl := MockImpl{}
		mockImpl.On("Impl", mock.Anything, mock.Anything).Once()

		expectedMessage := "message"
		expectedLevel := Fatal

		logger := New(mockImpl.Impl)

		ctx := logger.AddOptionsToContext(context.Background(), Level(Fatal))
		logger.Log(ctx, expectedMessage)

		assert.True(t, mockImpl.AssertExpectations(t))

		capturedMessage := mockImpl.Calls[0].Arguments.Get(1).(Message)
		assert.Equal(t, expectedLevel, capturedMessage.Level)
		assert.Equal(t, expectedMessage, capturedMessage.Message)
	})

	t.Run("multiple descendant contexts correctly build upon options", func(t *testing.T) {
		mockImpl := MockImpl{}
		mockImpl.On("Impl", mock.Anything, mock.Anything).Once()

		expectedMessage := "message"
		expectedLevel := Fatal

		expectedKey := "key"
		expectedValue := "value"

		logger := New(mockImpl.Impl)

		ctx := logger.AddOptionsToContext(context.Background(), Level(Fatal))
		cctx := logger.AddOptionsToContext(ctx, Field(expectedKey, expectedValue))
		logger.Log(cctx, expectedMessage)

		assert.True(t, mockImpl.AssertExpectations(t))

		capturedMessage := mockImpl.Calls[0].Arguments.Get(1).(Message)
		assert.Equal(t, expectedLevel, capturedMessage.Level)
		assert.Equal(t, expectedMessage, capturedMessage.Message)
		assert.Equal(t, expectedValue, capturedMessage.Fields[expectedKey])
	})

	t.Run("two separate loggers do not conflict within the same context", func(t *testing.T) {
		mockImpl := MockImpl{}
		mockImpl.On("Impl", mock.Anything, mock.Anything).Twice()

		expectedMessage := "message"
		expectedLevelOne := Fatal
		expectedLevelTwo := Fatal

		logger := New(mockImpl.Impl)
		loggerTwo := New(mockImpl.Impl)

		ctx := loggerTwo.AddOptionsToContext(context.Background(), Level(expectedLevelTwo))
		cctx := logger.AddOptionsToContext(ctx, Level(expectedLevelOne))

		logger.Log(cctx, expectedMessage)
		loggerTwo.Log(cctx, expectedMessage)

		assert.True(t, mockImpl.AssertExpectations(t))

		capturedMessage := mockImpl.Calls[0].Arguments.Get(1).(Message)
		assert.Equal(t, expectedLevelOne, capturedMessage.Level)

		capturedMessageTwo := mockImpl.Calls[1].Arguments.Get(1).(Message)
		assert.Equal(t, expectedLevelTwo, capturedMessageTwo.Level)
	})
}

type MockImpl struct {
	mock.Mock
}

func (l *MockImpl) Impl(ctx context.Context, msg Message) {
	l.Called(ctx, msg)
}
