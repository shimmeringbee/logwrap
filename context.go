package logwrap

import "context"

type contextKey struct {
	base   string
	unique uint64
}

func (l Logger) contextKey(base string) contextKey {
	return contextKey{
		base:   base,
		unique: l.unique,
	}
}

// AddOptionsToContext add default Option's to a context for this specific logger (i.e. two loggers will have different
// options on the same context). These are always processed first, before any Option's provided during Log.
func (l Logger) AddOptionsToContext(ctx context.Context, options ...Option) context.Context {
	optionsToAdd := append(l.getOptionsFromContext(ctx), options...)
	return context.WithValue(ctx, l.contextKey(contextKeyOptions), optionsToAdd)
}

func (l Logger) getOptionsFromContext(ctx context.Context) []Option {
	uncast := ctx.Value(l.contextKey(contextKeyOptions))
	if uncast == nil {
		return []Option{}
	}

	if options, ok := uncast.([]Option); ok {
		return options
	}

	return []Option{}
}
