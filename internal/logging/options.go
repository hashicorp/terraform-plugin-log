package logging

import (
	"io"
	"os"
	"regexp"

	"github.com/hashicorp/go-hclog"
)

// Option defines a modification to the configuration for a logger.
type Option func(LoggerOpts) LoggerOpts

// LoggerOpts is a collection of configuration settings for loggers.
type LoggerOpts struct {
	// Name is the name or "@module" of a logger.
	Name string

	// Level is the most verbose level that a logger will write logs for.
	Level hclog.Level

	// IncludeLocation indicates whether logs should include the location
	// of the logging statement or not.
	IncludeLocation bool

	// AdditionalLocationOffset is the number of additional stack levels to
	// skip when finding the file and line information for the log line.
	// Defaults to 1 to account for the tflog and tfsdklog logging functions.
	AdditionalLocationOffset int

	// Output dictates where logs are written to. Output should only ever
	// be set by tflog or tfsdklog, never by SDK authors or provider
	// developers. Where logs get written to is complex and delicate and
	// requires a deep understanding of Terraform's architecture, and it's
	// easy to mess up on accident.
	Output io.Writer

	// IncludeTime indicates whether logs should incude the time they were
	// written or not. It should only be set to true when testing tflog or
	// tfsdklog; providers and SDKs should always include the time logs
	// were written as part of the log.
	IncludeTime bool

	// IncludeRootFields indicates whether a new subsystem logger should
	// copy existing fields from the root logger. This is only performed
	// at the time of new subsystem creation.
	IncludeRootFields bool

	// OmitLogWithFieldKeys indicates that the logger should omit to write
	// any log when any of the given keys is found within the fields.
	//
	// Example:
	//
	//   OmitLogWithFieldKeys = `['foo', 'baz']`
	//
	//   log1 = `{ msg = "...", fields = { 'foo', '...', 'bar', '...' }`   -> omitted
	//   log2 = `{ msg = "...", fields = { 'bar', '...' }`                 -> printed
	//   log3 = `{ msg = "...", fields = { 'baz`', '...', 'boo', '...' }`  -> omitted
	//
	OmitLogWithFieldKeys []string

	// OmitLogWithMessageRegex indicates that the logger should omit to write
	// any log that matches any of the given *regexp.Regexp.
	//
	// Example:
	//
	//   OmitLogWithMessageRegex = `[regexp.MustCompile("(foo|bar)")]`
	//
	//   log1 = `{ msg = "banana apple foo", fields = {...}`     -> omitted
	//   log2 = `{ msg = "pineapple mango", fields = {...}`      -> printed
	//   log3 = `{ msg = "pineapple mango bar", fields = {...}`  -> omitted
	//
	OmitLogWithMessageRegex []*regexp.Regexp

	// OmitLogWithMessageStrings indicates that the logger should omit to write
	// any log that matches any of the given string.
	//
	// Example:
	//
	//   OmitLogWithMessageStrings = `['foo', 'bar']`
	//
	//   log1 = `{ msg = "banana apple foo", fields = {...}`     -> omitted
	//   log2 = `{ msg = "pineapple mango", fields = {...}`      -> printed
	//   log3 = `{ msg = "pineapple mango bar", fields = {...}`  -> omitted
	//
	OmitLogWithMessageStrings []string

	// MaskFieldValueWithFieldKeys indicates that the logger should mask with asterisks (`*`)
	// any field value where the key matches one of the given keys.
	//
	// Example:
	//
	//   MaskFieldValueWithFieldKeys = `['foo', 'baz']`
	//
	//   log1 = `{ msg = "...", fields = { 'foo', '***', 'bar', '...' }`   -> masked value
	//   log2 = `{ msg = "...", fields = { 'bar', '...' }`                 -> as-is value
	//   log3 = `{ msg = "...", fields = { 'baz`', '***', 'boo', '...' }`  -> masked value
	//
	MaskFieldValueWithFieldKeys []string

	// MaskMessageRegex indicates that the logger should replace, within
	// a log message, the portion matching one of the given *regexp.Regexp.
	//
	// Example:
	//
	//   MaskMessageRegex = `[regexp.MustCompile("(foo|bar)")]`
	//
	//   log1 = `{ msg = "banana apple ***", fields = {...}`     -> masked portion
	//   log2 = `{ msg = "pineapple mango", fields = {...}`      -> as-is
	//   log3 = `{ msg = "pineapple mango ***", fields = {...}`  -> masked portion
	//
	MaskMessageRegex []*regexp.Regexp

	// MaskMessageStrings indicates that the logger should replace, within
	// a log message, the portion matching one of the given strings.
	//
	// Example:
	//
	//   MaskMessageStrings = `['foo', 'bar']`
	//
	//   log1 = `{ msg = "banana apple ***", fields = {...}`     -> masked portion
	//   log2 = `{ msg = "pineapple mango", fields = {...}`      -> as-is
	//   log3 = `{ msg = "pineapple mango ***", fields = {...}`  -> masked portion
	//
	MaskMessageStrings []string
}

// ApplyLoggerOpts generates a LoggerOpts out of a list of Option
// implementations. By default, AdditionalLocationOffset is 1, IncludeLocation
// is true, IncludeTime is true, and Output is os.Stderr.
func ApplyLoggerOpts(opts ...Option) LoggerOpts {
	// set some defaults
	l := LoggerOpts{
		AdditionalLocationOffset: 1,
		IncludeLocation:          true,
		IncludeTime:              true,
		Output:                   os.Stderr,
	}
	for _, opt := range opts {
		l = opt(l)
	}
	return l
}

// WithAdditionalLocationOffset sets the WithAdditionalLocationOffset
// configuration option, allowing implementations to fix location information
// when implementing helper functions. The default offset of 1 is automatically
// added to the provided value to account for the tflog and tfsdk logging
// functions.
func WithAdditionalLocationOffset(additionalLocationOffset int) Option {
	return func(l LoggerOpts) LoggerOpts {
		l.AdditionalLocationOffset = additionalLocationOffset + 1
		return l
	}
}

// WithOutput sets the Output configuration option, controlling where logs get
// written to. This is mostly used for testing (to write to os.Stdout, so the
// test framework can compare it against the example output) and as a helper
// when implementing safe, specific output strategies in tfsdklog.
func WithOutput(output io.Writer) Option {
	return func(l LoggerOpts) LoggerOpts {
		l.Output = output
		return l
	}
}

// WithRootFields enables the copying of root logger fields to a new subsystem
// logger during creation.
func WithRootFields() Option {
	return func(l LoggerOpts) LoggerOpts {
		l.IncludeRootFields = true
		return l
	}
}

// WithoutLocation disables the location included with logging statements. It
// should only ever be used to make log output deterministic when testing
// terraform-plugin-log.
func WithoutLocation() Option {
	return func(l LoggerOpts) LoggerOpts {
		l.IncludeLocation = false
		return l
	}
}

// WithoutTimestamp disables the timestamp included with logging statements. It
// should only ever be used to make log output deterministic when testing
// terraform-plugin-log.
func WithoutTimestamp() Option {
	return func(l LoggerOpts) LoggerOpts {
		l.IncludeTime = false
		return l
	}
}

// WithOmitLogWithFieldKeys appends keys to the LoggerOpts.OmitLogWithFieldKeys field.
func WithOmitLogWithFieldKeys(keys ...string) Option {
	return func(l LoggerOpts) LoggerOpts {
		l.OmitLogWithFieldKeys = append(l.OmitLogWithFieldKeys, keys...)
		return l
	}
}

// WithOmitLogWithMessageRegex appends *regexp.Regexp to the LoggerOpts.OmitLogWithMessageRegex field.
func WithOmitLogWithMessageRegex(expressions ...*regexp.Regexp) Option {
	return func(l LoggerOpts) LoggerOpts {
		l.OmitLogWithMessageRegex = append(l.OmitLogWithMessageRegex, expressions...)
		return l
	}
}

// WithOmitLogWithMessageStrings appends string to the LoggerOpts.OmitLogWithMessageStrings field.
func WithOmitLogWithMessageStrings(matchingStrings ...string) Option {
	return func(l LoggerOpts) LoggerOpts {
		l.OmitLogWithMessageStrings = append(l.OmitLogWithMessageStrings, matchingStrings...)
		return l
	}
}

// WithMaskFieldValueWithFieldKeys appends keys to the LoggerOpts.MaskFieldValueWithFieldKeys field.
func WithMaskFieldValueWithFieldKeys(keys ...string) Option {
	return func(l LoggerOpts) LoggerOpts {
		l.MaskFieldValueWithFieldKeys = append(l.MaskFieldValueWithFieldKeys, keys...)
		return l
	}
}

// WithMaskMessageRegex appends *regexp.Regexp to the LoggerOpts.MaskMessageRegex field.
func WithMaskMessageRegex(expressions ...*regexp.Regexp) Option {
	return func(l LoggerOpts) LoggerOpts {
		l.MaskMessageRegex = append(l.MaskMessageRegex, expressions...)
		return l
	}
}

// WithMaskMessageStrings appends string to the LoggerOpts.MaskMessageStrings field.
func WithMaskMessageStrings(matchingStrings ...string) Option {
	return func(l LoggerOpts) LoggerOpts {
		l.MaskMessageStrings = append(l.MaskMessageStrings, matchingStrings...)
		return l
	}
}
