package tflog

import (
	"context"

	"github.com/hashicorp/go-hclog"
)

func getProviderSpaceRootLogger(ctx context.Context) hclog.Logger {
	logger := ctx.Value(providerSpaceRootLoggerKey)
	if logger == nil {
		return nil
	}
	return logger.(hclog.Logger)
}

func setProviderSpaceRootLogger(ctx context.Context, logger hclog.Logger) context.Context {
	return context.WithValue(ctx, providerSpaceRootLoggerKey, logger)
}

// New returns a new context.Context that contains a logger configured with the
// passed options.
func New(ctx context.Context, options ...Option) context.Context {
	opts := applyLoggerOpts(options...)
	return setProviderSpaceRootLogger(ctx, hclog.New(&hclog.LoggerOptions{
		Name:              opts.name,
		Level:             opts.level,
		JSONFormat:        true,
		IndependentLevels: true,
		IncludeLocation:   opts.includeLocation,
	}))
}

// With returns a new context.Context that has a modified logger in it which
// will include key and value as arguments in all its log output.
func With(ctx context.Context, key string, value interface{}) context.Context {
	logger := getProviderSpaceRootLogger(ctx)
	if logger == nil {
		// this essentially should never happen in production
		// the root logger for provider code should be injected
		// by whatever SDK the provider developer is using, so
		// really this is only likely in unit tests, at most
		// so just making this a no-op is fine
		return ctx
	}
	return setProviderSpaceRootLogger(ctx, logger.With(key, value))
}

// Trace logs `msg` at the trace level to the logger in `ctx`, with `args` as
// structured arguments in the log output. `args` is expected to be pairs of
// key and value, so Trace(ctx, "hello, world", "foo", 123) would have
// "foo=123" in its arguments.
func Trace(ctx context.Context, msg string, args ...interface{}) {
	logger := getProviderSpaceRootLogger(ctx)
	if logger == nil {
		// this essentially should never happen in production
		// the root logger for provider code should be injected
		// by whatever SDK the provider developer is using, so
		// really this is only likely in unit tests, at most
		// so just making this a no-op is fine
		return
	}
	logger.Trace(msg, args...)
}

// Debug logs `msg` at the debug level to the logger in `ctx`, with `args` as
// structured arguments in the log output. `args` is expected to be pairs of
// key and value, so Debug(ctx, "hello, world", "foo", 123) would have
// "foo=123" in its arguments.
func Debug(ctx context.Context, msg string, args ...interface{}) {
	logger := getProviderSpaceRootLogger(ctx)
	if logger == nil {
		// this essentially should never happen in production
		// the root logger for provider code should be injected
		// by whatever SDK the provider developer is using, so
		// really this is only likely in unit tests, at most
		// so just making this a no-op is fine
		return
	}
	logger.Debug(msg, args...)
}

// Info logs `msg` at the info level to the logger in `ctx`, with `args` as
// structured arguments in the log output. `args` is expected to be pairs of
// key and value, so Info(ctx, "hello, world", "foo", 123) would have
// "foo=123" in its arguments.
func Info(ctx context.Context, msg string, args ...interface{}) {
	logger := getProviderSpaceRootLogger(ctx)
	if logger == nil {
		// this essentially should never happen in production
		// the root logger for provider code should be injected
		// by whatever SDK the provider developer is using, so
		// really this is only likely in unit tests, at most
		// so just making this a no-op is fine
		return
	}
	logger.Info(msg, args...)
}

// Warn logs `msg` at the warn level to the logger in `ctx`, with `args` as
// structured arguments in the log output. `args` is expected to be pairs of
// key and value, so Warn(ctx, "hello, world", "foo", 123) would have
// "foo=123" in its arguments.
func Warn(ctx context.Context, msg string, args ...interface{}) {
	logger := getProviderSpaceRootLogger(ctx)
	if logger == nil {
		// this essentially should never happen in production
		// the root logger for provider code should be injected
		// by whatever SDK the provider developer is using, so
		// really this is only likely in unit tests, at most
		// so just making this a no-op is fine
		return
	}
	logger.Warn(msg, args...)
}

// Error logs `msg` at the error level to the logger in `ctx`, with `args` as
// structured arguments in the log output. `args` is expected to be pairs of
// key and value, so Error(ctx, "hello, world", "foo", 123) would have
// "foo=123" in its arguments.
func Error(ctx context.Context, msg string, args ...interface{}) {
	logger := getProviderSpaceRootLogger(ctx)
	if logger == nil {
		// this essentially should never happen in production
		// the root logger for provider code should be injected
		// by whatever SDK the provider developer is using, so
		// really this is only likely in unit tests, at most
		// so just making this a no-op is fine
		return
	}
	logger.Error(msg, args...)
}
