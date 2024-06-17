package supervisor

import (
	"context"
	"log/slog"
)

type ctxKey string

const (
	ctxKeyAttributes ctxKey = "attributes"
	ctxKeyChildName  ctxKey = "group.name"
)

type attributes struct {
	values map[string][]slog.Attr
}

func (a attributes) Set(key string, attr slog.Attr) {
	values, ok := a.values[key]
	if !ok {
		values = []slog.Attr{attr}
	} else {
		values = append(values, attr)
	}
	a.values[key] = values
}

func newAttributes() attributes {
	return attributes{
		values: map[string][]slog.Attr{},
	}
}

type ContextSupervisor struct {
	slog.Handler
}

func (h ContextSupervisor) Handle(ctx context.Context, r slog.Record) {
	h.Handler.Handle(ctx, r)
}

func contextAttributes(ctx context.Context) (attributes, bool) {
	v, ok := ctx.Value(ctxKeyAttributes).(attributes)
	return v, ok
}

func contextInit(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxKeyAttributes, newAttributes())
}

func NewChild(ctx context.Context, name string) context.Context {
	_, ok := contextAttributes(ctx)
	if !ok {
		ctx = contextInit(ctx)
	}

	ctx = context.WithValue(ctx, ctxKeyChildName, name)
	return ctx
}

func ContextAddAttributes(ctx context.Context, attr slog.Attr) {
	attrs, ok := contextAttributes(ctx)
	if !ok {
		attrs = newAttributes()
	}

	attrs.Set(attr.Key, attr)
}
