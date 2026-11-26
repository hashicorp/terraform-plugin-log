// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package logging_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-log/internal/logging"
)

func TestParseLogLevel(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		raw      string
		expected hclog.Level
	}{
		"empty": {
			raw:      "",
			expected: hclog.NoLevel,
		},
		"whitespace": {
			raw:      "   ",
			expected: hclog.NoLevel,
		},
		"invalid": {
			raw:      "INVALID",
			expected: hclog.NoLevel,
		},
		"trace": {
			raw:      "TRACE",
			expected: hclog.Trace,
		},
		"trace-lowercase-whitespace": {
			raw:      "  trace  ",
			expected: hclog.Trace,
		},
		"debug": {
			raw:      "DEBUG",
			expected: hclog.Debug,
		},
		"info": {
			raw:      "INFO",
			expected: hclog.Info,
		},
		"warn": {
			raw:      "WARN",
			expected: hclog.Warn,
		},
		"error": {
			raw:      "ERROR",
			expected: hclog.Error,
		},
		"off": {
			raw:      "OFF",
			expected: hclog.Off,
		},
		"json-uppercase": {
			raw:      "JSON",
			expected: hclog.Trace,
		},
		"json-lowercase": {
			raw:      "json",
			expected: hclog.Trace,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := logging.ParseLogLevel(testCase.raw)
			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestResolveRootLoggerEnvLevels_Matrix(t *testing.T) {
	testCases := map[string]struct {
		tfLogProvider     string
		tfLogProviderName string
		tfLogSDK          string
		expectedProvider  hclog.Level
		expectedSDK       hclog.Level
	}{
		"all-unset": {
			tfLogProvider:     "",
			tfLogProviderName: "",
			tfLogSDK:          "",
			expectedProvider:  hclog.Off,
			expectedSDK:       hclog.Off,
		},
		"tf-log-provider-debug": {
			tfLogProvider:     "DEBUG",
			tfLogProviderName: "",
			tfLogSDK:          "",
			expectedProvider:  hclog.Debug,
			expectedSDK:       hclog.Debug,
		},
		"tf-log-provider-xxx-debug-only": {
			tfLogProvider:     "",
			tfLogProviderName: "DEBUG",
			tfLogSDK:          "",
			expectedProvider:  hclog.Debug,
			expectedSDK:       hclog.Off,
		},
		"tf-log-sdk-debug-only": {
			tfLogProvider:     "",
			tfLogProviderName: "",
			tfLogSDK:          "DEBUG",
			expectedProvider:  hclog.Off,
			expectedSDK:       hclog.Debug,
		},
		"tf-log-provider-xxx-trace-and-tf-log-sdk-debug": {
			tfLogProvider:     "",
			tfLogProviderName: "TRACE",
			tfLogSDK:          "DEBUG",
			expectedProvider:  hclog.Trace,
			expectedSDK:       hclog.Debug,
		},
		"tf-log-provider-trace-and-tf-log-provider-xxx-debug": {
			tfLogProvider:     "TRACE",
			tfLogProviderName: "DEBUG",
			tfLogSDK:          "",
			expectedProvider:  hclog.Debug,
			expectedSDK:       hclog.Trace,
		},
		"tf-log-provider-trace-and-tf-log-sdk-debug": {
			tfLogProvider:     "TRACE",
			tfLogProviderName: "",
			tfLogSDK:          "DEBUG",
			expectedProvider:  hclog.Trace,
			expectedSDK:       hclog.Debug,
		},
		"tf-log-provider-trace-tf-log-provider-xxx-debug-tf-log-sdk-info": {
			tfLogProvider:     "TRACE",
			tfLogProviderName: "DEBUG",
			tfLogSDK:          "INFO",
			expectedProvider:  hclog.Debug,
			expectedSDK:       hclog.Info,
		},
		"tf-log-provider-xxx-invalid-defaults-trace": {
			tfLogProvider:     "",
			tfLogProviderName: "1",
			tfLogSDK:          "",
			expectedProvider:  hclog.Trace,
			expectedSDK:       hclog.Off,
		},
		"tf-log-provider-xxx-invalid-overrides-tf-log-provider-error": {
			tfLogProvider:     "ERROR",
			tfLogProviderName: "1",
			tfLogSDK:          "",
			expectedProvider:  hclog.Trace,
			expectedSDK:       hclog.Error,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Setenv("TF_LOG", "")
			t.Setenv("TF_LOG_SDK", testCase.tfLogSDK)
			t.Setenv("TF_LOG_PROVIDER", testCase.tfLogProvider)
			t.Setenv("TF_LOG_PROVIDER_XXX", testCase.tfLogProviderName)

			providerOpts := logging.ApplyLoggerOpts(logging.WithLevelFromEnv("TF_LOG_PROVIDER", "xxx"))
			gotProvider := providerOpts.Level
			if gotProvider == hclog.NoLevel && providerOpts.IsLevelFromEnv {
				gotProvider, _ = logging.ResolveRootProviderEnvLevel()
			}
			if diff := cmp.Diff(gotProvider, testCase.expectedProvider); diff != "" {
				t.Errorf("unexpected provider level difference: %s", diff)
			}

			sdkOpts := logging.ApplyLoggerOpts(logging.WithLevelFromEnv("TF_LOG_SDK"))
			gotSDK := sdkOpts.Level
			if gotSDK == hclog.NoLevel && sdkOpts.IsLevelFromEnv {
				gotSDK, _ = logging.ResolveRootSDKEnvLevel()
			}
			if diff := cmp.Diff(gotSDK, testCase.expectedSDK); diff != "" {
				t.Errorf("unexpected SDK level difference: %s", diff)
			}
		})
	}
}

func TestResolveRootSDKEnvLevel(t *testing.T) {
	testCases := map[string]struct {
		env        map[string]string
		expected   hclog.Level
		expectedOk bool
	}{
		"default-off-when-all-unset": {
			env:        map[string]string{},
			expected:   hclog.Off,
			expectedOk: false,
		},
		"tf-log-sdk-set": {
			env: map[string]string{
				"TF_LOG_SDK":      "INFO",
				"TF_LOG_PROVIDER": "WARN",
				"TF_LOG":          "ERROR",
			},
			expected:   hclog.Info,
			expectedOk: true,
		},
		"tf-log-sdk-off-overrides-tf-log-trace": {
			env: map[string]string{
				"TF_LOG_SDK": "OFF",
				"TF_LOG":     "TRACE",
			},
			expected:   hclog.Off,
			expectedOk: true,
		},
		"tf-log-sdk-invalid-overrides-tf-log-error": {
			env: map[string]string{
				"TF_LOG_SDK": "1",
				"TF_LOG":     "ERROR",
			},
			expected:   hclog.Trace,
			expectedOk: true,
		},
		"tf-log-provider-fallback": {
			env: map[string]string{
				"TF_LOG_PROVIDER": "WARN",
				"TF_LOG":          "ERROR",
			},
			expected:   hclog.Warn,
			expectedOk: true,
		},
		"tf-log-provider-off-overrides-tf-log-trace": {
			env: map[string]string{
				"TF_LOG_PROVIDER": "OFF",
				"TF_LOG":          "TRACE",
			},
			expected:   hclog.Off,
			expectedOk: true,
		},
		"tf-log-fallback": {
			env: map[string]string{
				"TF_LOG": "DEBUG",
			},
			expected:   hclog.Debug,
			expectedOk: true,
		},
		"tf-log-off-fallback": {
			env: map[string]string{
				"TF_LOG": "OFF",
			},
			expected:   hclog.Off,
			expectedOk: true,
		},
		"tf-log-json-fallback": {
			env: map[string]string{
				"TF_LOG": "JSON",
			},
			expected:   hclog.Trace,
			expectedOk: true,
		},
		"tf-acc-log-path-ignored": {
			env: map[string]string{
				"TF_ACC_LOG_PATH": "/tmp/acc.log",
			},
			expected:   hclog.Off,
			expectedOk: false,
		},
		"tf-log-path-mask-ignored": {
			env: map[string]string{
				"TF_LOG_PATH_MASK": "/tmp/test-%s.log",
			},
			expected:   hclog.Off,
			expectedOk: false,
		},
		"invalid-tf-log-defaults-trace": {
			env: map[string]string{
				"TF_LOG": "INVALID",
			},
			expected:   hclog.Trace,
			expectedOk: true,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Setenv("TF_LOG", "")
			t.Setenv("TF_LOG_SDK", "")
			t.Setenv("TF_LOG_PROVIDER", "")
			t.Setenv("TF_ACC_LOG_PATH", "")
			t.Setenv("TF_LOG_PATH_MASK", "")
			for k, v := range testCase.env {
				t.Setenv(k, v)
			}

			got, gotOk := logging.ResolveRootSDKEnvLevel()
			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
			if gotOk != testCase.expectedOk {
				t.Errorf("expected ResolveRootSDKEnvLevel ok=%v, got %v", testCase.expectedOk, gotOk)
			}
		})
	}
}

func TestResolveRootProviderEnvLevel(t *testing.T) {
	testCases := map[string]struct {
		env        map[string]string
		expected   hclog.Level
		expectedOk bool
	}{
		"default-off-when-all-unset": {
			env:        map[string]string{},
			expected:   hclog.Off,
			expectedOk: false,
		},
		"tf-log-provider-fallback": {
			env: map[string]string{
				"TF_LOG_PROVIDER": "WARN",
				"TF_LOG":          "ERROR",
			},
			expected:   hclog.Warn,
			expectedOk: true,
		},
		"tf-log-provider-off-overrides-tf-log-trace": {
			env: map[string]string{
				"TF_LOG_PROVIDER": "OFF",
				"TF_LOG":          "TRACE",
			},
			expected:   hclog.Off,
			expectedOk: true,
		},
		"tf-log-provider-invalid-overrides-tf-log-error": {
			env: map[string]string{
				"TF_LOG_PROVIDER": "1",
				"TF_LOG":          "ERROR",
			},
			expected:   hclog.Trace,
			expectedOk: true,
		},
		"tf-log-fallback": {
			env: map[string]string{
				"TF_LOG": "DEBUG",
			},
			expected:   hclog.Debug,
			expectedOk: true,
		},
		"tf-log-off-fallback": {
			env: map[string]string{
				"TF_LOG": "OFF",
			},
			expected:   hclog.Off,
			expectedOk: true,
		},
		"tf-log-json-fallback": {
			env: map[string]string{
				"TF_LOG": "JSON",
			},
			expected:   hclog.Trace,
			expectedOk: true,
		},
		"tf-acc-log-path-ignored": {
			env: map[string]string{
				"TF_ACC_LOG_PATH": "/tmp/acc.log",
			},
			expected:   hclog.Off,
			expectedOk: false,
		},
		"tf-log-path-mask-ignored": {
			env: map[string]string{
				"TF_LOG_PATH_MASK": "/tmp/test-%s.log",
			},
			expected:   hclog.Off,
			expectedOk: false,
		},
		"invalid-tf-log-defaults-trace": {
			env: map[string]string{
				"TF_LOG": "INVALID",
			},
			expected:   hclog.Trace,
			expectedOk: true,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Setenv("TF_LOG", "")
			t.Setenv("TF_LOG_SDK", "")
			t.Setenv("TF_LOG_PROVIDER", "")
			t.Setenv("TF_ACC_LOG_PATH", "")
			t.Setenv("TF_LOG_PATH_MASK", "")
			for k, v := range testCase.env {
				t.Setenv(k, v)
			}

			got, gotOk := logging.ResolveRootProviderEnvLevel()
			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
			if gotOk != testCase.expectedOk {
				t.Errorf("expected ResolveRootProviderEnvLevel ok=%v, got %v", testCase.expectedOk, gotOk)
			}
		})
	}
}
