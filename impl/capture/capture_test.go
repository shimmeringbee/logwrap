package capture

import (
	"context"
	"github.com/shimmeringbee/logwrap"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCapture(t *testing.T) {
	t.Run("returns a message that has been sent to the implementation", func(t *testing.T) {
		c := NewCapture()

		m := c.Messages()
		assert.Empty(t, m)

		expectedMessage := "debug message"

		c.Impl()(context.TODO(), logwrap.Message{
			Level:     logwrap.Debug,
			Message:   expectedMessage,
			Data:      make(map[string]interface{}),
			Timestamp: time.Time{},
			Sequence:  1,
		})

		m = c.messages
		assert.Len(t, m, 1)

		n := m[0]
		assert.Equal(t, expectedMessage, n.Message)
	})

	t.Run("clear wipes any captured messages", func(t *testing.T) {
		c := NewCapture()

		m := c.Messages()
		assert.Empty(t, m)

		c.Impl()(context.TODO(), logwrap.Message{})

		m = c.Messages()
		assert.NotEmpty(t, m)

		c.Clear()

		m = c.Messages()
		assert.Empty(t, m)
	})
}
