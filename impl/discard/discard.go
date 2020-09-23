package discard

import (
	"context"
	"github.com/shimmeringbee/logwrap"
)

// Discard implements a logger that discards all messages sent to it.
func Discard() logwrap.Impl {
	return func(ctx context.Context, message logwrap.Message) {}
}
