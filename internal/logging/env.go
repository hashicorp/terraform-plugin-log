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
)

// ParseLogLevel parses a Terraform log level environment variable value into
// an hclog.Level.
//
// Valid values are TRACE, DEBUG, INFO, WARN, ERROR, OFF, and the Terraform
// JSON pseudo-level (which maps to hclog.Trace). Empty or unrecognized values
// return hclog.NoLevel.
func ParseLogLevel(raw string) hclog.Level {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return hclog.NoLevel
	}

	if strings.EqualFold(trimmed, "JSON") {
		return hclog.Trace
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
// consistently across production and acceptance testing. Sinks maintain an
// internal level when no env vars are set, and need to know whether the level
// was explicitly configured or not.
//
// Environment variable precedence:
//  1. TF_LOG_SDK
//  2. TF_LOG_PROVIDER
//  3. TF_LOG
//  4. Default: hclog.Off
func ResolveRootSDKEnvLevel() (hclog.Level, bool) {
	return resolveRootEnvLevel(envLogSDK, envLogProvider, envLog)
}

// ResolveRootProviderEnvLevel resolves the effective log level for the root
// provider logger when configured via WithLevelFromEnv, allowing provider
// logging to behave consistently across production and acceptance testing.
// Sinks maintain an internal level when no env vars are set, and need to know
// whether the level was explicitly configured or not.
//
// Environment variable precedence:
//  1. TF_LOG_PROVIDER
//  2. TF_LOG
//  3. Default: hclog.Off
func ResolveRootProviderEnvLevel() (hclog.Level, bool) {
	return resolveRootEnvLevel(envLogProvider, envLog)
}

func resolveRootEnvLevel(envVars ...string) (hclog.Level, bool) {
	for _, envVar := range envVars {
		if lvl, ok := parseEnvLogLevel(envVar); ok {
			return lvl, true
		}
	}

	return hclog.Off, false
}
