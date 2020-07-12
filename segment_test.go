package logwrap

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestLogger_Segment(t *testing.T) {
	t.Run("starting a segment outputs a message, and closing a segment also outputs a message with fields indicating begin/end", func(t *testing.T) {
		mockImpl := MockImpl{}
		mockImpl.On("Impl", mock.Anything, mock.Anything).Times(3)

		expectedMessage := "message"
		expectedInnerMessage := "inner message"
		expectedKey := "key"
		expectedValue := "value"

		logger := New(mockImpl.Impl)
		ctx, end := logger.Segment(context.Background(), expectedMessage, Field(expectedKey, expectedValue))
		logger.Log(ctx, expectedInnerMessage)
		end()

		assert.True(t, mockImpl.AssertExpectations(t))

		var capturedMessage [3]Message
		capturedMessage[0] = mockImpl.Calls[0].Arguments.Get(1).(Message)
		capturedMessage[1] = mockImpl.Calls[1].Arguments.Get(1).(Message)
		capturedMessage[2] = mockImpl.Calls[2].Arguments.Get(1).(Message)

		assert.Equal(t, expectedMessage, capturedMessage[0].Message)
		assert.Equal(t, SegmentStartValue, capturedMessage[0].Fields[SegmentField])
		assert.Equal(t, expectedValue, capturedMessage[0].Fields[expectedKey])

		assert.Equal(t, expectedInnerMessage, capturedMessage[1].Message)
		assert.Equal(t, expectedValue, capturedMessage[1].Fields[expectedKey])

		assert.Equal(t, expectedMessage, capturedMessage[2].Message)
		assert.Equal(t, SegmentEndValue, capturedMessage[2].Fields[SegmentField])
		assert.Equal(t, expectedValue, capturedMessage[2].Fields[expectedKey])
	})

	t.Run("segment created has field with unique segment id", func(t *testing.T) {
		mockImpl := MockImpl{}
		mockImpl.On("Impl", mock.Anything, mock.Anything).Times(2)

		logger := New(mockImpl.Impl)
		logger.Segment(context.Background(), "")
		logger.Segment(context.Background(), "")

		var capturedMessage [2]Message
		capturedMessage[0] = mockImpl.Calls[0].Arguments.Get(1).(Message)
		capturedMessage[1] = mockImpl.Calls[1].Arguments.Get(1).(Message)

		assert.Equal(t, uint64(1), capturedMessage[0].Fields[SegmentIDField])
		assert.Equal(t, uint64(2), capturedMessage[1].Fields[SegmentIDField])
	})

	t.Run("segment created as child of another segment has the parents segment id", func(t *testing.T) {
		mockImpl := MockImpl{}
		mockImpl.On("Impl", mock.Anything, mock.Anything).Times(2)

		logger := New(mockImpl.Impl)
		ctx, _ := logger.Segment(context.Background(), "")
		logger.Segment(ctx, "")

		var capturedMessage [2]Message
		capturedMessage[0] = mockImpl.Calls[0].Arguments.Get(1).(Message)
		capturedMessage[1] = mockImpl.Calls[1].Arguments.Get(1).(Message)

		assert.Equal(t, uint64(1), capturedMessage[0].Fields[SegmentIDField])
		assert.Nil(t, capturedMessage[0].Fields[ParentSegmentIDField])

		assert.Equal(t, uint64(2), capturedMessage[1].Fields[SegmentIDField])
		assert.Equal(t, uint64(1), capturedMessage[1].Fields[ParentSegmentIDField])
	})
}
