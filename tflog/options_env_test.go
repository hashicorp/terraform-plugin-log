// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package tflog_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-log/internal/loggertest"
	"github.com/hashicorp/terraform-plugin-log/internal/logging"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-log/tfsdklog"
)

func TestWithLevelFromEnv(t *testing.T) {
	testCases := map[string]struct {
		env                      map[string]string
		rootOptions              []logging.Option
		subsystemOptions         []logging.Option
		expectedRootOutput       []map[string]interface{}
		expectedSubsystemOutput  []map[string]interface{}
		expectedRootIsTrace      bool
		expectedSubsystemIsTrace bool
	}{
		"default-off-when-all-unset": {
			env:                      map[string]string{},
			rootOptions:              []logging.Option{tflog.WithLevelFromEnv("TF_LOG_PROVIDER", "example")},
			subsystemOptions:         []logging.Option{tflog.WithLevelFromEnv("TF_LOG_PROVIDER", "example", "my_subsystem")},
			expectedRootOutput:       nil,
			expectedSubsystemOutput:  nil,
			expectedRootIsTrace:      false,
			expectedSubsystemIsTrace: false,
		},
		"tf-log-provider-fallback-debug": {
			env: map[string]string{
				"TF_LOG_PROVIDER": "DEBUG",
			},
			rootOptions:      []logging.Option{tflog.WithLevelFromEnv("TF_LOG_PROVIDER", "example")},
			subsystemOptions: []logging.Option{tflog.WithLevelFromEnv("TF_LOG_PROVIDER", "example", "my_subsystem")},
			expectedRootOutput: []map[string]interface{}{
				{
					"@level":   "debug",
					"@message": "root debug",
					"@module":  "provider",
				},
			},
			expectedSubsystemOutput: []map[string]interface{}{
				{
					"@level":   "debug",
					"@message": "subsystem debug",
					"@module":  "provider.my_subsystem",
				},
			},
			expectedRootIsTrace:      false,
			expectedSubsystemIsTrace: false,
		},
		"tf-log-fallback-debug": {
			env: map[string]string{
				"TF_LOG": "DEBUG",
			},
			rootOptions:      []logging.Option{tflog.WithLevelFromEnv("TF_LOG_PROVIDER", "example")},
			subsystemOptions: []logging.Option{tflog.WithLevelFromEnv("TF_LOG_PROVIDER", "example", "my_subsystem")},
			expectedRootOutput: []map[string]interface{}{
				{
					"@level":   "debug",
					"@message": "root debug",
					"@module":  "provider",
				},
			},
			expectedSubsystemOutput: []map[string]interface{}{
				{
					"@level":   "debug",
					"@message": "subsystem debug",
					"@module":  "provider.my_subsystem",
				},
			},
			expectedRootIsTrace:      false,
			expectedSubsystemIsTrace: false,
		},
		"subsystem-override-trace-when-root-off": {
			env: map[string]string{
				"TF_LOG_PROVIDER_EXAMPLE_MY_SUBSYSTEM": "TRACE",
			},
			rootOptions:        []logging.Option{tflog.WithLevelFromEnv("TF_LOG_PROVIDER", "example")},
			subsystemOptions:   []logging.Option{tflog.WithLevelFromEnv("TF_LOG_PROVIDER", "example", "my_subsystem")},
			expectedRootOutput: nil,
			expectedSubsystemOutput: []map[string]interface{}{
				{
					"@level":   "trace",
					"@message": "subsystem trace",
					"@module":  "provider.my_subsystem",
				},
				{
					"@level":   "debug",
					"@message": "subsystem debug",
					"@module":  "provider.my_subsystem",
				},
			},
			expectedRootIsTrace:      false,
			expectedSubsystemIsTrace: true,
		},
		"with-level-overrides-with-level-from-env": {
			env:         map[string]string{},
			rootOptions: []logging.Option{tflog.WithLevelFromEnv("TF_LOG_PROVIDER", "example"), tflog.WithLevel(hclog.Debug)},
			subsystemOptions: []logging.Option{
				tflog.WithLevelFromEnv("TF_LOG_PROVIDER", "example", "my_subsystem"),
			},
			expectedRootOutput: []map[string]interface{}{
				{
					"@level":   "debug",
					"@message": "root debug",
					"@module":  "provider",
				},
			},
			expectedSubsystemOutput: []map[string]interface{}{
				{
					"@level":   "debug",
					"@message": "subsystem debug",
					"@module":  "provider.my_subsystem",
				},
			},
			expectedRootIsTrace:      false,
			expectedSubsystemIsTrace: false,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Setenv("TF_LOG", "")
			t.Setenv("TF_LOG_PROVIDER", "")
			t.Setenv("TF_LOG_PROVIDER_EXAMPLE", "")
			t.Setenv("TF_LOG_PROVIDER_EXAMPLE_MY_SUBSYSTEM", "")
			t.Setenv("TF_ACC_LOG_PATH", "")
			for k, v := range testCase.env {
				t.Setenv(k, v)
			}

			var rootBuf bytes.Buffer
			ctx := context.Background()
			rootOpts := append(
				[]logging.Option{
					logging.WithoutLocation(),
					logging.WithoutTimestamp(),
					logging.WithOutput(&rootBuf),
				},
				testCase.rootOptions...,
			)
			ctx = tfsdklog.NewRootProviderLogger(ctx, rootOpts...)
			ctx = tflog.NewSubsystem(ctx, "my_subsystem", testCase.subsystemOptions...)

			if got := tflog.IsTrace(ctx); got != testCase.expectedRootIsTrace {
				t.Errorf("expected IsTrace=%v, got %v", testCase.expectedRootIsTrace, got)
			}
			if got := tflog.SubsystemIsTrace(ctx, "my_subsystem"); got != testCase.expectedSubsystemIsTrace {
				t.Errorf("expected SubsystemIsTrace=%v, got %v", testCase.expectedSubsystemIsTrace, got)
			}

			tflog.Trace(ctx, "root trace")
			tflog.Debug(ctx, "root debug")

			gotRoot, err := loggertest.MultilineJSONDecode(&rootBuf)
			if err != nil {
				t.Fatalf("unable to read root JSON: %s", err)
			}
			if diff := cmp.Diff(testCase.expectedRootOutput, gotRoot); diff != "" {
				t.Errorf("unexpected root output difference: %s", diff)
			}

			rootBuf.Reset()
			tflog.SubsystemTrace(ctx, "my_subsystem", "subsystem trace")
			tflog.SubsystemDebug(ctx, "my_subsystem", "subsystem debug")

			gotSub, err := loggertest.MultilineJSONDecode(&rootBuf)
			if err != nil {
				t.Fatalf("unable to read subsystem JSON: %s", err)
			}
			if diff := cmp.Diff(testCase.expectedSubsystemOutput, gotSub); diff != "" {
				t.Errorf("unexpected subsystem output difference: %s", diff)
			}
		})
	}
}
