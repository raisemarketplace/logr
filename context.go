package logr

import "context"

const ContextKey = "logr.Logger"

// ContextWithLogger returns a new context with the provided logr.Logger added as a value in the context.
func ContextWithLogger(ctx context.Context, l Logger) context.Context {
	return context.WithValue(ctx, ContextKey, l)
}

// FromContext returns a logr.Logger from the given context if it has one,
// it will return the default logr.Logger otherwise.
func FromContext(ctx context.Context) Logger {
	if l, ok := ctx.Value(ContextKey).(Logger); ok {
		return l
	}
	return logr
}
