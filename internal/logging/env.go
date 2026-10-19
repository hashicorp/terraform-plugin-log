// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package logging

import (
	"os"
	"strings"

	"github.com/hashicorp/go-hclog"
)

const (
	envLog         = "TF_LOG"
	envLogSDK      = "TF_LOG_SDK"
	envLogProvider = "TF_LOG_PROVIDER"
	envAccLogFile  = "TF_ACC_LOG_PATH"
)

// ParseLogLevel parses a Terraform log level environment variable value into
// an hclog.Level.
//
// Valid values are TRACE, DEBUG, INFO, WARN, ERROR, and OFF. Empty or
// unrecognized values return hclog.NoLevel.
func ParseLogLevel(raw string) hclog.Level {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return hclog.NoLevel
	}

	return hclog.LevelFromString(trimmed)
}

// parseEnvLogLevel reads envVar from the environment and resolves its log
// level. If the variable is unset or whitespace-only, it returns
// (hclog.NoLevel, false). If it is set to a recognized level, that level is
// returned; otherwise, non-empty unrecognized values default to hclog.Trace to
// match Terraform Core behavior (e.g. TF_LOG=1).
func parseEnvLogLevel(envVar string) (hclog.Level, bool) {
	raw := strings.TrimSpace(os.Getenv(envVar))
	if raw == "" {
		return hclog.NoLevel, false
	}

	if lvl := ParseLogLevel(raw); lvl != hclog.NoLevel {
		return lvl, true
	}

	return hclog.Trace, true
}

// ResolveRootSDKEnvLevel resolves the effective log level for the root SDK
// logger when configured via WithLevelFromEnv, allowing SDK logging to behave
// consistently across production and acceptance testing.
//
// Environment variable precedence:
//  1. TF_LOG_SDK
//  2. TF_LOG_PROVIDER
//  3. TF_LOG
//  4. TF_ACC_LOG_PATH (if non-empty, defaults to hclog.Trace)
//  5. Default: hclog.Off
func ResolveRootSDKEnvLevel() hclog.Level {
	return resolveRootEnvLevel(envLogSDK, envLogProvider, envLog)
}

// ResolveRootProviderEnvLevel resolves the effective log level for the root
// provider logger when configured via WithLevelFromEnv, allowing provider
// logging to behave consistently across production and acceptance testing.
//
// Environment variable precedence:
//  1. TF_LOG_PROVIDER
//  2. TF_LOG
//  3. TF_ACC_LOG_PATH (if non-empty, defaults to hclog.Trace)
//  4. Default: hclog.Off
func ResolveRootProviderEnvLevel() hclog.Level {
	return resolveRootEnvLevel(envLogProvider, envLog)
}

func resolveRootEnvLevel(envVars ...string) hclog.Level {
	for _, envVar := range envVars {
		if lvl, ok := parseEnvLogLevel(envVar); ok {
			return lvl
		}
	}

	// helper/resource makes this default to TRACE, so we should, too.
	if strings.TrimSpace(os.Getenv(envAccLogFile)) != "" {
		return hclog.Trace
	}

	return hclog.Off
}
