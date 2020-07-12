package tee

import (
	"context"
	"github.com/shimmeringbee/logwrap"
)

// Tee is an implementation that distributes log messages to multiple different implementations. Calls are made
// sequentially.
func Tee(destinations ...logwrap.Impl) logwrap.Impl {
	return func(ctx context.Context, message logwrap.Message) {
		for _, impl := range destinations {
			impl(ctx, message)
		}
	}
}
