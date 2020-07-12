package logwrap

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestField(t *testing.T) {
	t.Run("log sends a message to the implementation with an overridden log level", func(t *testing.T) {
		mockImpl := MockImpl{}
		mockImpl.On("Impl", mock.Anything, mock.Anything).Once()

		expectedKey := "key"
		expectedValue := "value"

		logger := New(mockImpl.Impl)
		logger.Log(context.Background(), "anything", Field(expectedKey, expectedValue))

		assert.True(t, mockImpl.AssertExpectations(t))

		capturedMessage := mockImpl.Calls[0].Arguments.Get(1).(Message)

		assert.Equal(t, expectedValue, capturedMessage.Fields[expectedKey])
	})
}

func TestFields(t *testing.T) {
	t.Run("log sends a message to the implementation with an overridden log level", func(t *testing.T) {
		mockImpl := MockImpl{}
		mockImpl.On("Impl", mock.Anything, mock.Anything).Once()

		expectedKeyOne := "key1"
		expectedValueOne := "value1"
		expectedKeyTwo := "key2"
		expectedValueTwo := "value2"

		logger := New(mockImpl.Impl)
		logger.Log(context.Background(), "anything", Fields(List{expectedKeyOne: expectedValueOne, expectedKeyTwo: expectedValueTwo}))

		assert.True(t, mockImpl.AssertExpectations(t))

		capturedMessage := mockImpl.Calls[0].Arguments.Get(1).(Message)

		assert.Equal(t, expectedValueOne, capturedMessage.Fields[expectedKeyOne])
		assert.Equal(t, expectedValueTwo, capturedMessage.Fields[expectedKeyTwo])
	})
}
