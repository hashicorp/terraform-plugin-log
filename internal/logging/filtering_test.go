package logging

import (
	"regexp"
	"testing"

	"github.com/google/go-cmp/cmp"
)

const testLogMsg = "System FOO has caused error BAR because of incorrectly configured BAZ"

func TestShouldOmit(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		lOpts          LoggerOpts
		msg            string
		hclogArgSlices [][]interface{}
		expectedToOmit bool
	}{
		"empty-opts": {
			lOpts: LoggerOpts{},
			msg:   testLogMsg,
			hclogArgSlices: [][]interface{}{
				{
					"k1", "v1",
					"k2", "v2",
				},
			},
			expectedToOmit: false,
		},
		"omit-log-by-key": {
			lOpts: LoggerOpts{
				OmitLogWithFieldKeys: []string{"k2"},
			},
			msg: testLogMsg,
			hclogArgSlices: [][]interface{}{
				{
					"k1", "v1",
					"k2", "v2",
				},
			},
			expectedToOmit: true,
		},
		"no-omit-log-by-key-if-case-mismatches": {
			lOpts: LoggerOpts{
				OmitLogWithFieldKeys: []string{"K2"},
			},
			msg: testLogMsg,
			hclogArgSlices: [][]interface{}{
				{
					"k1", "v1",
					"k2", "v2",
				},
			},
			expectedToOmit: false,
		},
		"do-not-omit-log-by-key": {
			lOpts: LoggerOpts{
				OmitLogWithFieldKeys: []string{"k3"},
			},
			msg: testLogMsg,
			hclogArgSlices: [][]interface{}{
				{
					"k1", "v1",
					"k2", "v2",
				},
			},
			expectedToOmit: false,
		},
		"omit-log-matching-regexp-case-insensitive": {
			lOpts: LoggerOpts{
				OmitLogWithMessageRegexes: []*regexp.Regexp{regexp.MustCompile("(?i)(foo|bar)")},
			},
			msg: testLogMsg,
			hclogArgSlices: [][]interface{}{
				{
					"k1", "v1",
					"k2", "v2",
				},
			},
			expectedToOmit: true,
		},
		"do-not-omit-log-matching-regexp-case-sensitive": {
			lOpts: LoggerOpts{
				OmitLogWithMessageRegexes: []*regexp.Regexp{regexp.MustCompile("(foo|bar)")},
			},
			msg: testLogMsg,
			hclogArgSlices: [][]interface{}{
				{
					"k1", "v1",
					"k2", "v2",
				},
			},
			expectedToOmit: false,
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := testCase.lOpts.ShouldOmit(&testCase.msg, testCase.hclogArgSlices...)

			if got != testCase.expectedToOmit {
				t.Errorf("expected ShouldOmit to return %t, got %t", testCase.expectedToOmit, got)
			}
		})
	}
}

func TestApplyMask(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		lOpts             LoggerOpts
		msg               string
		hclogArgSlices    [][]interface{}
		expectedMsg       string
		expectedArgSlices [][]interface{}
	}{
		"empty-opts": {
			lOpts: LoggerOpts{},
			msg:   testLogMsg,
			hclogArgSlices: [][]interface{}{
				{
					"k1", "v1",
					"k2", "v2",
				},
			},
			expectedMsg: "System FOO has caused error BAR because of incorrectly configured BAZ",
			expectedArgSlices: [][]interface{}{
				{
					"k1", "v1",
					"k2", "v2",
				},
			},
		},
		"mask-log-by-key": {
			lOpts: LoggerOpts{
				MaskFieldValuesWithFieldKeys: []string{"k2"},
			},
			msg: testLogMsg,
			hclogArgSlices: [][]interface{}{
				{
					"k1", "v1",
					"k2", "v2",
				},
			},
			expectedMsg: "System FOO has caused error BAR because of incorrectly configured BAZ",
			expectedArgSlices: [][]interface{}{
				{
					"k1", "v1",
					"k2", "***",
				},
			},
		},
		"no-mask-log-by-key-if-case-mismatches": {
			lOpts: LoggerOpts{
				MaskFieldValuesWithFieldKeys: []string{"K2"},
			},
			msg: testLogMsg,
			hclogArgSlices: [][]interface{}{
				{
					"k1", "v1",
					"k2", "v2",
				},
			},
			expectedMsg: "System FOO has caused error BAR because of incorrectly configured BAZ",
			expectedArgSlices: [][]interface{}{
				{
					"k1", "v1",
					"k2", "v2",
				},
			},
		},
		"mask-log-by-non-even-args-cannot-mask-missing-value": {
			lOpts: LoggerOpts{
				MaskFieldValuesWithFieldKeys: []string{"k2", "k4"},
			},
			msg: testLogMsg,
			hclogArgSlices: [][]interface{}{
				{
					"k1", "v1",
					"k2", "v2",
				},
				{
					"k3", "v3",
					"k4",
				},
			},
			expectedMsg: "System FOO has caused error BAR because of incorrectly configured BAZ",
			expectedArgSlices: [][]interface{}{
				{
					"k1", "v1",
					"k2", "***",
				},
				{
					"k3", "v3",
					"k4",
				},
			},
		},
		"mask-log-by-non-even-args": {
			lOpts: LoggerOpts{
				MaskFieldValuesWithFieldKeys: []string{"k2"},
			},
			msg: testLogMsg,
			hclogArgSlices: [][]interface{}{
				{
					"k1", "v1",
					"k2", "v2",
					"k3", "v3",
					"k4",
				},
			},
			expectedMsg: "System FOO has caused error BAR because of incorrectly configured BAZ",
			expectedArgSlices: [][]interface{}{
				{
					"k1", "v1",
					"k2", "***",
					"k3", "v3",
					"k4",
				},
			},
		},
		"mask-log-matching-regexp-case-insensitive": {
			lOpts: LoggerOpts{
				MaskMessageRegexes: []*regexp.Regexp{regexp.MustCompile("(?i)(foo|bar)")},
			},
			msg: testLogMsg,
			hclogArgSlices: [][]interface{}{
				{
					"k1", "v1",
					"k2", "v2",
				},
			},
			expectedMsg: "System *** has caused error *** because of incorrectly configured BAZ",
			expectedArgSlices: [][]interface{}{
				{
					"k1", "v1",
					"k2", "v2",
				},
			},
		},
		"mask-log-matching-regexp-case-sensitive": {
			lOpts: LoggerOpts{
				MaskMessageRegexes: []*regexp.Regexp{regexp.MustCompile("incorrectly configured BAZ")},
			},
			msg: testLogMsg,
			hclogArgSlices: [][]interface{}{
				{
					"k1", "v1",
					"k2", "v2",
				},
			},
			expectedMsg: "System FOO has caused error BAR because of ***",
			expectedArgSlices: [][]interface{}{
				{
					"k1", "v1",
					"k2", "v2",
				},
			},
		},
		"mask-log-by-key-and-matching-regexp": {
			lOpts: LoggerOpts{
				MaskMessageRegexes:           []*regexp.Regexp{regexp.MustCompile("incorrectly configured BAZ")},
				MaskFieldValuesWithFieldKeys: []string{"k1", "k2"},
			},
			msg: testLogMsg,
			hclogArgSlices: [][]interface{}{
				{
					"k1", "v1",
					"k2", "v2",
				},
			},
			expectedMsg: "System FOO has caused error BAR because of ***",
			expectedArgSlices: [][]interface{}{
				{
					"k1", "***",
					"k2", "***",
				},
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			testCase.lOpts.ApplyMask(&testCase.msg, testCase.hclogArgSlices...)

			if diff := cmp.Diff(testCase.msg, testCase.expectedMsg); diff != "" {
				t.Errorf("unexpected difference detected in log message: %s", diff)
			}

			if diff := cmp.Diff(testCase.hclogArgSlices, testCase.expectedArgSlices); diff != "" {
				t.Errorf("unexpected difference detected in log arguments: %s", diff)
			}
		})
	}
}
