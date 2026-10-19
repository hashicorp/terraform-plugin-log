// Copyright IBM Corp. 2021, 2026
// SPDX-License-Identifier: MPL-2.0

package tfsdklog_test

import (
	"bytes"
	"context"
	"log"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-log/internal/loggertest"
	"github.com/hashicorp/terraform-plugin-log/internal/logging"
	"github.com/hashicorp/terraform-plugin-log/tfsdklog"
)

func TestWithLevelFromEnv(t *testing.T) {
	testCases := map[string]struct {
		env                      map[string]string
		rootOption               logging.Option
		subsystemOption          logging.Option
		expectedRootOutput       []map[string]interface{}
		expectedSubsystemOutput  []map[string]interface{}
		expectedRootIsTrace      bool
		expectedSubsystemIsTrace bool
	}{
		"default-off-when-all-unset": {
			env:                      map[string]string{},
			rootOption:               tfsdklog.WithLevelFromEnv("TF_LOG_SDK"),
			subsystemOption:          tfsdklog.WithLevelFromEnv("TF_LOG_SDK", "framework"),
			expectedRootOutput:       nil,
			expectedSubsystemOutput:  nil,
			expectedRootIsTrace:      false,
			expectedSubsystemIsTrace: false,
		},
		"tf-log-fallback-debug": {
			env: map[string]string{
				"TF_LOG": "DEBUG",
			},
			rootOption:      tfsdklog.WithLevelFromEnv("TF_LOG_SDK"),
			subsystemOption: tfsdklog.WithLevelFromEnv("TF_LOG_SDK", "framework"),
			expectedRootOutput: []map[string]interface{}{
				{
					"@level":    "debug",
					"@message":  "root debug",
					"@module":   "sdk",
					"tf_req_id": "123",
				},
			},
			expectedSubsystemOutput: []map[string]interface{}{
				{
					"@level":    "debug",
					"@message":  "subsystem debug",
					"@module":   "sdk.framework",
					"tf_req_id": "123",
					"sub_key":   "sub_val",
				},
			},
			expectedRootIsTrace:      false,
			expectedSubsystemIsTrace: false,
		},
		"tf-log-provider-fallback-debug": {
			env: map[string]string{
				"TF_LOG_PROVIDER": "DEBUG",
			},
			rootOption:      tfsdklog.WithLevelFromEnv("TF_LOG_SDK"),
			subsystemOption: tfsdklog.WithLevelFromEnv("TF_LOG_SDK", "framework"),
			expectedRootOutput: []map[string]interface{}{
				{
					"@level":    "debug",
					"@message":  "root debug",
					"@module":   "sdk",
					"tf_req_id": "123",
				},
			},
			expectedSubsystemOutput: []map[string]interface{}{
				{
					"@level":    "debug",
					"@message":  "subsystem debug",
					"@module":   "sdk.framework",
					"tf_req_id": "123",
					"sub_key":   "sub_val",
				},
			},
			expectedRootIsTrace:      false,
			expectedSubsystemIsTrace: false,
		},
		"subsystem-override-trace-when-root-off": {
			env: map[string]string{
				"TF_LOG_SDK_FRAMEWORK": "TRACE",
			},
			rootOption:         tfsdklog.WithLevelFromEnv("TF_LOG_SDK"),
			subsystemOption:    tfsdklog.WithLevelFromEnv("TF_LOG_SDK", "framework"),
			expectedRootOutput: nil,
			expectedSubsystemOutput: []map[string]interface{}{
				{
					"@level":    "trace",
					"@message":  "subsystem trace",
					"@module":   "sdk.framework",
					"tf_req_id": "123",
					"sub_key":   "sub_val",
				},
				{
					"@level":    "debug",
					"@message":  "subsystem debug",
					"@module":   "sdk.framework",
					"tf_req_id": "123",
					"sub_key":   "sub_val",
				},
			},
			expectedRootIsTrace:      false,
			expectedSubsystemIsTrace: true,
		},
		"subsystem-override-off-when-root-trace": {
			env: map[string]string{
				"TF_LOG":               "TRACE",
				"TF_LOG_SDK_FRAMEWORK": "OFF",
			},
			rootOption:      tfsdklog.WithLevelFromEnv("TF_LOG_SDK"),
			subsystemOption: tfsdklog.WithLevelFromEnv("TF_LOG_SDK", "framework"),
			expectedRootOutput: []map[string]interface{}{
				{
					"@level":    "trace",
					"@message":  "root trace",
					"@module":   "sdk",
					"tf_req_id": "123",
				},
				{
					"@level":    "debug",
					"@message":  "root debug",
					"@module":   "sdk",
					"tf_req_id": "123",
				},
			},
			expectedSubsystemOutput:  nil,
			expectedRootIsTrace:      true,
			expectedSubsystemIsTrace: false,
		},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			t.Cleanup(func() {
				var restoreBuf bytes.Buffer
				restoreCtx := loggertest.SDKRoot(context.Background(), &restoreBuf)
				_ = tfsdklog.NewSubsystem(restoreCtx, "framework")
			})

			t.Setenv("TF_LOG", "")
			t.Setenv("TF_LOG_SDK", "")
			t.Setenv("TF_LOG_SDK_FRAMEWORK", "")
			t.Setenv("TF_LOG_PROVIDER", "")
			t.Setenv("TF_ACC_LOG_PATH", "")
			for k, v := range testCase.env {
				t.Setenv(k, v)
			}

			var rootBuf bytes.Buffer
			ctx := context.Background()
			ctx = tfsdklog.NewRootSDKLogger(
				ctx,
				logging.WithoutLocation(),
				logging.WithoutTimestamp(),
				logging.WithOutput(&rootBuf),
				testCase.rootOption,
			)
			ctx = tfsdklog.SetField(ctx, "tf_req_id", "123")
			ctx = tfsdklog.NewSubsystem(
				ctx,
				"framework",
				testCase.subsystemOption,
				tfsdklog.WithRootFields(),
			)
			ctx = tfsdklog.SubsystemSetField(ctx, "framework", "sub_key", "sub_val")

			if got := tfsdklog.IsTrace(ctx); got != testCase.expectedRootIsTrace {
				t.Errorf("expected IsTrace=%v, got %v", testCase.expectedRootIsTrace, got)
			}
			if got := tfsdklog.SubsystemIsTrace(ctx, "framework"); got != testCase.expectedSubsystemIsTrace {
				t.Errorf("expected SubsystemIsTrace=%v, got %v", testCase.expectedSubsystemIsTrace, got)
			}

			tfsdklog.Trace(ctx, "root trace")
			tfsdklog.Debug(ctx, "root debug")

			gotRoot, err := loggertest.MultilineJSONDecode(&rootBuf)
			if err != nil {
				t.Fatalf("unable to read root JSON: %s", err)
			}
			if diff := cmp.Diff(testCase.expectedRootOutput, gotRoot); diff != "" {
				t.Errorf("unexpected root output difference: %s", diff)
			}

			rootBuf.Reset()
			tfsdklog.SubsystemTrace(ctx, "framework", "subsystem trace")
			tfsdklog.SubsystemDebug(ctx, "framework", "subsystem debug")

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
	t.Cleanup(func() {
		var restoreBuf bytes.Buffer
		restoreCtx := loggertest.SDKRoot(context.Background(), &restoreBuf)
		_ = tfsdklog.NewSubsystem(restoreCtx, "sink_sub")
	})

	t.Run("ContextWithTestLogging-honors-TF_LOG_PROVIDER", func(t *testing.T) {
		t.Setenv("TF_LOG", "")
		t.Setenv("TF_LOG_SDK", "")
		t.Setenv("TF_LOG_PROVIDER", "DEBUG")
		t.Setenv("TF_ACC_LOG_PATH", "")

		ctx := tfsdklog.ContextWithTestLogging(context.Background(), t.Name())
		ctx = tfsdklog.NewRootSDKLogger(ctx, tfsdklog.WithLevelFromEnv("TF_LOG_SDK"))
		ctx = tfsdklog.NewSubsystem(ctx, "sink_sub", tfsdklog.WithLevelFromEnv("TF_LOG_SDK", "sink_sub"))

		if tfsdklog.IsTrace(ctx) {
			t.Errorf("expected IsTrace=false when TF_LOG_PROVIDER=DEBUG")
		}
		if !tfsdklog.IsDebug(ctx) {
			t.Errorf("expected IsDebug=true when TF_LOG_PROVIDER=DEBUG")
		}
		if tfsdklog.SubsystemIsTrace(ctx, "sink_sub") {
			t.Errorf("expected SubsystemIsTrace=false when TF_LOG_PROVIDER=DEBUG")
		}
		if !tfsdklog.SubsystemIsDebug(ctx, "sink_sub") {
			t.Errorf("expected SubsystemIsDebug=true when TF_LOG_PROVIDER=DEBUG")
		}
	})

	t.Run("ContextWithStandardLogging-preserves-default-info-when-env-unset", func(t *testing.T) {
		t.Setenv("TF_LOG", "")
		t.Setenv("TF_LOG_SDK", "")
		t.Setenv("TF_LOG_PROVIDER", "")
		t.Setenv("TF_ACC_LOG_PATH", "")

		var stdBuf bytes.Buffer
		origOut := log.Default().Writer()
		log.Default().SetOutput(&stdBuf)
		t.Cleanup(func() {
			log.Default().SetOutput(origOut)
		})

		ctx := tfsdklog.ContextWithStandardLogging(context.Background(), t.Name())
		ctx = tfsdklog.NewRootSDKLogger(ctx, tfsdklog.WithLevelFromEnv("TF_LOG_SDK"))
		ctx = tfsdklog.NewSubsystem(ctx, "sink_sub", tfsdklog.WithLevelFromEnv("TF_LOG_SDK", "sink_sub"))

		tfsdklog.Debug(ctx, "root debug should be filtered by default info")
		tfsdklog.Info(ctx, "root info should be emitted")

		if !tfsdklog.SubsystemIsInfo(ctx, "sink_sub") {
			t.Errorf("expected SubsystemIsInfo=true when standard logging sink is active")
		}

		out := stdBuf.String()
		if strings.Contains(out, "should be filtered") {
			t.Errorf("unexpected debug output in standard logger sink: %s", out)
		}
		if !strings.Contains(out, "root info should be emitted") {
			t.Errorf("expected info output in standard logger sink, got: %s", out)
		}
	})
}

func BenchmarkSubsystemTraceDisabledFromEnv(b *testing.B) {
	b.Cleanup(func() {
		var restoreBuf bytes.Buffer
		restoreCtx := loggertest.SDKRoot(context.Background(), &restoreBuf)
		_ = tfsdklog.NewSubsystem(restoreCtx, "bench_framework")
	})

	b.Setenv("TF_LOG", "")
	b.Setenv("TF_LOG_SDK", "")
	b.Setenv("TF_LOG_SDK_FRAMEWORK", "")
	b.Setenv("TF_LOG_PROVIDER", "")
	b.Setenv("TF_ACC_LOG_PATH", "")

	var buf bytes.Buffer
	ctx := context.Background()
	ctx = tfsdklog.NewRootSDKLogger(
		ctx,
		logging.WithoutLocation(),
		logging.WithoutTimestamp(),
		logging.WithOutput(&buf),
		tfsdklog.WithLevelFromEnv("TF_LOG_SDK"),
	)
	ctx = tfsdklog.NewSubsystem(
		ctx,
		"bench_framework",
		tfsdklog.WithLevelFromEnv("TF_LOG_SDK_FRAMEWORK"),
	)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx = tfsdklog.SubsystemSetField(ctx, "bench_framework", "tf_req_id", "12345")
		tfsdklog.SubsystemTrace(ctx, "bench_framework", "Checking SchemaTypes")
	}
}
