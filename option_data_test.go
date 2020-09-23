package logwrap

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestField(t *testing.T) {
	t.Run("the datum option inserts one key value", func(t *testing.T) {
		mockImpl := MockImpl{}
		mockImpl.On("Impl", mock.Anything, mock.Anything).Once()

		expectedKey := "key"
		expectedValue := "value"

		logger := New(mockImpl.Impl)
		logger.Log(context.Background(), "anything", Datum(expectedKey, expectedValue))

		assert.True(t, mockImpl.AssertExpectations(t))

		capturedMessage := mockImpl.Calls[0].Arguments.Get(1).(Message)

		assert.Equal(t, expectedValue, capturedMessage.Data[expectedKey])
	})
}

func TestFields(t *testing.T) {
	t.Run("the data option inserts the map of key/values", func(t *testing.T) {
		mockImpl := MockImpl{}
		mockImpl.On("Impl", mock.Anything, mock.Anything).Once()

		expectedKeyOne := "key1"
		expectedValueOne := "value1"
		expectedKeyTwo := "key2"
		expectedValueTwo := "value2"

		logger := New(mockImpl.Impl)
		logger.Log(context.Background(), "anything", Data(List{expectedKeyOne: expectedValueOne, expectedKeyTwo: expectedValueTwo}))

		assert.True(t, mockImpl.AssertExpectations(t))

		capturedMessage := mockImpl.Calls[0].Arguments.Get(1).(Message)

		assert.Equal(t, expectedValueOne, capturedMessage.Data[expectedKeyOne])
		assert.Equal(t, expectedValueTwo, capturedMessage.Data[expectedKeyTwo])
	})
}
