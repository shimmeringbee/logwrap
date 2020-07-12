package tee

import (
	"context"
	"github.com/shimmeringbee/logwrap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockImpl struct {
	mock.Mock
}

func (l *MockImpl) Impl(ctx context.Context, msg logwrap.Message) {
	l.Called(ctx, msg)
}

func TestTee(t *testing.T) {
	t.Run("tee distributes the log message to each implementation passed", func(t *testing.T) {
		mockImplOne := MockImpl{}
		mockImplOne.On("Impl", mock.Anything, mock.Anything).Once()

		mockImplTwo := MockImpl{}
		mockImplTwo.On("Impl", mock.Anything, mock.Anything).Once()

		tee := Tee(mockImplOne.Impl, mockImplTwo.Impl)

		expectedMessage := logwrap.Message{Message: "message"}
		tee(context.Background(), expectedMessage)

		assert.Equal(t, expectedMessage, mockImplOne.Calls[0].Arguments.Get(1).(logwrap.Message))
		assert.Equal(t, expectedMessage, mockImplTwo.Calls[0].Arguments.Get(1).(logwrap.Message))
	})
}
