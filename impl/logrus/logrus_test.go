package logrus

import (
	"context"
	"github.com/shimmeringbee/logwrap"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
	"time"
)

func Test_mapLogLevels(t *testing.T) {
	t.Run("maps logwrap to logrus log levels", func(t *testing.T) {
		assert.Equal(t, logrus.PanicLevel, mapLogLevels(logwrap.Panic))
		assert.Equal(t, logrus.FatalLevel, mapLogLevels(logwrap.Fatal))
		assert.Equal(t, logrus.ErrorLevel, mapLogLevels(logwrap.Error))
		assert.Equal(t, logrus.WarnLevel, mapLogLevels(logwrap.Warn))
		assert.Equal(t, logrus.InfoLevel, mapLogLevels(logwrap.Info))
		assert.Equal(t, logrus.DebugLevel, mapLogLevels(logwrap.Debug))
		assert.Equal(t, logrus.TraceLevel, mapLogLevels(logwrap.Trace))

		// Unknown logwrap level.
		assert.Equal(t, logrus.InfoLevel, mapLogLevels(logwrap.LogLevel(math.MaxUint64)))
	})
}

func TestWrap(t *testing.T) {
	t.Run("wrap correctly sends log with message, level and fields", func(t *testing.T) {
		logger, hook := test.NewNullLogger()

		expectedLevel := logrus.WarnLevel
		expectedMessage := "message"
		expectedKey := "key"
		expectedValue := "value"
		expectedTime := time.Now()

		logrusWrap := Wrap(logger)
		logrusWrap(context.Background(), logwrap.Message{
			Level:     logwrap.Warn,
			Message:   expectedMessage,
			Data:      map[string]interface{}{expectedKey: expectedValue},
			Timestamp: expectedTime,
			Sequence:  0,
		})

		assert.Equal(t, 1, len(hook.Entries))
		entry := hook.LastEntry()

		assert.Equal(t, expectedLevel, entry.Level)
		assert.Equal(t, expectedMessage, entry.Message)
		assert.Equal(t, expectedTime, entry.Time)
		assert.Equal(t, expectedValue, entry.Data[expectedKey])
	})
}
