package filter

import (
	"context"
	"github.com/shimmeringbee/logwrap"
)

// Filter is an Implementation that permits filtering of messages outbound to an Implementation.
func Filter(impl logwrap.Impl, filter func(message logwrap.Message) bool) logwrap.Impl {
	return func(ctx context.Context, message logwrap.Message) {
		if filter(message) {
			impl(ctx, message)
		}
	}
}
