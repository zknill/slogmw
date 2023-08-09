package slogmw_test

import (
	"bytes"
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zknill/slogmw"
)

func TestFormat(t *testing.T) {
	t.Parallel()

	now, _ := time.Parse(time.DateTime, "2023-08-09 12:12:04")

	opts := &slog.HandlerOptions{
		ReplaceAttr: slogmw.FormatChain(
			staticTime(now),
			slogmw.FormatTime(slog.TimeKey, time.DateTime),
			slogmw.FormatField(slog.MessageKey, func(a slog.Attr) slog.Attr {
				a.Value = slog.StringValue(a.Value.String() + " edited")
				return a
			}),
			slogmw.FormatKey(slog.MessageKey, "log"),
		),
	}

	b := bytes.Buffer{}

	logger := slog.New(slog.NewTextHandler(&b, opts))

	logger.Info("hello world")

	want := "time=\"2023-08-09 12:12:04\" level=INFO log=\"hello world edited\"\n"
	assert.Equal(t, want, b.String())

}

func staticTime(t time.Time) slogmw.AttrFormatFunc {
	return func(a slog.Attr) slog.Attr {
		if a.Key != slog.TimeKey {
			return a
		}

		a.Value = slog.TimeValue(t)

		return a
	}
}

func TestWrapHandler(t *testing.T) {
	t.Parallel()

	ctx := context.WithValue(context.Background(), "key", "my-value")

	now, _ := time.Parse(time.DateTime, "2023-08-09 12:12:04")
	opts := &slog.HandlerOptions{
		ReplaceAttr: slogmw.FormatChain(
			staticTime(now),
			slogmw.FormatTime(slog.TimeKey, time.DateTime),
		),
	}

	b := bytes.Buffer{}
	h := slogmw.WrapHandler(
		slog.NewJSONHandler(&b, opts),
		slogmw.IncludeStatic(slog.String("a", "b")),
		slogmw.IncludeContext(func(ctx context.Context) []slog.Attr {
			v := ctx.Value("key").(string)
			return []slog.Attr{slog.String("my-key", v)}
		}),
	)

	logger := slog.New(h)

	logger.InfoContext(ctx, "hello world")

	want := `{
		"time":"2023-08-09 12:12:04", 
		"a": "b", 
		"level": "INFO", 
		"msg": "hello world", 
		"my-key": "my-value"
	}`

	assert.JSONEq(t, want, b.String())
}
