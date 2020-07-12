package logwrap

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestLevel(t *testing.T) {
	t.Run("log sends a message to the implementation with an overridden log level", func(t *testing.T) {
		mockImpl := MockImpl{}
		mockImpl.On("Impl", mock.Anything, mock.Anything).Once()

		expectedLevel := Fatal

		logger := New(mockImpl.Impl)
		logger.Log(context.Background(), "anything", Level(expectedLevel))

		assert.True(t, mockImpl.AssertExpectations(t))

		capturedMessage := mockImpl.Calls[0].Arguments.Get(1).(Message)

		assert.Equal(t, expectedLevel, capturedMessage.Level)
	})
}
