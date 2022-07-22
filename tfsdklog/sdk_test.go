package tfsdklog_test

import (
	"bytes"
	"context"
	"regexp"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-log/internal/loggertest"
	"github.com/hashicorp/terraform-plugin-log/tfsdklog"
)

func TestSetField(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		key              string
		value            interface{}
		logMessage       string
		additionalFields []map[string]interface{}
		expectedOutput   []map[string]interface{}
	}{
		"no-log-fields": {
			key:        "test-with-key",
			value:      "test-with-value",
			logMessage: "test message",
			expectedOutput: []map[string]interface{}{
				{
					"@level":        "trace",
					"@message":      "test message",
					"@module":       "sdk",
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
					"@module":                 "sdk",
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
					"@module":        "sdk",
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
			ctx = tfsdklog.SetField(ctx, testCase.key, testCase.value)

			tfsdklog.Trace(ctx, testCase.logMessage, testCase.additionalFields...)

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

func TestTrace(t *testing.T) {
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
					"@module":  "sdk",
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
					"@module":    "sdk",
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
					"@module":    "sdk",
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

			tfsdklog.Trace(ctx, testCase.message, testCase.additionalFields...)

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

func TestDebug(t *testing.T) {
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
					"@module":  "sdk",
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
					"@module":    "sdk",
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
					"@module":    "sdk",
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

			tfsdklog.Debug(ctx, testCase.message, testCase.additionalFields...)

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

func TestInfo(t *testing.T) {
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
					"@module":  "sdk",
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
					"@module":    "sdk",
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
					"@module":    "sdk",
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

			tfsdklog.Info(ctx, testCase.message, testCase.additionalFields...)

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

func TestWarn(t *testing.T) {
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
					"@module":  "sdk",
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
					"@module":    "sdk",
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
					"@module":    "sdk",
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

			tfsdklog.Warn(ctx, testCase.message, testCase.additionalFields...)

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

func TestError(t *testing.T) {
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
					"@module":  "sdk",
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
					"@module":    "sdk",
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
					"@module":    "sdk",
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

			tfsdklog.Error(ctx, testCase.message, testCase.additionalFields...)

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

const testLogMsg = "System FOO has caused error BAR because of incorrectly configured BAZ"

func TestOmitLogWithFieldKeys(t *testing.T) {
	testCases := map[string]struct {
		msg              string
		additionalFields []map[string]interface{}
		omitLogKeys      []string
		expectedOutput   []map[string]interface{}
	}{
		"no-omission": {
			msg: testLogMsg,
			additionalFields: []map[string]interface{}{
				{
					"k1": "v1",
					"k2": "v2",
				},
			},
			omitLogKeys: []string{},
			expectedOutput: []map[string]interface{}{
				{
					"@level":   "warn",
					"@message": "System FOO has caused error BAR because of incorrectly configured BAZ",
					"@module":  "sdk",
					"k1":       "v1",
					"k2":       "v2",
				},
			},
		},
		"no-key-matches": {
			msg: testLogMsg,
			additionalFields: []map[string]interface{}{
				{
					"k1": "v1",
					"k2": "v2",
				},
			},
			omitLogKeys: []string{"k3", "K3"},
			expectedOutput: []map[string]interface{}{
				{
					"@level":   "warn",
					"@message": "System FOO has caused error BAR because of incorrectly configured BAZ",
					"@module":  "sdk",
					"k1":       "v1",
					"k2":       "v2",
				},
			},
		},
		"omit-log-by-key": {
			msg: testLogMsg,
			additionalFields: []map[string]interface{}{
				{
					"k1": "v1",
					"k2": "v2",
				},
			},
			omitLogKeys:    []string{"k1"},
			expectedOutput: nil,
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var outputBuffer bytes.Buffer

			ctx := context.Background()
			ctx = loggertest.SDKRoot(ctx, &outputBuffer)
			ctx = tfsdklog.OmitLogWithFieldKeys(ctx, testCase.omitLogKeys...)

			tfsdklog.Warn(ctx, testCase.msg, testCase.additionalFields...)

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

func TestOmitLogWithMessageRegexes(t *testing.T) {
	testCases := map[string]struct {
		msg                   string
		additionalFields      []map[string]interface{}
		omitLogMatchingRegexp []*regexp.Regexp
		expectedOutput        []map[string]interface{}
	}{
		"no-omission": {
			msg: testLogMsg,
			additionalFields: []map[string]interface{}{
				{
					"k1": "v1",
					"k2": "v2",
				},
			},
			omitLogMatchingRegexp: []*regexp.Regexp{},
			expectedOutput: []map[string]interface{}{
				{
					"@level":   "debug",
					"@message": "System FOO has caused error BAR because of incorrectly configured BAZ",
					"@module":  "sdk",
					"k1":       "v1",
					"k2":       "v2",
				},
			},
		},
		"no-matches": {
			msg: testLogMsg,
			additionalFields: []map[string]interface{}{
				{
					"k1": "v1",
					"k2": "v2",
				},
			},
			omitLogMatchingRegexp: []*regexp.Regexp{regexp.MustCompile("(?i)BaAnAnA")},
			expectedOutput: []map[string]interface{}{
				{
					"@level":   "debug",
					"@message": "System FOO has caused error BAR because of incorrectly configured BAZ",
					"@module":  "sdk",
					"k1":       "v1",
					"k2":       "v2",
				},
			},
		},
		"omit-log-matching-regexp": {
			msg: testLogMsg,
			additionalFields: []map[string]interface{}{
				{
					"k1": "v1",
					"k2": "v2",
				},
			},
			omitLogMatchingRegexp: []*regexp.Regexp{regexp.MustCompile("BAZ$")},
			expectedOutput:        nil,
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var outputBuffer bytes.Buffer

			ctx := context.Background()
			ctx = loggertest.SDKRoot(ctx, &outputBuffer)
			ctx = tfsdklog.OmitLogWithMessageRegexes(ctx, testCase.omitLogMatchingRegexp...)

			tfsdklog.Debug(ctx, testCase.msg, testCase.additionalFields...)

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

func TestOmitLogWithMessageStrings(t *testing.T) {
	testCases := map[string]struct {
		msg                   string
		additionalFields      []map[string]interface{}
		omitLogMatchingString []string
		expectedOutput        []map[string]interface{}
	}{
		"no-omission": {
			msg: testLogMsg,
			additionalFields: []map[string]interface{}{
				{
					"k1": "v1",
					"k2": "v2",
				},
			},
			omitLogMatchingString: []string{},
			expectedOutput: []map[string]interface{}{
				{
					"@level":   "debug",
					"@message": "System FOO has caused error BAR because of incorrectly configured BAZ",
					"@module":  "sdk",
					"k1":       "v1",
					"k2":       "v2",
				},
			},
		},
		"no-matches": {
			msg: testLogMsg,
			additionalFields: []map[string]interface{}{
				{
					"k1": "v1",
					"k2": "v2",
				},
			},
			omitLogMatchingString: []string{"BaAnAnA"},
			expectedOutput: []map[string]interface{}{
				{
					"@level":   "debug",
					"@message": "System FOO has caused error BAR because of incorrectly configured BAZ",
					"@module":  "sdk",
					"k1":       "v1",
					"k2":       "v2",
				},
			},
		},
		"omit-log-matching-string": {
			msg: testLogMsg,
			additionalFields: []map[string]interface{}{
				{
					"k1": "v1",
					"k2": "v2",
				},
			},
			omitLogMatchingString: []string{"BAZ"},
			expectedOutput:        nil,
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var outputBuffer bytes.Buffer

			ctx := context.Background()
			ctx = loggertest.SDKRoot(ctx, &outputBuffer)
			ctx = tfsdklog.OmitLogWithMessageStrings(ctx, testCase.omitLogMatchingString...)

			tfsdklog.Debug(ctx, testCase.msg, testCase.additionalFields...)

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

func TestMaskFieldValuesWithFieldKeys(t *testing.T) {
	testCases := map[string]struct {
		msg              string
		additionalFields []map[string]interface{}
		maskLogKeys      []string
		expectedOutput   []map[string]interface{}
	}{
		"no-masking": {
			msg: testLogMsg,
			additionalFields: []map[string]interface{}{
				{
					"k1": "v1",
					"k2": "v2",
				},
			},
			maskLogKeys: []string{},
			expectedOutput: []map[string]interface{}{
				{
					"@level":   "error",
					"@message": "System FOO has caused error BAR because of incorrectly configured BAZ",
					"@module":  "sdk",
					"k1":       "v1",
					"k2":       "v2",
				},
			},
		},
		"no-key-matches": {
			msg: testLogMsg,
			additionalFields: []map[string]interface{}{
				{
					"k1": "v1",
					"k2": "v2",
				},
			},
			maskLogKeys: []string{"k3", "K3"},
			expectedOutput: []map[string]interface{}{
				{
					"@level":   "error",
					"@message": "System FOO has caused error BAR because of incorrectly configured BAZ",
					"@module":  "sdk",
					"k1":       "v1",
					"k2":       "v2",
				},
			},
		},
		"mask-log-by-key": {
			msg: testLogMsg,
			additionalFields: []map[string]interface{}{
				{
					"k1": "v1",
					"k2": "v2",
				},
			},
			maskLogKeys: []string{"k1", "k2", "k3"},
			expectedOutput: []map[string]interface{}{
				{
					"@level":   "error",
					"@message": "System FOO has caused error BAR because of incorrectly configured BAZ",
					"@module":  "sdk",
					"k1":       "***",
					"k2":       "***",
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
			ctx = tfsdklog.MaskFieldValuesWithFieldKeys(ctx, testCase.maskLogKeys...)

			tfsdklog.Error(ctx, testCase.msg, testCase.additionalFields...)

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

func TestMaskAllFieldValuesRegexes(t *testing.T) {
	testCases := map[string]struct {
		msg              string
		additionalFields []map[string]interface{}
		expressions      []*regexp.Regexp
		expectedOutput   []map[string]interface{}
	}{
		"no-masking": {
			msg: testLogMsg,
			additionalFields: []map[string]interface{}{
				{
					"k1": "v1",
					"k2": "v2",
				},
			},
			expressions: []*regexp.Regexp{},
			expectedOutput: []map[string]interface{}{
				{
					"@level":   "trace",
					"@message": "System FOO has caused error BAR because of incorrectly configured BAZ",
					"@module":  "sdk",
					"k1":       "v1",
					"k2":       "v2",
				},
			},
		},
		"no-matches": {
			msg: testLogMsg,
			additionalFields: []map[string]interface{}{
				{
					"k1": "v1",
					"k2": "v2",
				},
			},
			expressions: []*regexp.Regexp{regexp.MustCompile("v3"), regexp.MustCompile("v4")},
			expectedOutput: []map[string]interface{}{
				{
					"@level":   "trace",
					"@message": "System FOO has caused error BAR because of incorrectly configured BAZ",
					"@module":  "sdk",
					"k1":       "v1",
					"k2":       "v2",
				},
			},
		},
		"mask-matching-regexp": {
			msg: testLogMsg,
			additionalFields: []map[string]interface{}{
				{
					"k1": "v1 plus text",
					"k2": "v2 more text",
				},
			},
			expressions: []*regexp.Regexp{regexp.MustCompile("v1|v2")},
			expectedOutput: []map[string]interface{}{
				{
					"@level":   "trace",
					"@message": "System FOO has caused error BAR because of incorrectly configured BAZ",
					"@module":  "sdk",
					"k1":       "*** plus text",
					"k2":       "*** more text",
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
			ctx = tfsdklog.MaskAllFieldValuesRegexes(ctx, testCase.expressions...)

			tfsdklog.Trace(ctx, testCase.msg, testCase.additionalFields...)

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

func TestMaskAllFieldValuesStrings(t *testing.T) {
	testCases := map[string]struct {
		msg              string
		additionalFields []map[string]interface{}
		matchingStrings  []string
		expectedOutput   []map[string]interface{}
	}{
		"no-masking": {
			msg: testLogMsg,
			additionalFields: []map[string]interface{}{
				{
					"k1": "v1",
					"k2": "v2",
				},
			},
			matchingStrings: []string{},
			expectedOutput: []map[string]interface{}{
				{
					"@level":   "trace",
					"@message": "System FOO has caused error BAR because of incorrectly configured BAZ",
					"@module":  "sdk",
					"k1":       "v1",
					"k2":       "v2",
				},
			},
		},
		"no-matches": {
			msg: testLogMsg,
			additionalFields: []map[string]interface{}{
				{
					"k1": "v1",
					"k2": "v2",
				},
			},
			matchingStrings: []string{"v3", "v4"},
			expectedOutput: []map[string]interface{}{
				{
					"@level":   "trace",
					"@message": "System FOO has caused error BAR because of incorrectly configured BAZ",
					"@module":  "sdk",
					"k1":       "v1",
					"k2":       "v2",
				},
			},
		},
		"mask-matching-strings": {
			msg: testLogMsg,
			additionalFields: []map[string]interface{}{
				{
					"k1": "v1 plus text",
					"k2": "v2 more text",
				},
			},
			matchingStrings: []string{"v1", "v2"},
			expectedOutput: []map[string]interface{}{
				{
					"@level":   "trace",
					"@message": "System FOO has caused error BAR because of incorrectly configured BAZ",
					"@module":  "sdk",
					"k1":       "*** plus text",
					"k2":       "*** more text",
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
			ctx = tfsdklog.MaskAllFieldValuesStrings(ctx, testCase.matchingStrings...)

			tfsdklog.Trace(ctx, testCase.msg, testCase.additionalFields...)

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

func TestMaskMessageRegexes(t *testing.T) {
	testCases := map[string]struct {
		msg                   string
		additionalFields      []map[string]interface{}
		maskLogMatchingRegexp []*regexp.Regexp
		expectedOutput        []map[string]interface{}
	}{
		"no-masking": {
			msg: testLogMsg,
			additionalFields: []map[string]interface{}{
				{
					"k1": "v1",
					"k2": "v2",
				},
			},
			maskLogMatchingRegexp: []*regexp.Regexp{},
			expectedOutput: []map[string]interface{}{
				{
					"@level":   "trace",
					"@message": "System FOO has caused error BAR because of incorrectly configured BAZ",
					"@module":  "sdk",
					"k1":       "v1",
					"k2":       "v2",
				},
			},
		},
		"no-matches": {
			msg: testLogMsg,
			additionalFields: []map[string]interface{}{
				{
					"k1": "v1",
					"k2": "v2",
				},
			},
			maskLogMatchingRegexp: []*regexp.Regexp{regexp.MustCompile("(?i)BaAnAnA")},
			expectedOutput: []map[string]interface{}{
				{
					"@level":   "trace",
					"@message": "System FOO has caused error BAR because of incorrectly configured BAZ",
					"@module":  "sdk",
					"k1":       "v1",
					"k2":       "v2",
				},
			},
		},
		"mask-log-matching-regexp": {
			msg: testLogMsg,
			additionalFields: []map[string]interface{}{
				{
					"k1": "v1",
					"k2": "v2",
				},
			},
			maskLogMatchingRegexp: []*regexp.Regexp{regexp.MustCompile("BAZ$")},
			expectedOutput: []map[string]interface{}{
				{
					"@level":   "trace",
					"@message": "System FOO has caused error BAR because of incorrectly configured ***",
					"@module":  "sdk",
					"k1":       "v1",
					"k2":       "v2",
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
			ctx = tfsdklog.MaskMessageRegexes(ctx, testCase.maskLogMatchingRegexp...)

			tfsdklog.Trace(ctx, testCase.msg, testCase.additionalFields...)

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

func TestMaskMessageStrings(t *testing.T) {
	testCases := map[string]struct {
		msg                   string
		additionalFields      []map[string]interface{}
		maskLogMatchingString []string
		expectedOutput        []map[string]interface{}
	}{
		"no-masking": {
			msg: testLogMsg,
			additionalFields: []map[string]interface{}{
				{
					"k1": "v1",
					"k2": "v2",
				},
			},
			maskLogMatchingString: []string{},
			expectedOutput: []map[string]interface{}{
				{
					"@level":   "info",
					"@message": "System FOO has caused error BAR because of incorrectly configured BAZ",
					"@module":  "sdk",
					"k1":       "v1",
					"k2":       "v2",
				},
			},
		},
		"no-matches": {
			msg: testLogMsg,
			additionalFields: []map[string]interface{}{
				{
					"k1": "v1",
					"k2": "v2",
				},
			},
			maskLogMatchingString: []string{"BaAnAnA"},
			expectedOutput: []map[string]interface{}{
				{
					"@level":   "info",
					"@message": "System FOO has caused error BAR because of incorrectly configured BAZ",
					"@module":  "sdk",
					"k1":       "v1",
					"k2":       "v2",
				},
			},
		},
		"mask-log-matching-string": {
			msg: testLogMsg,
			additionalFields: []map[string]interface{}{
				{
					"k1": "v1",
					"k2": "v2",
				},
			},
			maskLogMatchingString: []string{"incorrectly configured BAZ"},
			expectedOutput: []map[string]interface{}{
				{
					"@level":   "info",
					"@message": "System FOO has caused error BAR because of ***",
					"@module":  "sdk",
					"k1":       "v1",
					"k2":       "v2",
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
			ctx = tfsdklog.MaskMessageStrings(ctx, testCase.maskLogMatchingString...)

			tfsdklog.Info(ctx, testCase.msg, testCase.additionalFields...)

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

func TestMaskLogRegexes(t *testing.T) {
	testCases := map[string]struct {
		msg              string
		additionalFields []map[string]interface{}
		expressions      []*regexp.Regexp
		expectedOutput   []map[string]interface{}
	}{
		"no-masking": {
			msg: testLogMsg,
			additionalFields: []map[string]interface{}{
				{
					"k1": "v1",
					"k2": "v2",
				},
			},
			expressions: []*regexp.Regexp{},
			expectedOutput: []map[string]interface{}{
				{
					"@level":   "trace",
					"@message": "System FOO has caused error BAR because of incorrectly configured BAZ",
					"@module":  "sdk",
					"k1":       "v1",
					"k2":       "v2",
				},
			},
		},
		"no-matches": {
			msg: testLogMsg,
			additionalFields: []map[string]interface{}{
				{
					"k1": "v1",
					"k2": "v2",
				},
			},
			expressions: []*regexp.Regexp{regexp.MustCompile("v3"), regexp.MustCompile("v4"), regexp.MustCompile("(?i)BaAnAnA")},
			expectedOutput: []map[string]interface{}{
				{
					"@level":   "trace",
					"@message": "System FOO has caused error BAR because of incorrectly configured BAZ",
					"@module":  "sdk",
					"k1":       "v1",
					"k2":       "v2",
				},
			},
		},
		"mask-matching-regexp": {
			msg: testLogMsg,
			additionalFields: []map[string]interface{}{
				{
					"k1": "v1 plus text",
					"k2": "v2 more text",
				},
			},
			expressions: []*regexp.Regexp{regexp.MustCompile("v1|v2"), regexp.MustCompile("FOO|BAR|BAZ")},
			expectedOutput: []map[string]interface{}{
				{
					"@level":   "trace",
					"@message": "System *** has caused error *** because of incorrectly configured ***",
					"@module":  "sdk",
					"k1":       "*** plus text",
					"k2":       "*** more text",
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
			ctx = tfsdklog.MaskLogRegexes(ctx, testCase.expressions...)

			tfsdklog.Trace(ctx, testCase.msg, testCase.additionalFields...)

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

func TestMaskLogStrings(t *testing.T) {
	testCases := map[string]struct {
		msg              string
		additionalFields []map[string]interface{}
		matchingStrings  []string
		expectedOutput   []map[string]interface{}
	}{
		"no-masking": {
			msg: testLogMsg,
			additionalFields: []map[string]interface{}{
				{
					"k1": "v1",
					"k2": "v2",
				},
			},
			matchingStrings: []string{},
			expectedOutput: []map[string]interface{}{
				{
					"@level":   "trace",
					"@message": "System FOO has caused error BAR because of incorrectly configured BAZ",
					"@module":  "sdk",
					"k1":       "v1",
					"k2":       "v2",
				},
			},
		},
		"no-matches": {
			msg: testLogMsg,
			additionalFields: []map[string]interface{}{
				{
					"k1": "v1",
					"k2": "v2",
				},
			},
			matchingStrings: []string{"v3", "v4"},
			expectedOutput: []map[string]interface{}{
				{
					"@level":   "trace",
					"@message": "System FOO has caused error BAR because of incorrectly configured BAZ",
					"@module":  "sdk",
					"k1":       "v1",
					"k2":       "v2",
				},
			},
		},
		"mask-matching-strings": {
			msg: testLogMsg,
			additionalFields: []map[string]interface{}{
				{
					"k1": "v1 plus text",
					"k2": "v2 more text",
				},
			},
			matchingStrings: []string{"v1", "v2", "FOO", "BAR", "BAZ"},
			expectedOutput: []map[string]interface{}{
				{
					"@level":   "trace",
					"@message": "System *** has caused error *** because of incorrectly configured ***",
					"@module":  "sdk",
					"k1":       "*** plus text",
					"k2":       "*** more text",
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
			ctx = tfsdklog.MaskLogStrings(ctx, testCase.matchingStrings...)

			tfsdklog.Trace(ctx, testCase.msg, testCase.additionalFields...)

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
