package tfsdklog_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-log/internal/loggertest"
	"github.com/hashicorp/terraform-plugin-log/tfsdklog"
)

const (
	testSubsystem       = "test_subsystem"
	testSubsystemModule = "sdk." + testSubsystem
)

func TestSubsystemWith(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		key                 string
		value               interface{}
		logMessage          string
		logadditionalFields []map[string]interface{}
		expectedOutput      []map[string]interface{}
	}{
		"no-log-fields": {
			key:        "test-with-key",
			value:      "test-with-value",
			logMessage: "test message",
			expectedOutput: []map[string]interface{}{
				{
					"@level":        "trace",
					"@message":      "test message",
					"@module":       testSubsystemModule,
					"test-with-key": "test-with-value",
				},
			},
		},
		"mismatched-with-field": {
			key:        "unfielded-test-with-key",
			logMessage: "test message",
			expectedOutput: []map[string]interface{}{
				{
					"@level":                  "trace",
					"@message":                "test message",
					"@module":                 testSubsystemModule,
					"unfielded-test-with-key": nil,
				},
			},
		},
		"with-and-log-fields": {
			key:        "test-with-key",
			value:      "test-with-value",
			logMessage: "test message",
			additionalFields: []map[string]interface{}{
				{
					"test-log-key-1": "test-log-value-1",
					"test-log-key-2": "test-log-value-2",
					"test-log-key-3": "test-log-value-3",
				},
			},
			expectedOutput: []map[string]interface{}{
				{
					"@level":         "trace",
					"@message":       "test message",
					"@module":        testSubsystemModule,
					"test-log-key-1": "test-log-value-1",
					"test-log-key-2": "test-log-value-2",
					"test-log-key-3": "test-log-value-3",
					"test-with-key":  "test-with-value",
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
			ctx = loggertest.SDKRoot(ctx, &outputBuffer)
			ctx = tfsdklog.NewSubsystem(ctx, testSubsystem)
			ctx = tfsdklog.SubsystemWith(ctx, testSubsystem, testCase.key, testCase.value)

			tfsdklog.SubsystemTrace(ctx, testSubsystem, testCase.logMessage, testCase.additionalFields...)

			got, err := loggertest.MultilineJSONDecode(&outputBuffer)

			if err != nil {
				t.Fatalf("unable to read multiple line JSON: %s", err)
			}

			if diff := cmp.Diff(testCase.expectedOutput, got); diff != "" {
				t.Errorf("unexpected output difference: %s", diff)
			}
		})
	}
}

func TestSubsystemTrace(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		message          string
		additionalFields []map[string]interface{}
		expectedOutput   []map[string]interface{}
	}{
		"no-fields": {
			message: "test message",
			expectedOutput: []map[string]interface{}{
				{
					"@level":   "trace",
					"@message": "test message",
					"@module":  testSubsystemModule,
				},
			},
		},
		"fields-single-map": {
			message: "test message",
			additionalFields: []map[string]interface{}{
				{
					"test-key-1": "test-value-1",
					"test-key-2": "test-value-2",
					"test-key-3": "test-value-3",
				},
			},
			expectedOutput: []map[string]interface{}{
				{
					"@level":     "trace",
					"@message":   "test message",
					"@module":    testSubsystemModule,
					"test-key-1": "test-value-1",
					"test-key-2": "test-value-2",
					"test-key-3": "test-value-3",
				},
			},
		},
		"fields-multiple-maps": {
			message: "test message",
			additionalFields: []map[string]interface{}{
				{
					"test-key-1": "test-value-1-map1",
					"test-key-2": "test-value-2-map1",
					"test-key-3": "test-value-3-map1",
				},
				{
					"test-key-4": "test-value-4-map2",
					"test-key-1": "test-value-1-map2",
					"test-key-5": "test-value-5-map2",
				},
			},
			expectedOutput: []map[string]interface{}{
				{
					"@level":     "trace",
					"@message":   "test message",
					"@module":    testSubsystemModule,
					"test-key-1": "test-value-1-map2",
					"test-key-2": "test-value-2-map1",
					"test-key-3": "test-value-3-map1",
					"test-key-4": "test-value-4-map2",
					"test-key-5": "test-value-5-map2",
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
			ctx = loggertest.SDKRoot(ctx, &outputBuffer)
			ctx = tfsdklog.NewSubsystem(ctx, testSubsystem)

			tfsdklog.SubsystemTrace(ctx, testSubsystem, testCase.message, testCase.additionalFields...)

			got, err := loggertest.MultilineJSONDecode(&outputBuffer)

			if err != nil {
				t.Fatalf("unable to read multiple line JSON: %s", err)
			}

			if diff := cmp.Diff(testCase.expectedOutput, got); diff != "" {
				t.Errorf("unexpected output difference: %s", diff)
			}
		})
	}
}

func TestSubsystemDebug(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		message          string
		additionalFields []map[string]interface{}
		expectedOutput   []map[string]interface{}
	}{
		"no-fields": {
			message: "test message",
			expectedOutput: []map[string]interface{}{
				{
					"@level":   "debug",
					"@message": "test message",
					"@module":  testSubsystemModule,
				},
			},
		},
		"fields-single-map": {
			message: "test message",
			additionalFields: []map[string]interface{}{
				{
					"test-key-1": "test-value-1",
					"test-key-2": "test-value-2",
					"test-key-3": "test-value-3",
				},
			},
			expectedOutput: []map[string]interface{}{
				{
					"@level":     "debug",
					"@message":   "test message",
					"@module":    testSubsystemModule,
					"test-key-1": "test-value-1",
					"test-key-2": "test-value-2",
					"test-key-3": "test-value-3",
				},
			},
		},
		"fields-multiple-maps": {
			message: "test message",
			additionalFields: []map[string]interface{}{
				{
					"test-key-1": "test-value-1-map1",
					"test-key-2": "test-value-2-map1",
					"test-key-3": "test-value-3-map1",
				},
				{
					"test-key-4": "test-value-4-map2",
					"test-key-1": "test-value-1-map2",
					"test-key-5": "test-value-5-map2",
				},
			},
			expectedOutput: []map[string]interface{}{
				{
					"@level":     "debug",
					"@message":   "test message",
					"@module":    testSubsystemModule,
					"test-key-1": "test-value-1-map2",
					"test-key-2": "test-value-2-map1",
					"test-key-3": "test-value-3-map1",
					"test-key-4": "test-value-4-map2",
					"test-key-5": "test-value-5-map2",
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
			ctx = loggertest.SDKRoot(ctx, &outputBuffer)
			ctx = tfsdklog.NewSubsystem(ctx, testSubsystem)

			tfsdklog.SubsystemDebug(ctx, testSubsystem, testCase.message, testCase.additionalFields...)

			got, err := loggertest.MultilineJSONDecode(&outputBuffer)

			if err != nil {
				t.Fatalf("unable to read multiple line JSON: %s", err)
			}

			if diff := cmp.Diff(testCase.expectedOutput, got); diff != "" {
				t.Errorf("unexpected output difference: %s", diff)
			}
		})
	}
}

func TestSubsystemInfo(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		message          string
		additionalFields []map[string]interface{}
		expectedOutput   []map[string]interface{}
	}{
		"no-fields": {
			message: "test message",
			expectedOutput: []map[string]interface{}{
				{
					"@level":   "info",
					"@message": "test message",
					"@module":  testSubsystemModule,
				},
			},
		},
		"fields-single-map": {
			message: "test message",
			additionalFields: []map[string]interface{}{
				{
					"test-key-1": "test-value-1",
					"test-key-2": "test-value-2",
					"test-key-3": "test-value-3",
				},
			},
			expectedOutput: []map[string]interface{}{
				{
					"@level":     "info",
					"@message":   "test message",
					"@module":    testSubsystemModule,
					"test-key-1": "test-value-1",
					"test-key-2": "test-value-2",
					"test-key-3": "test-value-3",
				},
			},
		},
		"fields-multiple-maps": {
			message: "test message",
			additionalFields: []map[string]interface{}{
				{
					"test-key-1": "test-value-1-map1",
					"test-key-2": "test-value-2-map1",
					"test-key-3": "test-value-3-map1",
				},
				{
					"test-key-4": "test-value-4-map2",
					"test-key-1": "test-value-1-map2",
					"test-key-5": "test-value-5-map2",
				},
			},
			expectedOutput: []map[string]interface{}{
				{
					"@level":     "info",
					"@message":   "test message",
					"@module":    testSubsystemModule,
					"test-key-1": "test-value-1-map2",
					"test-key-2": "test-value-2-map1",
					"test-key-3": "test-value-3-map1",
					"test-key-4": "test-value-4-map2",
					"test-key-5": "test-value-5-map2",
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
			ctx = loggertest.SDKRoot(ctx, &outputBuffer)
			ctx = tfsdklog.NewSubsystem(ctx, testSubsystem)

			tfsdklog.SubsystemInfo(ctx, testSubsystem, testCase.message, testCase.additionalFields...)

			got, err := loggertest.MultilineJSONDecode(&outputBuffer)

			if err != nil {
				t.Fatalf("unable to read multiple line JSON: %s", err)
			}

			if diff := cmp.Diff(testCase.expectedOutput, got); diff != "" {
				t.Errorf("unexpected output difference: %s", diff)
			}
		})
	}
}

func TestSubsystemWarn(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		message          string
		additionalFields []map[string]interface{}
		expectedOutput   []map[string]interface{}
	}{
		"no-fields": {
			message: "test message",
			expectedOutput: []map[string]interface{}{
				{
					"@level":   "warn",
					"@message": "test message",
					"@module":  testSubsystemModule,
				},
			},
		},
		"fields-single-map": {
			message: "test message",
			additionalFields: []map[string]interface{}{
				{
					"test-key-1": "test-value-1",
					"test-key-2": "test-value-2",
					"test-key-3": "test-value-3",
				},
			},
			expectedOutput: []map[string]interface{}{
				{
					"@level":     "warn",
					"@message":   "test message",
					"@module":    testSubsystemModule,
					"test-key-1": "test-value-1",
					"test-key-2": "test-value-2",
					"test-key-3": "test-value-3",
				},
			},
		},
		"fields-multiple-maps": {
			message: "test message",
			additionalFields: []map[string]interface{}{
				{
					"test-key-1": "test-value-1-map1",
					"test-key-2": "test-value-2-map1",
					"test-key-3": "test-value-3-map1",
				},
				{
					"test-key-4": "test-value-4-map2",
					"test-key-1": "test-value-1-map2",
					"test-key-5": "test-value-5-map2",
				},
			},
			expectedOutput: []map[string]interface{}{
				{
					"@level":     "warn",
					"@message":   "test message",
					"@module":    testSubsystemModule,
					"test-key-1": "test-value-1-map2",
					"test-key-2": "test-value-2-map1",
					"test-key-3": "test-value-3-map1",
					"test-key-4": "test-value-4-map2",
					"test-key-5": "test-value-5-map2",
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
			ctx = loggertest.SDKRoot(ctx, &outputBuffer)
			ctx = tfsdklog.NewSubsystem(ctx, testSubsystem)

			tfsdklog.SubsystemWarn(ctx, testSubsystem, testCase.message, testCase.additionalFields...)

			got, err := loggertest.MultilineJSONDecode(&outputBuffer)

			if err != nil {
				t.Fatalf("unable to read multiple line JSON: %s", err)
			}

			if diff := cmp.Diff(testCase.expectedOutput, got); diff != "" {
				t.Errorf("unexpected output difference: %s", diff)
			}
		})
	}
}

func TestSubsystemError(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		message          string
		additionalFields []map[string]interface{}
		expectedOutput   []map[string]interface{}
	}{
		"no-fields": {
			message: "test message",
			expectedOutput: []map[string]interface{}{
				{
					"@level":   "error",
					"@message": "test message",
					"@module":  testSubsystemModule,
				},
			},
		},
		"fields-single-map": {
			message: "test message",
			additionalFields: []map[string]interface{}{
				{
					"test-key-1": "test-value-1",
					"test-key-2": "test-value-2",
					"test-key-3": "test-value-3",
				},
			},
			expectedOutput: []map[string]interface{}{
				{
					"@level":     "error",
					"@message":   "test message",
					"@module":    testSubsystemModule,
					"test-key-1": "test-value-1",
					"test-key-2": "test-value-2",
					"test-key-3": "test-value-3",
				},
			},
		},
		"fields-multiple-maps": {
			message: "test message",
			additionalFields: []map[string]interface{}{
				{
					"test-key-1": "test-value-1-map1",
					"test-key-2": "test-value-2-map1",
					"test-key-3": "test-value-3-map1",
				},
				{
					"test-key-4": "test-value-4-map2",
					"test-key-1": "test-value-1-map2",
					"test-key-5": "test-value-5-map2",
				},
			},
			expectedOutput: []map[string]interface{}{
				{
					"@level":     "error",
					"@message":   "test message",
					"@module":    testSubsystemModule,
					"test-key-1": "test-value-1-map2",
					"test-key-2": "test-value-2-map1",
					"test-key-3": "test-value-3-map1",
					"test-key-4": "test-value-4-map2",
					"test-key-5": "test-value-5-map2",
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
			ctx = loggertest.SDKRoot(ctx, &outputBuffer)
			ctx = tfsdklog.NewSubsystem(ctx, testSubsystem)

			tfsdklog.SubsystemError(ctx, testSubsystem, testCase.message, testCase.additionalFields...)

			got, err := loggertest.MultilineJSONDecode(&outputBuffer)

			if err != nil {
				t.Fatalf("unable to read multiple line JSON: %s", err)
			}

			if diff := cmp.Diff(testCase.expectedOutput, got); diff != "" {
				t.Errorf("unexpected output difference: %s", diff)
			}
		})
	}
}
