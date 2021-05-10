// Package tfsdklog provides helper functions for logging from SDKs and
// frameworks for building Terraform plugins.
//
// Plugin authors shouldn't need to use this package; it is meant for authors
// of the frameworks and SDKs for plugins. Plugin authors should use the tflog
// package.
//
// This package provides very similar functionality to tflog, except it uses a
// separate namespace for its logs.
package tfsdklog

import (
	"context"
	"io"
	"os"

	"github.com/hashicorp/go-hclog"
)

type loggerKey string

var (
	rootLoggerKey loggerKey = "sdk"
	stderr        io.Writer
)

func init() {
	// When go-plugin.Serve is called, it overwrites our os.Stderr with a
	// gRPC stream which Terraform ignores. This tends to be before our
	// loggers get set up, as go-plugin has no way to pass in a base
	// context, and our loggers are passed around via contexts. This leaves
	// our loggers writing to an output that is never read by anything,
	// meaning the logs get blackholed. This isn't ideal, for log output,
	// so this is our workaround: we copy stderr on init, before Serve can
	// be called, and offer an option to write to that instead of the
	// os.Stderr available at runtime.
	//
	// Ideally, this is a short-term fix until Terraform starts reading
	// from go-plugin's gRPC-streamed stderr channel, but for the moment it
	// works.
	stderr = os.Stderr
}

func getRootLogger(ctx context.Context) hclog.Logger {
	logger := ctx.Value(rootLoggerKey)
	if logger == nil {
		return nil
	}
	return logger.(hclog.Logger)
}

func setRootLogger(ctx context.Context, logger hclog.Logger) context.Context {
	return context.WithValue(ctx, rootLoggerKey, logger)
}

// New returns a new context.Context that contains a logger configured with the
// passed options.
func New(ctx context.Context, options ...Option) context.Context {
	opts := applyLoggerOpts(options...)
	if opts.level == hclog.NoLevel {
		opts.level = hclog.Trace
	}
	loggerOptions := &hclog.LoggerOptions{
		Name:              opts.name,
		Level:             opts.level,
		JSONFormat:        true,
		IndependentLevels: true,
		IncludeLocation:   opts.includeLocation,
		DisableTime:       !opts.includeTime,
		Output:            opts.output,
	}
	return setRootLogger(ctx, hclog.New(loggerOptions))
}

// With returns a new context.Context that has a modified logger in it which
// will include key and value as arguments in all its log output.
func With(ctx context.Context, key string, value interface{}) context.Context {
	logger := getRootLogger(ctx)
	if logger == nil {
		// this essentially should never happen in production the root
		// logger for  code should be injected by the  in
		// question, so really this is only likely in unit tests, at
		// most so just making this a no-op is fine
		return ctx
	}
	return setRootLogger(ctx, logger.With(key, value))
}

// Trace logs `msg` at the trace level to the logger in `ctx`, with `args` as
// structured arguments in the log output. `args` is expected to be pairs of
// key and value.
func Trace(ctx context.Context, msg string, args ...interface{}) {
	logger := getRootLogger(ctx)
	if logger == nil {
		// this essentially should never happen in production the root
		// logger for  code should be injected by the  in
		// question, so really this is only likely in unit tests, at
		// most so just making this a no-op is fine
		return
	}
	logger.Trace(msg, args...)
}

// Debug logs `msg` at the debug level to the logger in `ctx`, with `args` as
// structured arguments in the log output. `args` is expected to be pairs of
// key and value.
func Debug(ctx context.Context, msg string, args ...interface{}) {
	logger := getRootLogger(ctx)
	if logger == nil {
		// this essentially should never happen in production the root
		// logger for  code should be injected by the  in
		// question, so really this is only likely in unit tests, at
		// most so just making this a no-op is fine
		return
	}
	logger.Debug(msg, args...)
}

// Info logs `msg` at the info level to the logger in `ctx`, with `args` as
// structured arguments in the log output. `args` is expected to be pairs of
// key and value.
func Info(ctx context.Context, msg string, args ...interface{}) {
	logger := getRootLogger(ctx)
	if logger == nil {
		// this essentially should never happen in production the root
		// logger for  code should be injected by the  in
		// question, so really this is only likely in unit tests, at
		// most so just making this a no-op is fine
		return
	}
	logger.Info(msg, args...)
}

// Warn logs `msg` at the warn level to the logger in `ctx`, with `args` as
// structured arguments in the log output. `args` is expected to be pairs of
// key and value.
func Warn(ctx context.Context, msg string, args ...interface{}) {
	logger := getRootLogger(ctx)
	if logger == nil {
		// this essentially should never happen in production the root
		// logger for  code should be injected by the  in
		// question, so really this is only likely in unit tests, at
		// most so just making this a no-op is fine
		return
	}
	logger.Warn(msg, args...)
}

// Error logs `msg` at the error level to the logger in `ctx`, with `args` as
// structured arguments in the log output. `args` is expected to be pairs of
// key and value.
func Error(ctx context.Context, msg string, args ...interface{}) {
	logger := getRootLogger(ctx)
	if logger == nil {
		// this essentially should never happen in production the root
		// logger for  code should be injected by the  in
		// question, so really this is only likely in unit tests, at
		// most so just making this a no-op is fine
		return
	}
	logger.Error(msg, args...)
}
