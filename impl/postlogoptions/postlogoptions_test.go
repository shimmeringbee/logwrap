package postlogoptions

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

func TestPostLogOptions(t *testing.T) {
	t.Run("post log options runs options against the log message provided and forwards onto the implementation", func(t *testing.T) {
		mockImplOne := MockImpl{}
		mockImplOne.On("Impl", mock.Anything, mock.Anything).Once()

		expectedKey := "key"
		expectedValue := "value"

		postlog := PostLogOptions(mockImplOne.Impl, logwrap.Datum(expectedKey, expectedValue))

		expectedMessage := logwrap.Message{Message: "message", Data: map[string]interface{}{}}
		postlog(context.Background(), expectedMessage)

		capturedMessage := mockImplOne.Calls[0].Arguments.Get(1).(logwrap.Message)

		assert.Equal(t, expectedMessage.Message, capturedMessage.Message)
		assert.Equal(t, expectedValue, capturedMessage.Data[expectedKey])
	})
}
