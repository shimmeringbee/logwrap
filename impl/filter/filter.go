package filter

import (
	"context"
	"github.com/shimmeringbee/logwrap"
)

func Filter(impl logwrap.Impl, filter func(message logwrap.Message) bool) logwrap.Impl {
	return func(ctx context.Context, message logwrap.Message) {
		if filter(message) {
			impl(ctx, message)
		}
	}
}
