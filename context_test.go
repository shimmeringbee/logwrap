package logwrap

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestLogger_AddOptionsToContext(t *testing.T) {
	t.Run("adds an option to the context which is used in following logs", func(t *testing.T) {
		mockImpl := MockImpl{}
		mockImpl.On("Impl", mock.Anything, mock.Anything).Once()

		expectedMessage := "message"
		expectedLevel := Fatal

		logger := New(mockImpl.Impl)

		ctx := logger.AddOptionsToContext(context.Background(), Level(Fatal))
		logger.Log(ctx, expectedMessage)

		assert.True(t, mockImpl.AssertExpectations(t))

		capturedMessage := mockImpl.Calls[0].Arguments.Get(1).(Message)
		assert.Equal(t, expectedLevel, capturedMessage.Level)
		assert.Equal(t, expectedMessage, capturedMessage.Message)
	})

	t.Run("multiple descendant contexts correctly build upon options", func(t *testing.T) {
		mockImpl := MockImpl{}
		mockImpl.On("Impl", mock.Anything, mock.Anything).Once()

		expectedMessage := "message"
		expectedLevel := Fatal

		expectedKey := "key"
		expectedValue := "value"

		logger := New(mockImpl.Impl)

		ctx := logger.AddOptionsToContext(context.Background(), Level(Fatal))
		cctx := logger.AddOptionsToContext(ctx, Datum(expectedKey, expectedValue))
		logger.Log(cctx, expectedMessage)

		assert.True(t, mockImpl.AssertExpectations(t))

		capturedMessage := mockImpl.Calls[0].Arguments.Get(1).(Message)
		assert.Equal(t, expectedLevel, capturedMessage.Level)
		assert.Equal(t, expectedMessage, capturedMessage.Message)
		assert.Equal(t, expectedValue, capturedMessage.Data[expectedKey])
	})

	t.Run("two separate loggers do not conflict within the same context", func(t *testing.T) {
		mockImpl := MockImpl{}
		mockImpl.On("Impl", mock.Anything, mock.Anything).Twice()

		expectedMessage := "message"
		expectedLevelOne := Fatal
		expectedLevelTwo := Fatal

		logger := New(mockImpl.Impl)
		loggerTwo := New(mockImpl.Impl)

		ctx := loggerTwo.AddOptionsToContext(context.Background(), Level(expectedLevelTwo))
		cctx := logger.AddOptionsToContext(ctx, Level(expectedLevelOne))

		logger.Log(cctx, expectedMessage)
		loggerTwo.Log(cctx, expectedMessage)

		assert.True(t, mockImpl.AssertExpectations(t))

		capturedMessage := mockImpl.Calls[0].Arguments.Get(1).(Message)
		assert.Equal(t, expectedLevelOne, capturedMessage.Level)

		capturedMessageTwo := mockImpl.Calls[1].Arguments.Get(1).(Message)
		assert.Equal(t, expectedLevelTwo, capturedMessageTwo.Level)
	})
}
