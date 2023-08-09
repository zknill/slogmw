package slogmw

import (
	"context"
	"log/slog"
)

type AttrFormatFunc func(a slog.Attr) slog.Attr

func FormatChain(fns ...AttrFormatFunc) func(groups []string, a slog.Attr) slog.Attr {
	return func(_ []string, a slog.Attr) slog.Attr {
		for _, f := range fns {
			a = f(a)
		}

		return a
	}
}

func FormatTime(key string, format string) func(a slog.Attr) slog.Attr {
	return func(a slog.Attr) slog.Attr {
		if a.Key != key {
			return a
		}

		t := a.Value.Time()

		a.Value = slog.StringValue(t.Format(format))

		return a
	}
}

func FormatField(key string, f AttrFormatFunc) AttrFormatFunc {
	return func(a slog.Attr) slog.Attr {
		if a.Key != key {
			return a
		}

		return f(a)
	}
}

func FormatKey(key, newKey string) AttrFormatFunc {
	return func(a slog.Attr) slog.Attr {
		if a.Key != key {
			return a
		}

		a.Key = newKey

		return a
	}
}

type wrapper struct {
	slog.Handler
	attrFns []AttrValueFunc
}

type AttrValueFunc func(ctx context.Context, r slog.Record) []slog.Attr

func (w *wrapper) Handle(ctx context.Context, r slog.Record) error {
	var extra []slog.Attr

	for _, f := range w.attrFns {
		extra = append(extra, f(ctx, r)...)
	}

	return w.
		Handler.
		WithAttrs(extra).
		Handle(ctx, r)
}

func WrapHandler(handler slog.Handler, fns ...AttrValueFunc) slog.Handler {
	return &wrapper{
		Handler: handler,
		attrFns: fns,
	}
}

func IncludeContext(fieldExtractor func(ctx context.Context) []slog.Attr) AttrValueFunc {
	return func(ctx context.Context, r slog.Record) []slog.Attr {
		return fieldExtractor(ctx)
	}
}

func IncludeStatic(static ...slog.Attr) AttrValueFunc {
	return func(_ context.Context, _ slog.Record) []slog.Attr {
		return static
	}
}
