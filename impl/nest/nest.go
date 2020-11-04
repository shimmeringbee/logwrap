package nest

import (
	"context"
	"github.com/shimmeringbee/logwrap"
)

// Nest is a wrapper around logwrap, allows passing messages back to a parent implementation.
func Wrap(dest *logwrap.Logger) logwrap.Impl {
	return func(ctx context.Context, message logwrap.Message) {
		dest.Log(ctx, message.Message, logwrap.Data(message.Data), logwrap.Level(message.Level), logwrap.Source(message.Source))
	}
}
