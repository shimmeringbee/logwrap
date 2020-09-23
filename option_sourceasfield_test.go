package logwrap

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestSourceAsField(t *testing.T) {
	t.Run("log sends a message to the implementation with source copied to field", func(t *testing.T) {
		mockImpl := MockImpl{}
		mockImpl.On("Impl", mock.Anything, mock.Anything).Once()

		expectedSource := "this-source"

		logger := New(mockImpl.Impl)
		logger.Log(context.Background(), "anything", Source(expectedSource), SourceAsField)

		assert.True(t, mockImpl.AssertExpectations(t))

		capturedMessage := mockImpl.Calls[0].Arguments.Get(1).(Message)

		assert.Equal(t, expectedSource, capturedMessage.Data["source"])
	})
}
