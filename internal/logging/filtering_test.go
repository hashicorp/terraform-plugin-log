package logging_test

import (
	"regexp"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-log/internal/logging"
)

const testLogMsg = "System FOO has caused error BAR because of incorrectly configured BAZ"

func TestShouldOmit(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		lOpts          logging.LoggerOpts
		msg            string
		fieldMaps      []map[string]interface{}
		expectedToOmit bool
	}{
		"empty-opts": {
			lOpts: logging.LoggerOpts{},
			msg:   testLogMsg,
			fieldMaps: []map[string]interface{}{
				{
					"k1": "v1",
					"k2": "v2",
				},
			},
			expectedToOmit: false,
		},
		"omit-log-by-key": {
			lOpts: logging.LoggerOpts{
				OmitLogWithFieldKeys: []string{"k2"},
			},
			msg: testLogMsg,
			fieldMaps: []map[string]interface{}{
				{
					"k1": "v1",
					"k2": "v2",
				},
			},
			expectedToOmit: true,
		},
		"no-omit-log-by-key-if-case-mismatches": {
			lOpts: logging.LoggerOpts{
				OmitLogWithFieldKeys: []string{"K2"},
			},
			msg: testLogMsg,
			fieldMaps: []map[string]interface{}{
				{
					"k1": "v1",
					"k2": "v2",
				},
			},
			expectedToOmit: false,
		},
		"do-not-omit-log-by-key": {
			lOpts: logging.LoggerOpts{
				OmitLogWithFieldKeys: []string{"k3"},
			},
			msg: testLogMsg,
			fieldMaps: []map[string]interface{}{
				{
					"k1": "v1",
					"k2": "v2",
				},
			},
			expectedToOmit: false,
		},
		"omit-log-matching-regexp-case-insensitive": {
			lOpts: logging.LoggerOpts{
				OmitLogWithMessageRegexes: []*regexp.Regexp{regexp.MustCompile("(?i)(foo|bar)")},
			},
			msg: testLogMsg,
			fieldMaps: []map[string]interface{}{
				{
					"k1": "v1",
					"k2": "v2",
				},
			},
			expectedToOmit: true,
		},
		"do-not-omit-log-matching-regexp-case-sensitive": {
			lOpts: logging.LoggerOpts{
				OmitLogWithMessageRegexes: []*regexp.Regexp{regexp.MustCompile("(foo|bar)")},
			},
			msg: testLogMsg,
			fieldMaps: []map[string]interface{}{
				{
					"k1": "v1",
					"k2": "v2",
				},
			},
			expectedToOmit: false,
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := testCase.lOpts.ShouldOmit(&testCase.msg, testCase.fieldMaps...)

			if got != testCase.expectedToOmit {
				t.Errorf("expected ShouldOmit to return %t, got %t", testCase.expectedToOmit, got)
			}
		})
	}
}

func TestApplyMask(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		lOpts             logging.LoggerOpts
		msg               string
		fieldMaps         []map[string]interface{}
		expectedMsg       string
		expectedFieldMaps []map[string]interface{}
	}{
		"empty-opts": {
			lOpts: logging.LoggerOpts{},
			msg:   testLogMsg,
			fieldMaps: []map[string]interface{}{
				{
					"k1": "v1",
					"k2": "v2",
				},
			},
			expectedMsg: "System FOO has caused error BAR because of incorrectly configured BAZ",
			expectedFieldMaps: []map[string]interface{}{
				{
					"k1": "v1",
					"k2": "v2",
				},
			},
		},
		"mask-log-by-key": {
			lOpts: logging.LoggerOpts{
				MaskFieldValuesWithFieldKeys: []string{"k2"},
			},
			msg: testLogMsg,
			fieldMaps: []map[string]interface{}{
				{
					"k1": "v1",
					"k2": "v2",
				},
			},
			expectedMsg: "System FOO has caused error BAR because of incorrectly configured BAZ",
			expectedFieldMaps: []map[string]interface{}{
				{
					"k1": "v1",
					"k2": "***",
				},
			},
		},
		"no-mask-log-by-key-if-case-mismatches": {
			lOpts: logging.LoggerOpts{
				MaskFieldValuesWithFieldKeys: []string{"K2"},
			},
			msg: testLogMsg,
			fieldMaps: []map[string]interface{}{
				{
					"k1": "v1",
					"k2": "v2",
				},
			},
			expectedMsg: "System FOO has caused error BAR because of incorrectly configured BAZ",
			expectedFieldMaps: []map[string]interface{}{
				{
					"k1": "v1",
					"k2": "v2",
				},
			},
		},
		"mask-log-matching-regexp-case-insensitive": {
			lOpts: logging.LoggerOpts{
				MaskMessageRegexes: []*regexp.Regexp{regexp.MustCompile("(?i)(foo|bar)")},
			},
			msg: testLogMsg,
			fieldMaps: []map[string]interface{}{
				{
					"k1": "v1",
					"k2": "v2",
				},
			},
			expectedMsg: "System *** has caused error *** because of incorrectly configured BAZ",
			expectedFieldMaps: []map[string]interface{}{
				{
					"k1": "v1",
					"k2": "v2",
				},
			},
		},
		"mask-log-matching-regexp-case-sensitive": {
			lOpts: logging.LoggerOpts{
				MaskMessageRegexes: []*regexp.Regexp{regexp.MustCompile("incorrectly configured BAZ")},
			},
			msg: testLogMsg,
			fieldMaps: []map[string]interface{}{
				{
					"k1": "v1",
					"k2": "v2",
				},
			},
			expectedMsg: "System FOO has caused error BAR because of ***",
			expectedFieldMaps: []map[string]interface{}{
				{
					"k1": "v1",
					"k2": "v2",
				},
			},
		},
		"mask-log-by-key-and-matching-regexp": {
			lOpts: logging.LoggerOpts{
				MaskMessageRegexes:           []*regexp.Regexp{regexp.MustCompile("incorrectly configured BAZ")},
				MaskFieldValuesWithFieldKeys: []string{"k1", "k2"},
			},
			msg: testLogMsg,
			fieldMaps: []map[string]interface{}{
				{
					"k1": "v1",
					"k2": "v2",
				},
			},
			expectedMsg: "System FOO has caused error BAR because of ***",
			expectedFieldMaps: []map[string]interface{}{
				{
					"k1": "***",
					"k2": "***",
				},
			},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			testCase.lOpts.ApplyMask(&testCase.msg, testCase.fieldMaps...)

			if diff := cmp.Diff(testCase.msg, testCase.expectedMsg); diff != "" {
				t.Errorf("unexpected difference detected in log message: %s", diff)
			}

			if diff := cmp.Diff(testCase.fieldMaps, testCase.expectedFieldMaps); diff != "" {
				t.Errorf("unexpected difference detected in log arguments: %s", diff)
			}
		})
	}
}
