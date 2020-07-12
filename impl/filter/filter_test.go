package filter

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

func TestFilter(t *testing.T) {
	t.Run("filter will only call impl if a message matches", func(t *testing.T) {
		mockImplOne := MockImpl{}
		mockImplOne.On("Impl", mock.Anything, mock.Anything).Once()

		matchingMessage := logwrap.Message{Message: "match"}
		notMatchingMessage := logwrap.Message{Message: "nomatch"}

		filter := Filter(mockImplOne.Impl, func(message logwrap.Message) bool {
			return message.Message == "match"
		})

		filter(context.Background(), matchingMessage)
		filter(context.Background(), notMatchingMessage)

		capturedMessage := mockImplOne.Calls[0].Arguments.Get(1).(logwrap.Message)

		assert.Equal(t, matchingMessage.Message, capturedMessage.Message)
	})
}
