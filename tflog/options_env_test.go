// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package tflog_test

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
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
		"test-log-path-env-vars-ignored": {
			env: map[string]string{
				"TF_ACC_LOG_PATH":  "/tmp/acc.log",
				"TF_LOG_PATH_MASK": "/tmp/test-%s.log",
			},
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
		"tf-log-provider-example-json": {
			env: map[string]string{
				"TF_LOG_PROVIDER_EXAMPLE": "JSON",
			},
			rootOptions:      []logging.Option{tflog.WithLevelFromEnv("TF_LOG_PROVIDER", "example")},
			subsystemOptions: []logging.Option{tflog.WithLevelFromEnv("TF_LOG_PROVIDER", "example", "my_subsystem")},
			expectedRootOutput: []map[string]interface{}{
				{
					"@level":   "trace",
					"@message": "root trace",
					"@module":  "provider",
				},
				{
					"@level":   "debug",
					"@message": "root debug",
					"@module":  "provider",
				},
			},
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
			expectedRootIsTrace:      true,
			expectedSubsystemIsTrace: true,
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
			t.Setenv("TF_LOG_PATH_MASK", "")
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

func TestWithLevelFromEnv_Sinks(t *testing.T) {
	testLoggingCases := map[string]struct {
		env                map[string]string
		useAccLogPath      bool
		useLogPathMask     bool
		expectedIsTrace    bool
		expectedIsDebug    bool
		expectedIsInfo     bool
		expectedEmitted    []string
		expectedSuppressed []string
	}{
		"all-unset": {
			env:                map[string]string{},
			expectedIsTrace:    false,
			expectedIsDebug:    false,
			expectedIsInfo:     false,
			expectedSuppressed: []string{"root trace", "root debug", "root info", "sub trace", "sub debug"},
		},
		"tf-acc-log-path": {
			env:             map[string]string{},
			useAccLogPath:   true,
			expectedIsTrace: true,
			expectedIsDebug: true,
			expectedIsInfo:  true,
			expectedEmitted: []string{"root trace", "root debug", "root info", "sub trace", "sub debug"},
		},
		"tf-log-path-mask-all-unset": {
			env:                map[string]string{},
			useLogPathMask:     true,
			expectedIsTrace:    false,
			expectedIsDebug:    false,
			expectedIsInfo:     false,
			expectedSuppressed: []string{"root trace", "root debug", "root info", "sub trace", "sub debug"},
		},
		"tf-log-path-mask-with-tf-log-provider-debug": {
			env: map[string]string{
				"TF_LOG_PROVIDER": "DEBUG",
			},
			useLogPathMask:     true,
			expectedIsTrace:    false,
			expectedIsDebug:    true,
			expectedIsInfo:     true,
			expectedEmitted:    []string{"root debug", "root info", "sub debug"},
			expectedSuppressed: []string{"root trace", "sub trace"},
		},
		"tf-log-debug": {
			env: map[string]string{
				"TF_LOG": "DEBUG",
			},
			expectedIsTrace:    false,
			expectedIsDebug:    true,
			expectedIsInfo:     true,
			expectedEmitted:    []string{"root debug", "root info", "sub debug"},
			expectedSuppressed: []string{"root trace", "sub trace"},
		},
		"tf-log-provider-debug": {
			env: map[string]string{
				"TF_LOG_PROVIDER": "DEBUG",
			},
			expectedIsTrace:    false,
			expectedIsDebug:    true,
			expectedIsInfo:     true,
			expectedEmitted:    []string{"root debug", "root info", "sub debug"},
			expectedSuppressed: []string{"root trace", "sub trace"},
		},
		"tf-log-truthy": {
			env: map[string]string{
				"TF_LOG": "1",
			},
			expectedIsTrace: true,
			expectedIsDebug: true,
			expectedIsInfo:  true,
			expectedEmitted: []string{"root trace", "root debug", "root info", "sub trace", "sub debug"},
		},
		"tf-log-truthy-with-tf-log-provider-debug": {
			env: map[string]string{
				"TF_LOG":          "1",
				"TF_LOG_PROVIDER": "DEBUG",
			},
			expectedIsTrace:    false,
			expectedIsDebug:    true,
			expectedIsInfo:     true,
			expectedEmitted:    []string{"root debug", "root info", "sub debug"},
			expectedSuppressed: []string{"root trace", "sub trace"},
		},
	}

	for name, testCase := range testLoggingCases {
		t.Run("ContextWithTestLogging-"+name, func(t *testing.T) {
			tempDir := t.TempDir()
			logPath := filepath.Join(tempDir, "test.log")
			t.Setenv("TF_LOG", "")
			t.Setenv("TF_LOG_PROVIDER", "")
			t.Setenv("TF_LOG_PROVIDER_EXAMPLE", "")
			t.Setenv("TF_LOG_PROVIDER_EXAMPLE_SINK_SUB", "")
			t.Setenv("TF_ACC_LOG_PATH", "")
			t.Setenv("TF_LOG_PATH_MASK", "")
			t.Setenv("TF_LOG_PATH", "")
			switch {
			case testCase.useAccLogPath:
				t.Setenv("TF_ACC_LOG_PATH", logPath)
			case testCase.useLogPathMask:
				mask := filepath.Join(tempDir, "mask-%s.log")
				logPath = fmt.Sprintf(mask, strings.ReplaceAll(t.Name(), "/", "__"))
				t.Setenv("TF_LOG_PATH_MASK", mask)
			default:
				t.Setenv("TF_LOG_PATH", logPath)
			}
			for k, v := range testCase.env {
				t.Setenv(k, v)
			}

			ctx := tfsdklog.ContextWithTestLogging(context.Background(), t.Name())
			ctx = tfsdklog.NewRootProviderLogger(ctx, tflog.WithLevelFromEnv("TF_LOG_PROVIDER", "example"))
			ctx = tflog.NewSubsystem(ctx, "sink_sub", tflog.WithLevelFromEnv("TF_LOG_PROVIDER", "example", "sink_sub"))

			if got := tflog.IsTrace(ctx); got != testCase.expectedIsTrace {
				t.Errorf("expected IsTrace=%v, got %v", testCase.expectedIsTrace, got)
			}
			if got := tflog.IsDebug(ctx); got != testCase.expectedIsDebug {
				t.Errorf("expected IsDebug=%v, got %v", testCase.expectedIsDebug, got)
			}
			if got := tflog.IsInfo(ctx); got != testCase.expectedIsInfo {
				t.Errorf("expected IsInfo=%v, got %v", testCase.expectedIsInfo, got)
			}
			if got := tflog.SubsystemIsTrace(ctx, "sink_sub"); got != testCase.expectedIsTrace {
				t.Errorf("expected SubsystemIsTrace=%v, got %v", testCase.expectedIsTrace, got)
			}
			if got := tflog.SubsystemIsDebug(ctx, "sink_sub"); got != testCase.expectedIsDebug {
				t.Errorf("expected SubsystemIsDebug=%v, got %v", testCase.expectedIsDebug, got)
			}
			if got := tflog.SubsystemIsInfo(ctx, "sink_sub"); got != testCase.expectedIsInfo {
				t.Errorf("expected SubsystemIsInfo=%v, got %v", testCase.expectedIsInfo, got)
			}

			tflog.Trace(ctx, "root trace")
			tflog.Debug(ctx, "root debug")
			tflog.Info(ctx, "root info")
			tflog.SubsystemTrace(ctx, "sink_sub", "sub trace")
			tflog.SubsystemDebug(ctx, "sink_sub", "sub debug")

			rawOut, err := os.ReadFile(logPath)
			if err != nil {
				t.Fatalf("unable to read TF_LOG_PATH file: %s", err)
			}
			out := string(rawOut)
			for _, msg := range testCase.expectedEmitted {
				if !strings.Contains(out, msg) {
					t.Errorf("expected log output to contain %q, got: %s", msg, out)
				}
			}
			for _, msg := range testCase.expectedSuppressed {
				if strings.Contains(out, msg) {
					t.Errorf("expected log output not to contain %q, got: %s", msg, out)
				}
			}
		})
	}

	standardLoggingCases := map[string]struct {
		env                map[string]string
		expectedIsTrace    bool
		expectedIsDebug    bool
		expectedIsInfo     bool
		expectedEmitted    []string
		expectedSuppressed []string
	}{
		"all-unset-preserves-default-info": {
			env:                map[string]string{},
			expectedIsTrace:    false,
			expectedIsDebug:    false,
			expectedIsInfo:     true,
			expectedEmitted:    []string{"root info"},
			expectedSuppressed: []string{"root trace", "root debug"},
		},
		"test-log-path-env-vars-ignored": {
			env: map[string]string{
				"TF_ACC_LOG_PATH":  filepath.Join(t.TempDir(), "unused-acc.log"),
				"TF_LOG_PATH_MASK": filepath.Join(t.TempDir(), "unused-%s.log"),
			},
			expectedIsTrace:    false,
			expectedIsDebug:    false,
			expectedIsInfo:     true,
			expectedEmitted:    []string{"root info"},
			expectedSuppressed: []string{"root trace", "root debug"},
		},
		"tf-log-off": {
			env: map[string]string{
				"TF_LOG": "OFF",
			},
			expectedIsTrace:    false,
			expectedIsDebug:    false,
			expectedIsInfo:     false,
			expectedSuppressed: []string{"root trace", "root debug", "root info"},
		},
		"tf-log-debug": {
			env: map[string]string{
				"TF_LOG": "DEBUG",
			},
			expectedIsTrace:    false,
			expectedIsDebug:    true,
			expectedIsInfo:     true,
			expectedEmitted:    []string{"root debug", "root info"},
			expectedSuppressed: []string{"root trace"},
		},
		"tf-log-provider-debug": {
			env: map[string]string{
				"TF_LOG_PROVIDER": "DEBUG",
			},
			expectedIsTrace:    false,
			expectedIsDebug:    true,
			expectedIsInfo:     true,
			expectedEmitted:    []string{"root debug", "root info"},
			expectedSuppressed: []string{"root trace"},
		},
		"tf-log-truthy": {
			env: map[string]string{
				"TF_LOG": "1",
			},
			expectedIsTrace: true,
			expectedIsDebug: true,
			expectedIsInfo:  true,
			expectedEmitted: []string{"root trace", "root debug", "root info"},
		},
		"tf-log-truthy-with-tf-log-provider-debug": {
			env: map[string]string{
				"TF_LOG":          "1",
				"TF_LOG_PROVIDER": "DEBUG",
			},
			expectedIsTrace:    false,
			expectedIsDebug:    true,
			expectedIsInfo:     true,
			expectedEmitted:    []string{"root debug", "root info"},
			expectedSuppressed: []string{"root trace"},
		},
	}

	for name, testCase := range standardLoggingCases {
		t.Run("ContextWithStandardLogging-"+name, func(t *testing.T) {
			t.Setenv("TF_LOG", "")
			t.Setenv("TF_LOG_PROVIDER", "")
			t.Setenv("TF_LOG_PROVIDER_EXAMPLE", "")
			t.Setenv("TF_LOG_PROVIDER_EXAMPLE_SINK_SUB", "")
			t.Setenv("TF_ACC_LOG_PATH", "")
			t.Setenv("TF_LOG_PATH_MASK", "")
			for k, v := range testCase.env {
				t.Setenv(k, v)
			}

			var stdBuf bytes.Buffer
			origOut := log.Default().Writer()
			log.Default().SetOutput(&stdBuf)
			t.Cleanup(func() {
				log.Default().SetOutput(origOut)
			})

			ctx := tfsdklog.ContextWithStandardLogging(context.Background(), t.Name())
			ctx = tfsdklog.NewRootProviderLogger(ctx, tflog.WithLevelFromEnv("TF_LOG_PROVIDER", "example"))
			ctx = tflog.NewSubsystem(ctx, "sink_sub", tflog.WithLevelFromEnv("TF_LOG_PROVIDER", "example", "sink_sub"))

			if got := tflog.IsTrace(ctx); got != testCase.expectedIsTrace {
				t.Errorf("expected IsTrace=%v, got %v", testCase.expectedIsTrace, got)
			}
			if got := tflog.IsDebug(ctx); got != testCase.expectedIsDebug {
				t.Errorf("expected IsDebug=%v, got %v", testCase.expectedIsDebug, got)
			}
			if got := tflog.IsInfo(ctx); got != testCase.expectedIsInfo {
				t.Errorf("expected IsInfo=%v, got %v", testCase.expectedIsInfo, got)
			}
			if got := tflog.SubsystemIsTrace(ctx, "sink_sub"); got != testCase.expectedIsTrace {
				t.Errorf("expected SubsystemIsTrace=%v, got %v", testCase.expectedIsTrace, got)
			}
			if got := tflog.SubsystemIsDebug(ctx, "sink_sub"); got != testCase.expectedIsDebug {
				t.Errorf("expected SubsystemIsDebug=%v, got %v", testCase.expectedIsDebug, got)
			}
			if got := tflog.SubsystemIsInfo(ctx, "sink_sub"); got != testCase.expectedIsInfo {
				t.Errorf("expected SubsystemIsInfo=%v, got %v", testCase.expectedIsInfo, got)
			}

			tflog.Trace(ctx, "root trace")
			tflog.Debug(ctx, "root debug")
			tflog.Info(ctx, "root info")

			out := stdBuf.String()
			for _, msg := range testCase.expectedEmitted {
				if !strings.Contains(out, msg) {
					t.Errorf("expected standard logger output to contain %q, got: %s", msg, out)
				}
			}
			for _, msg := range testCase.expectedSuppressed {
				if strings.Contains(out, msg) {
					t.Errorf("expected standard logger output not to contain %q, got: %s", msg, out)
				}
			}
		})
	}
}
