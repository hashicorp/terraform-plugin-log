package tfsdklog

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestParseJSON(t *testing.T) {
	t.Parallel()

	type testCase struct {
		input       string
		expected    *logEntry
		expectedErr error
	}

	tests := map[string]testCase{
		"basic": {
			input: `{"@level":"DEBUG", "@message": "this is a test", "@module": "test", "@timestamp": "2021-09-15T03:07:23.123456Z"}`,
			expected: &logEntry{
				Message:   "this is a test",
				Level:     "DEBUG",
				Timestamp: time.Date(2021, time.September, 15, 3, 7, 23, 123456000, time.UTC),
				Module:    "test",
			},
		},
		"kv-pair": {
			input: `{"@level":"DEBUG", "@message": "this is a test", "@module": "test", "@timestamp": "2021-09-15T03:07:23.123456Z", "kv-pair-1": ["testing", 123]}`,
			expected: &logEntry{
				Message:   "this is a test",
				Level:     "DEBUG",
				Timestamp: time.Date(2021, time.September, 15, 3, 7, 23, 123456000, time.UTC),
				Module:    "test",
				KVPairs: []*logEntryKV{
					{
						Key:   "kv-pair-1",
						Value: []interface{}{"testing", float64(123)},
					},
				},
			},
		},
		"kv-pairs": {
			input: `{"@level":"DEBUG", "@message": "this is a test", "@module": "test", "@timestamp": "2021-09-15T03:07:23.123456Z", "kv-pair-1": ["testing", 123], "kv-pair-2": "foo", "kv-pair-3": 123, "kv-pair-4": {"hello": "world"}}`,
			expected: &logEntry{
				Message:   "this is a test",
				Level:     "DEBUG",
				Timestamp: time.Date(2021, time.September, 15, 3, 7, 23, 123456000, time.UTC),
				Module:    "test",
				KVPairs: []*logEntryKV{
					{
						Key:   "kv-pair-1",
						Value: []interface{}{"testing", float64(123)},
					},
					{
						Key:   "kv-pair-2",
						Value: "foo",
					},
					{
						Key:   "kv-pair-3",
						Value: float64(123),
					},
					{
						Key: "kv-pair-4",
						Value: map[string]interface{}{
							"hello": "world",
						},
					},
				},
			},
		},
		"unterminated-json": {
			input:       `{"@level":`,
			expectedErr: fmt.Errorf("unexpected end of JSON input"),
		},
	}

	for name, tc := range tests {
		name, tc := name, tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := parseJSON([]byte(tc.input))

			if err != nil && tc.expectedErr != nil && err.Error() != tc.expectedErr.Error() {
				t.Fatalf("expected error to be %v, got %v", tc.expectedErr, err)
			} else if err == nil && tc.expectedErr != nil {
				t.Fatalf("expected error to be %v, didn't get an error", tc.expectedErr)
			} else if err != nil && tc.expectedErr == nil {
				t.Fatalf("unexpected error %v", err)
			}

			if err != nil {
				return
			}

			if diff := cmp.Diff(tc.expected, got); diff != "" {
				t.Fatalf("unexpected diff (-wanted, +got): %s", diff)
			}
		})
	}
}
