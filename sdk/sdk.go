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

	"github.com/hashicorp/go-hclog"
)

type loggerKey string

var (
	rootLoggerKey loggerKey = "sdk"
)

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
