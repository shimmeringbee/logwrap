package impl

import (
	"context"
	"github.com/shimmeringbee/logwrap"
	"github.com/shimmeringbee/logwrap/impl/capture"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestWrap(t *testing.T) {
	t.Run("wrap correctly sends log with message, level and fields", func(t *testing.T) {
		captImpl := capture.NewCapture()
		logger := logwrap.New(captImpl.Impl())

		wrapImpl := Wrap(&logger)

		expectedLevel := logwrap.Warn
		expectedMessage := "message"
		expectedKey := "key"
		expectedValue := "value"
		expectedTime := time.Now()
		expectedSource := "source"

		wrapImpl(context.Background(), logwrap.Message{
			Level:     expectedLevel,
			Message:   expectedMessage,
			Data:      map[string]interface{}{expectedKey: expectedValue},
			Timestamp: expectedTime,
			Sequence:  0,
			Source:    expectedSource,
		})

		m := captImpl.Messages()
		assert.NotEmpty(t, m)

		entry := m[0]

		assert.Equal(t, expectedLevel, entry.Level)
		assert.Equal(t, expectedMessage, entry.Message)
		assert.Equal(t, expectedValue, entry.Data[expectedKey])
		assert.Equal(t, expectedSource, entry.Source)
	})
}
