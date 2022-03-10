package tflog_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-log/internal/loggertest"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func testSubsystemTraceHelper(ctx context.Context, message string) {
	tflog.SubsystemTrace(ctx, testSubsystem, message)
}

func TestWithAdditionalLocationOffset(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		additionalLocationOffset int
		logImpl                  func(context.Context)
		expectedOutput           []map[string]interface{}
	}{
		"0-no-helper": {
			additionalLocationOffset: 0,
			logImpl: func(ctx context.Context) {
				tflog.SubsystemTrace(ctx, testSubsystem, "test message")
			},
			expectedOutput: []map[string]interface{}{
				{
					// Caller line (number after colon) should match
					// tflog.SubsystemTrace() line in test case implementation.
					"@caller":  "/tflog/options_test.go:30",
					"@level":   hclog.Trace.String(),
					"@message": "test message",
					"@module":  testSubsystemModule,
				},
			},
		},
		"0-one-helper": {
			additionalLocationOffset: 0,
			logImpl: func(ctx context.Context) {
				testSubsystemTraceHelper(ctx, "test message")
			},
			expectedOutput: []map[string]interface{}{
				{
					// Caller line (number after colon) should match
					// tflog.SubsystemTrace() line in testSubsystemTraceHelper
					// function implementation.
					"@caller":  "/tflog/options_test.go:16",
					"@level":   hclog.Trace.String(),
					"@message": "test message",
					"@module":  testSubsystemModule,
				},
			},
		},
		"1-one-helper": {
			additionalLocationOffset: 1,
			logImpl: func(ctx context.Context) {
				testSubsystemTraceHelper(ctx, "test message")
			},
			expectedOutput: []map[string]interface{}{
				{
					// Caller line (number after colon) should match
					// testSubsystemTraceHelper() line in test case
					// implementation.
					"@caller":  "/tflog/options_test.go:63",
					"@level":   hclog.Trace.String(),
					"@message": "test message",
					"@module":  testSubsystemModule,
				},
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var outputBuffer bytes.Buffer

			ctx := context.Background()
			ctx = loggertest.ProviderRootWithLocation(ctx, &outputBuffer)
			ctx = tflog.NewSubsystem(ctx, testSubsystem, tflog.WithAdditionalLocationOffset(testCase.additionalLocationOffset))

			testCase.logImpl(ctx)

			got, err := loggertest.MultilineJSONDecode(&outputBuffer)

			if err != nil {
				t.Fatalf("unable to read multiple line JSON: %s", err)
			}

			// Strip non-deterministic caller information up to this package, e.g.
			// /Users/example/src/github.com/hashicorp/terraform-plugin-log/tflog/...
			for _, gotEntry := range got {
				caller, ok := gotEntry["@caller"].(string)

				if !ok {
					continue
				}

				packageIndex := strings.Index(caller, "/tflog/")

				if packageIndex == -1 {
					continue
				}

				gotEntry["@caller"] = caller[packageIndex:]
			}

			if diff := cmp.Diff(testCase.expectedOutput, got); diff != "" {
				t.Errorf("unexpected output difference: %s", diff)
			}
		})
	}
}
