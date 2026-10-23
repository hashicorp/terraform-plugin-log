// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package tflog

import (
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-log/internal/logging"
)

// Options are a collection of logging options, useful for collecting arguments
// to NewSubsystem prior to calling it.
type Options []logging.Option

// WithAdditionalLocationOffset returns an option that allowing implementations
// to fix location information when implementing helper functions. The default
// offset of 1 is automatically added to the provided value to account for the
// tflog logging functions.
func WithAdditionalLocationOffset(additionalLocationOffset int) logging.Option {
	return logging.WithAdditionalLocationOffset(additionalLocationOffset)
}

// WithLevelFromEnv returns an option that will set the level of the logger
// based on the string in an environment variable. The environment variable
// checked will be `name` and `subsystems`, joined by _ and in all caps.
//
// Valid environment variable values are TRACE, DEBUG, INFO, WARN, ERROR, OFF,
// and JSON (which maps to hclog.Trace). For root provider loggers, if the
// target environment variable is unset, the level falls back to
// TF_LOG_PROVIDER, TF_LOG, and TF_ACC_LOG_PATH before defaulting to hclog.Off.
// For subsystem loggers, an unset environment variable inherits the root
// logger's level.
func WithLevelFromEnv(name string, subsystems ...string) logging.Option {
	return logging.WithLevelFromEnv(name, subsystems...)
}

// WithLevel returns an option that will set the level of the logger.
func WithLevel(level hclog.Level) logging.Option {
	return logging.WithLevel(level)
}

// WithRootFields enables the copying of root logger fields to a new subsystem
// logger during creation.
func WithRootFields() logging.Option {
	return logging.WithRootFields()
}

// WithoutLocation returns an option that disables including the location of
// the log line in the log output, which is on by default. This has no effect
// when used with NewSubsystem.
func WithoutLocation() logging.Option {
	return func(l logging.LoggerOpts) logging.LoggerOpts {
		l.IncludeLocation = false
		return l
	}
}
