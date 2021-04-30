package tflog

import (
	"context"

	"github.com/hashicorp/go-hclog"
)

func getSubsystemLogger(ctx context.Context, subsystem string) hclog.Logger {
	logger := ctx.Value(providerSpaceRootLoggerKey + loggerKey("."+subsystem))
	if logger == nil {
		return nil
	}
	return logger.(hclog.Logger)
}

func setSubsystemLogger(ctx context.Context, subsystem string, logger hclog.Logger) context.Context {
	return context.WithValue(ctx, providerSpaceRootLoggerKey+loggerKey("."+subsystem), logger)
}

func NewSubsystem(ctx context.Context, subsystem string, options ...Option) context.Context {
	logger := getProviderSpaceRootLogger(ctx)
	if logger == nil {
		// this essentially should never happen in production
		// the root logger for provider code should be injected
		// by whatever SDK the provider developer is using, so
		// really this is only likely in unit tests, at most
		// so just making this a no-op is fine
		return ctx
	}
	subLogger := logger.Named(subsystem)
	opts := applyLoggerOpts(options...)
	if opts.level != hclog.NoLevel {
		subLogger.SetLevel(opts.level)
	}
	return setSubsystemLogger(ctx, subsystem, logger.Named(subsystem))
}

func SubsystemWith(ctx context.Context, subsystem, key string, value interface{}) context.Context {
	logger := getSubsystemLogger(ctx, subsystem)
	if logger == nil {
		// create a new logger if one doesn't exist
		logger = getSubsystemLogger(NewSubsystem(ctx, subsystem), subsystem)
	}
	return setSubsystemLogger(ctx, subsystem, logger.With(key, value))
}

func SubsystemTrace(ctx context.Context, subsystem, msg string, args ...interface{}) {
	logger := getSubsystemLogger(ctx, subsystem)
	if logger == nil {
		// create a new logger if one doesn't exist
		logger = getSubsystemLogger(NewSubsystem(ctx, subsystem), subsystem)
	}
	logger.Trace(msg, args...)
}

func SubsystemDebug(ctx context.Context, subsystem, msg string, args ...interface{}) {
	logger := getSubsystemLogger(ctx, subsystem)
	if logger == nil {
		// create a new logger if one doesn't exist
		logger = getSubsystemLogger(NewSubsystem(ctx, subsystem), subsystem)
	}
	logger.Debug(msg, args...)
}

func SubsystemInfo(ctx context.Context, subsystem, msg string, args ...interface{}) {
	logger := getSubsystemLogger(ctx, subsystem)
	if logger == nil {
		// create a new logger if one doesn't exist
		logger = getSubsystemLogger(NewSubsystem(ctx, subsystem), subsystem)
	}
	logger.Info(msg, args...)
}

func SubsystemWarn(ctx context.Context, subsystem, msg string, args ...interface{}) {
	logger := getSubsystemLogger(ctx, subsystem)
	if logger == nil {
		// create a new logger if one doesn't exist
		logger = getSubsystemLogger(NewSubsystem(ctx, subsystem), subsystem)
	}
	logger.Warn(msg, args...)
}

func SubsystemError(ctx context.Context, subsystem, msg string, args ...interface{}) {
	logger := getSubsystemLogger(ctx, subsystem)
	if logger == nil {
		// create a new logger if one doesn't exist
		logger = getSubsystemLogger(NewSubsystem(ctx, subsystem), subsystem)
	}
	logger.Error(msg, args...)
}
