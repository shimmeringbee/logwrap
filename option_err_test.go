package logwrap

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestErr(t *testing.T) {
	t.Run("log sends a message to the implementation with an error in fields as err", func(t *testing.T) {
		mockImpl := MockImpl{}
		mockImpl.On("Impl", mock.Anything, mock.Anything).Once()

		expectedErr := errors.New("expected error")

		logger := New(mockImpl.Impl)
		logger.Log(context.Background(), "anything", Err(expectedErr))

		assert.True(t, mockImpl.AssertExpectations(t))

		capturedMessage := mockImpl.Calls[0].Arguments.Get(1).(Message)

		assert.Equal(t, expectedErr.Error(), capturedMessage.Data["err"])
	})
}
