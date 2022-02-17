package logwrap

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestTrail(t *testing.T) {
	t.Run("the first trail populates the trail data field with the name provided", func(t *testing.T) {
		mockImpl := MockImpl{}
		mockImpl.On("Impl", mock.Anything, mock.Anything).Once()

		expectedValue := "root"

		logger := New(mockImpl.Impl)
		logger.Log(context.Background(), "anything", Trail("root"))

		assert.True(t, mockImpl.AssertExpectations(t))

		capturedMessage := mockImpl.Calls[0].Arguments.Get(1).(Message)

		assert.Equal(t, expectedValue, capturedMessage.Data[trailField])
	})

	t.Run("multiple trails result in a period delimited field", func(t *testing.T) {
		mockImpl := MockImpl{}
		mockImpl.On("Impl", mock.Anything, mock.Anything).Once()

		expectedValue := "root.sub"

		logger := New(mockImpl.Impl)
		logger.Log(context.Background(), "anything", Trail("root"), Trail("sub"))

		assert.True(t, mockImpl.AssertExpectations(t))

		capturedMessage := mockImpl.Calls[0].Arguments.Get(1).(Message)

		assert.Equal(t, expectedValue, capturedMessage.Data[trailField])
	})

}
