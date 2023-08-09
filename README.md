# Go `log/slog` middleware

The `slogmw` package makes it easier to use the standard library `log/slog` package. 

This package provides middleware for adding, changing, or editing fields of the standard library structured logger.


The slog package, and the slogmw package comes with two hooks:
- `slog.HandlerOptions` has a `ReplaceAttr` function. We use `slogwm.FormatChain(...)` to create a chain of formatters that edit attributes using this hook.
- `slog.Handler`. You can use `slowmw.WrapHandler(...)` to wrap a handler and edit the log events.

## Examples

Format the time field of the log events: 
```go
opts := &slog.HandlerOptions{
  // Create a chain of slog formatters
  ReplaceAttr: slogmw.FormatChain(
    // pass the key to edit "slog.TimeKey" and 
    // the date time format to use "time.DateTime"
    slogmw.FormatTime(slog.TimeKey, time.DateTime)
  ),
}

slog.SetDefault(slog.NewJSONHandler(os.Stdout, opts))
```

Change the key of a log field:
```go
opts := &slog.HandlerOptions{
  // Create a chain of slog formatters
  ReplaceAttr: slogmw.FormatChain(
    // The default log event message key is "msg"
    // this is exported as the constant slog.MessageKey.
    // Here we can change the msg key to "log". All
    // log events will print out the log field as "log"
    // and not "msg".
    slog.FormatKey(slog.MessageKey, "log")
  ),
}

slog.SetDefault(slog.NewJSONHandler(os.Stdout, opts))
slog.Info("hello world")

// Output: {"level": "INFO", "log": "hello world" ...trimmed}
```

Include a field from the `context.Context` in every log event:
```go
// create a function that can extract the
// attributes from the context passed.
ctxValueFn := func(ctx context.Context) []slog.Attr {
  v := ctx.Value("userid").(string)
  return []slog.Attr{slog.String("user_id", v)}
}

h := slogmw.WrapHandler(
  // Pass the built in JSONHandler to be wrapped
  slog.NewJSONHandler(os.Stdout, opts),
  // Pass the IncludeContext middleware to extract
  // values from the context and include them as log fields.
  slogmw.IncludeContext(ctxValueFn),
)

slog.SetDefault(h)

// The context passed to InfoContext will be
// available to the function in slogmw.IncludeContext. 
// If just slog.Info(...) is called, the context will
// be context.Background()
slog.InfoContext(ctx, "hello world")

// Output: {"level": "INFO", "msg": "hello world", "user_id": "..."}
```

Include a static field in every log event:
```go
h := slogmw.WrapHandler(
  // Pass the built in JSONHandler to be wrapped
  slog.NewJSONHandler(os.Stdout, opts),
  // Include env=prod in every log event
  slogmw.IncludeStatic(slog.String("env", "prod"))
)

slog.SetDefault(h)

slog.Info("hello world")
// Output: {"level": "INFO", "msg": "hello world", "env": "prod" .. trimmed}
```

More example usage can be found in the [`middleware_test.go`](https://github.com/zknill/slogmw/blob/main/middleware_test.go) file.
