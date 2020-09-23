package golog

import (
	"bytes"
	"context"
	"github.com/shimmeringbee/logwrap"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func TestWrap(t *testing.T) {
	t.Run("wrap correctly sends log with message, level and fields", func(t *testing.T) {
		var outputBuffer bytes.Buffer
		goLogger := log.New(&outputBuffer, "", 0)
		gologWrap := Wrap(goLogger)

		logger := logwrap.New(gologWrap)

		logger.Log(context.Background(), "message with spaces", logwrap.Datum("key", "value"))

		expectedMessage := "[INFO] \"message with spaces\" {\"key\":\"value\"}\n"
		actualMessage := outputBuffer.String()
		assert.Equal(t, expectedMessage, actualMessage)
	})
}
