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
