package logwrap

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestSource(t *testing.T) {
	t.Run("log sends a message to the implementation with a source populated on the message", func(t *testing.T) {
		mockImpl := MockImpl{}
		mockImpl.On("Impl", mock.Anything, mock.Anything).Once()

		expectedSource := "this-component"

		logger := New(mockImpl.Impl)
		logger.Log(context.Background(), "anything", Source(expectedSource))

		assert.True(t, mockImpl.AssertExpectations(t))

		capturedMessage := mockImpl.Calls[0].Arguments.Get(1).(Message)

		assert.Equal(t, expectedSource, capturedMessage.Source)
	})
}
