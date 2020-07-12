package logwrap

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestSequenceAsField(t *testing.T) {
	t.Run("log sends a message to the implementation with sequence copied to field", func(t *testing.T) {
		mockImpl := MockImpl{}
		mockImpl.On("Impl", mock.Anything, mock.Anything).Once()

		logger := New(mockImpl.Impl)
		logger.Log(context.Background(), "anything", SequenceAsField)

		assert.True(t, mockImpl.AssertExpectations(t))

		capturedMessage := mockImpl.Calls[0].Arguments.Get(1).(Message)

		assert.Equal(t, uint64(1), capturedMessage.Fields["sequence"])
	})
}
