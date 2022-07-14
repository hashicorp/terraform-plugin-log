package tfsdklog

import (
	"context"
	"regexp"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-log/internal/hclogutils"
	"github.com/hashicorp/terraform-plugin-log/internal/logging"
)

// NewRootSDKLogger returns a new context.Context that contains an SDK logger
// configured with the passed options.
func NewRootSDKLogger(ctx context.Context, options ...logging.Option) context.Context {
	opts := logging.ApplyLoggerOpts(options...)
	if opts.Name == "" {
		opts.Name = logging.DefaultSDKRootLoggerName
	}
	if sink := logging.GetSink(ctx); sink != nil {
		logger := sink.Named(opts.Name)
		sinkLoggerOptions := logging.GetSinkOptions(ctx)
		sdkLoggerOptions := hclogutils.LoggerOptionsCopy(sinkLoggerOptions)
		sdkLoggerOptions.Name = opts.Name

		if opts.Level != hclog.NoLevel {
			logger.SetLevel(opts.Level)
			sdkLoggerOptions.Level = opts.Level
		}

		ctx = logging.SetSDKRootLogger(ctx, logger)
		ctx = logging.SetSDKRootLoggerOptions(ctx, sdkLoggerOptions)

		return ctx
	}
	if opts.Level == hclog.NoLevel {
		opts.Level = hclog.Trace
	}
	loggerOptions := &hclog.LoggerOptions{
		Name:                     opts.Name,
		Level:                    opts.Level,
		JSONFormat:               true,
		IndependentLevels:        true,
		IncludeLocation:          opts.IncludeLocation,
		DisableTime:              !opts.IncludeTime,
		Output:                   opts.Output,
		AdditionalLocationOffset: opts.AdditionalLocationOffset,
	}

	ctx = logging.SetSDKRootLogger(ctx, hclog.New(loggerOptions))
	ctx = logging.SetSDKRootLoggerOptions(ctx, loggerOptions)

	return ctx
}

// NewRootProviderLogger returns a new context.Context that contains a provider
// logger configured with the passed options.
func NewRootProviderLogger(ctx context.Context, options ...logging.Option) context.Context {
	opts := logging.ApplyLoggerOpts(options...)
	if opts.Name == "" {
		opts.Name = logging.DefaultProviderRootLoggerName
	}
	if sink := logging.GetSink(ctx); sink != nil {
		logger := sink.Named(opts.Name)
		sinkLoggerOptions := logging.GetSinkOptions(ctx)
		providerLoggerOptions := hclogutils.LoggerOptionsCopy(sinkLoggerOptions)
		providerLoggerOptions.Name = opts.Name

		if opts.Level != hclog.NoLevel {
			logger.SetLevel(opts.Level)
			providerLoggerOptions.Level = opts.Level
		}

		ctx = logging.SetProviderRootLogger(ctx, logger)
		ctx = logging.SetProviderRootLoggerOptions(ctx, providerLoggerOptions)

		return ctx
	}
	if opts.Level == hclog.NoLevel {
		opts.Level = hclog.Trace
	}
	loggerOptions := &hclog.LoggerOptions{
		Name:                     opts.Name,
		Level:                    opts.Level,
		JSONFormat:               true,
		IndependentLevels:        true,
		IncludeLocation:          opts.IncludeLocation,
		DisableTime:              !opts.IncludeTime,
		Output:                   opts.Output,
		AdditionalLocationOffset: opts.AdditionalLocationOffset,
	}

	ctx = logging.SetProviderRootLogger(ctx, hclog.New(loggerOptions))
	ctx = logging.SetProviderRootLoggerOptions(ctx, loggerOptions)

	return ctx
}

// With returns a new context.Context that has a modified logger in it which
func With(ctx context.Context, key string, value interface{}) context.Context {
// will include key and value as fields in all its log output.
	logger := logging.GetSDKRootLogger(ctx)
	if logger == nil {
		// this essentially should never happen in production the root
		// logger for  code should be injected by the  in
		// question, so really this is only likely in unit tests, at
		// most so just making this a no-op is fine
		return ctx
	}
	return logging.SetSDKRootLogger(ctx, logger.With(key, value))
}

// Trace logs `msg` at the trace level to the logger in `ctx`, with optional
// `additionalFields` structured key-value fields in the log output. Fields are
// shallow merged with any defined on the logger, e.g. by the `With()` function,
// and across multiple maps.
func Trace(ctx context.Context, msg string, additionalFields ...map[string]interface{}) {
	logger := logging.GetSDKRootLogger(ctx)
	if logger == nil {
		// this essentially should never happen in production the root
		// logger for  code should be injected by the  in
		// question, so really this is only likely in unit tests, at
		// most so just making this a no-op is fine
		return
	}

	additionalArgs, shouldOmit := omitOrMask(ctx, logger, &msg, additionalFields)
	if shouldOmit {
		return
	}

	logger.Trace(msg, additionalArgs...)
}

// Debug logs `msg` at the debug level to the logger in `ctx`, with optional
// `additionalFields` structured key-value fields in the log output. Fields are
// shallow merged with any defined on the logger, e.g. by the `With()` function,
// and across multiple maps.
func Debug(ctx context.Context, msg string, additionalFields ...map[string]interface{}) {
	logger := logging.GetSDKRootLogger(ctx)
	if logger == nil {
		// this essentially should never happen in production the root
		// logger for  code should be injected by the  in
		// question, so really this is only likely in unit tests, at
		// most so just making this a no-op is fine
		return
	}

	additionalArgs, shouldOmit := omitOrMask(ctx, logger, &msg, additionalFields)
	if shouldOmit {
		return
	}

	logger.Debug(msg, additionalArgs...)
}

// Info logs `msg` at the info level to the logger in `ctx`, with optional
// `additionalFields` structured key-value fields in the log output. Fields are
// shallow merged with any defined on the logger, e.g. by the `With()` function,
// and across multiple maps.
func Info(ctx context.Context, msg string, additionalFields ...map[string]interface{}) {
	logger := logging.GetSDKRootLogger(ctx)
	if logger == nil {
		// this essentially should never happen in production the root
		// logger for  code should be injected by the  in
		// question, so really this is only likely in unit tests, at
		// most so just making this a no-op is fine
		return
	}

	additionalArgs, shouldOmit := omitOrMask(ctx, logger, &msg, additionalFields)
	if shouldOmit {
		return
	}

	logger.Info(msg, additionalArgs...)
}

// Warn logs `msg` at the warn level to the logger in `ctx`, with optional
// `additionalFields` structured key-value fields in the log output. Fields are
// shallow merged with any defined on the logger, e.g. by the `With()` function,
// and across multiple maps.
func Warn(ctx context.Context, msg string, additionalFields ...map[string]interface{}) {
	logger := logging.GetSDKRootLogger(ctx)
	if logger == nil {
		// this essentially should never happen in production the root
		// logger for  code should be injected by the  in
		// question, so really this is only likely in unit tests, at
		// most so just making this a no-op is fine
		return
	}

	additionalArgs, shouldOmit := omitOrMask(ctx, logger, &msg, additionalFields)
	if shouldOmit {
		return
	}

	logger.Warn(msg, additionalArgs...)
}

// Error logs `msg` at the error level to the logger in `ctx`, with optional
// `additionalFields` structured key-value fields in the log output. Fields are
// shallow merged with any defined on the logger, e.g. by the `With()` function,
// and across multiple maps.
func Error(ctx context.Context, msg string, additionalFields ...map[string]interface{}) {
	logger := logging.GetSDKRootLogger(ctx)
	if logger == nil {
		// this essentially should never happen in production the root
		// logger for  code should be injected by the  in
		// question, so really this is only likely in unit tests, at
		// most so just making this a no-op is fine
		return
	}

	additionalArgs, shouldOmit := omitOrMask(ctx, logger, &msg, additionalFields)
	if shouldOmit {
		return
	}

	logger.Error(msg, additionalArgs...)
}

func omitOrMask(ctx context.Context, logger hclog.Logger, msg *string, additionalFields []map[string]interface{}) ([]interface{}, bool) {
	tfLoggerOpts := logging.GetSDKRootTFLoggerOpts(ctx)
	additionalArgs := hclogutils.MapsToArgs(additionalFields...)
	impliedArgs := logger.ImpliedArgs()

	// Apply the provider root LoggerOpts to determine if this log should be omitted
	if tfLoggerOpts.ShouldOmit(msg, impliedArgs, additionalArgs) {
		return nil, true
	}

	// Apply the provider root LoggerOpts to apply masking to this log
	tfLoggerOpts.ApplyMask(msg, impliedArgs, additionalArgs)
	return additionalArgs, false
}

// WithOmitLogWithFieldKeys returns a new context.Context that has a modified logger
// that will omit to write any log when any of the given keys is found
// within its fields.
//
// Each call to this function is additive:
// the keys to omit by are added to the existing configuration.
//
// Example:
//
//   configuration = `['foo', 'baz']`
//
//   log1 = `{ msg = "...", fields = { 'foo', '...', 'bar', '...' }`   -> omitted
//   log2 = `{ msg = "...", fields = { 'bar', '...' }`                 -> printed
//   log3 = `{ msg = "...", fields = { 'baz`', '...', 'boo', '...' }`  -> omitted
//
func WithOmitLogWithFieldKeys(ctx context.Context, keys ...string) context.Context {
	lOpts := logging.GetSDKRootTFLoggerOpts(ctx)

	lOpts = logging.WithOmitLogWithFieldKeys(keys...)(lOpts)

	return logging.SetSDKRootTFLoggerOpts(ctx, lOpts)
}

// WithOmitLogWithMessageRegex returns a new context.Context that has a modified logger
// that will omit to write any log that has a message matching any of the
// given *regexp.Regexp.
//
// Each call to this function is additive:
// the regexp to omit by are added to the existing configuration.
//
// Example:
//
//   configuration = `[regexp.MustCompile("(foo|bar)")]`
//
//   log1 = `{ msg = "banana apple foo", fields = {...}`     -> omitted
//   log2 = `{ msg = "pineapple mango", fields = {...}`      -> printed
//   log3 = `{ msg = "pineapple mango bar", fields = {...}`  -> omitted
//
func WithOmitLogWithMessageRegex(ctx context.Context, expressions ...*regexp.Regexp) context.Context {
	lOpts := logging.GetSDKRootTFLoggerOpts(ctx)

	lOpts = logging.WithOmitLogWithMessageRegex(expressions...)(lOpts)

	return logging.SetSDKRootTFLoggerOpts(ctx, lOpts)
}

// WithOmitLogMatchingString  returns a new context.Context that has a modified logger
// that will omit to write any log that matches any of the given string.
//
// Each call to this function is additive:
// the string to omit by are added to the existing configuration.
//
// Example:
//
//   configuration = `['foo', 'bar']`
//
//   log1 = `{ msg = "banana apple foo", fields = {...}`     -> omitted
//   log2 = `{ msg = "pineapple mango", fields = {...}`      -> printed
//   log3 = `{ msg = "pineapple mango bar", fields = {...}`  -> omitted
//
func WithOmitLogMatchingString(ctx context.Context, matchingStrings ...string) context.Context {
	lOpts := logging.GetSDKRootTFLoggerOpts(ctx)

	lOpts = logging.WithOmitLogWithMessageStrings(matchingStrings...)(lOpts)

	return logging.SetSDKRootTFLoggerOpts(ctx, lOpts)
}

// WithMaskFieldValueWithFieldKeys returns a new context.Context that has a modified logger
// that masks (replaces) with asterisks (`***`) any field value where the
// key matches one of the given keys.
//
// Each call to this function is additive:
// the keys to mask by are added to the existing configuration.
//
// Example:
//
//   configuration = `['foo', 'baz']`
//
//   log1 = `{ msg = "...", fields = { 'foo', '***', 'bar', '...' }`   -> masked value
//   log2 = `{ msg = "...", fields = { 'bar', '...' }`                 -> as-is value
//   log3 = `{ msg = "...", fields = { 'baz`', '***', 'boo', '...' }`  -> masked value
//
func WithMaskFieldValueWithFieldKeys(ctx context.Context, keys ...string) context.Context {
	lOpts := logging.GetSDKRootTFLoggerOpts(ctx)

	lOpts = logging.WithMaskFieldValueWithFieldKeys(keys...)(lOpts)

	return logging.SetSDKRootTFLoggerOpts(ctx, lOpts)
}

// WithMaskMessageRegex returns a new context.Context that has a modified logger
// that masks (replaces) with asterisks (`***`) all message substrings matching one
// of the given strings.
//
// Each call to this function is additive:
// the regexp to mask by are added to the existing configuration.
//
// Example:
//
//   configuration = `[regexp.MustCompile("(foo|bar)")]`
//
//   log1 = `{ msg = "banana apple ***", fields = {...}`     -> masked portion
//   log2 = `{ msg = "pineapple mango", fields = {...}`      -> as-is
//   log3 = `{ msg = "pineapple mango ***", fields = {...}`  -> masked portion
//
func WithMaskMessageRegex(ctx context.Context, expressions ...*regexp.Regexp) context.Context {
	lOpts := logging.GetSDKRootTFLoggerOpts(ctx)

	lOpts = logging.WithMaskMessageRegex(expressions...)(lOpts)

	return logging.SetSDKRootTFLoggerOpts(ctx, lOpts)
}

// WithMaskLogMatchingString returns a new context.Context that has a modified logger
// that masks (replace) with asterisks (`***`) all message substrings equal to one
// of the given strings.
//
// Each call to this function is additive:
// the string to mask by are added to the existing configuration.
//
// Example:
//
//   configuration = `['foo', 'bar']`
//
//   log1 = `{ msg = "banana apple ***", fields = {...}`     -> masked portion
//   log2 = `{ msg = "pineapple mango", fields = {...}`      -> as-is
//   log3 = `{ msg = "pineapple mango ***", fields = {...}`  -> masked portion
//
func WithMaskLogMatchingString(ctx context.Context, matchingStrings ...string) context.Context {
	lOpts := logging.GetSDKRootTFLoggerOpts(ctx)

	lOpts = logging.WithMaskMessageStrings(matchingStrings...)(lOpts)

	return logging.SetSDKRootTFLoggerOpts(ctx, lOpts)
}
