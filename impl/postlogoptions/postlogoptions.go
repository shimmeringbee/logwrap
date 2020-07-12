package postlogoptions

import (
	"context"
	"github.com/shimmeringbee/logwrap"
)

// PostLogOptions is an implementation that has the ability to modify a log message before being sent onwards to another
// implementation.
func PostLogOptions(impl logwrap.Impl, options ...logwrap.Option) logwrap.Impl {
	return func(ctx context.Context, message logwrap.Message) {
		for _, option := range options {
			option(&message)
		}

		impl(ctx, message)
	}
}
