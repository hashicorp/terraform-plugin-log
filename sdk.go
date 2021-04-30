package tflog

import (
	"context"

	"github.com/hashicorp/go-hclog"
)

func getSDKRootLogger(ctx context.Context) hclog.Logger {
	logger := ctx.Value(sdkRootLoggerKey)
	if logger == nil {
		return nil
	}
	return logger.(hclog.Logger)
}

func setSDKRootLogger(ctx context.Context, logger hclog.Logger) context.Context {
	return context.WithValue(ctx, sdkRootLoggerKey, logger)
}

func getSDKSubsystemLogger(ctx context.Context, subsystem string) hclog.Logger {
	logger := ctx.Value(sdkRootLoggerKey + loggerKey("."+subsystem))
	if logger == nil {
		return nil
	}
	return logger.(hclog.Logger)
}

func setSDKSubsystemLogger(ctx context.Context, subsystem string, logger hclog.Logger) context.Context {
	return context.WithValue(ctx, sdkRootLoggerKey+loggerKey("."+subsystem), logger)
}

func NewSDK(ctx context.Context, options ...Option) context.Context {
	opts := applyLoggerOpts(options...)
	return setSDKRootLogger(ctx, hclog.New(&hclog.LoggerOptions{
		Name:              "sdk",
		Level:             opts.level,
		JSONFormat:        true,
		IndependentLevels: true,
		IncludeLocation:   opts.includeLocation,
	}))
}

func SDKWith(ctx context.Context, key string, value interface{}) context.Context {
	logger := getSDKRootLogger(ctx)
	if logger == nil {
		// this essentially should never happen in production the root
		// logger for SDK code should be injected by the SDK in
		// question, so really this is only likely in unit tests, at
		// most so just making this a no-op is fine
		return ctx
	}
	return setSDKRootLogger(ctx, logger.With(key, value))
}

func SDKTrace(ctx context.Context, msg string, args ...interface{}) {
	logger := getSDKRootLogger(ctx)
	if logger == nil {
		// this essentially should never happen in production the root
		// logger for SDK code should be injected by the SDK in
		// question, so really this is only likely in unit tests, at
		// most so just making this a no-op is fine
		return
	}
	logger.Trace(msg, args...)
}

func SDKDebug(ctx context.Context, msg string, args ...interface{}) {
	logger := getSDKRootLogger(ctx)
	if logger == nil {
		// this essentially should never happen in production the root
		// logger for SDK code should be injected by the SDK in
		// question, so really this is only likely in unit tests, at
		// most so just making this a no-op is fine
		return
	}
	logger.Debug(msg, args...)
}

func SDKInfo(ctx context.Context, msg string, args ...interface{}) {
	logger := getSDKRootLogger(ctx)
	if logger == nil {
		// this essentially should never happen in production the root
		// logger for SDK code should be injected by the SDK in
		// question, so really this is only likely in unit tests, at
		// most so just making this a no-op is fine
		return
	}
	logger.Info(msg, args...)
}

func SDKWarn(ctx context.Context, msg string, args ...interface{}) {
	logger := getSDKRootLogger(ctx)
	if logger == nil {
		// this essentially should never happen in production the root
		// logger for SDK code should be injected by the SDK in
		// question, so really this is only likely in unit tests, at
		// most so just making this a no-op is fine
		return
	}
	logger.Warn(msg, args...)
}

func SDKError(ctx context.Context, msg string, args ...interface{}) {
	logger := getSDKRootLogger(ctx)
	if logger == nil {
		// this essentially should never happen in production the root
		// logger for SDK code should be injected by the SDK in
		// question, so really this is only likely in unit tests, at
		// most so just making this a no-op is fine
		return
	}
	logger.Error(msg, args...)
}

func NewSDKSubsystem(ctx context.Context, subsystem string, options ...Option) context.Context {
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
	return setSDKSubsystemLogger(ctx, subsystem, logger.Named(subsystem))
}

func SDKSubsystemWith(ctx context.Context, subsystem, key string, value interface{}) context.Context {
	logger := getSDKSubsystemLogger(ctx, subsystem)
	if logger == nil {
		// create a new logger if one doesn't exist
		logger = getSDKSubsystemLogger(NewSDKSubsystem(ctx, subsystem), subsystem)
	}
	return setSDKSubsystemLogger(ctx, subsystem, logger.With(key, value))
}

func SDKSubsystemTrace(ctx context.Context, subsystem, msg string, args ...interface{}) {
	logger := getSDKSubsystemLogger(ctx, subsystem)
	if logger == nil {
		// create a new logger if one doesn't exist
		logger = getSDKSubsystemLogger(NewSDKSubsystem(ctx, subsystem), subsystem)
	}
	logger.Trace(msg, args...)
}

func SDKSubsystemDebug(ctx context.Context, subsystem, msg string, args ...interface{}) {
	logger := getSDKSubsystemLogger(ctx, subsystem)
	if logger == nil {
		// create a new logger if one doesn't exist
		logger = getSDKSubsystemLogger(NewSDKSubsystem(ctx, subsystem), subsystem)
	}
	logger.Debug(msg, args...)
}

func SDKSubsystemInfo(ctx context.Context, subsystem, msg string, args ...interface{}) {
	logger := getSDKSubsystemLogger(ctx, subsystem)
	if logger == nil {
		// create a new logger if one doesn't exist
		logger = getSDKSubsystemLogger(NewSDKSubsystem(ctx, subsystem), subsystem)
	}
	logger.Info(msg, args...)
}

func SDKSubsystemWarn(ctx context.Context, subsystem, msg string, args ...interface{}) {
	logger := getSDKSubsystemLogger(ctx, subsystem)
	if logger == nil {
		// create a new logger if one doesn't exist
		logger = getSDKSubsystemLogger(NewSDKSubsystem(ctx, subsystem), subsystem)
	}
	logger.Warn(msg, args...)
}

func SDKSubsystemError(ctx context.Context, subsystem, msg string, args ...interface{}) {
	logger := getSDKSubsystemLogger(ctx, subsystem)
	if logger == nil {
		// create a new logger if one doesn't exist
		logger = getSDKSubsystemLogger(NewSDKSubsystem(ctx, subsystem), subsystem)
	}
	logger.Error(msg, args...)
}
