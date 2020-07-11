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
