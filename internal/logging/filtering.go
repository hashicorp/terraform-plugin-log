package logging

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/internal/hclogutils"
)

const logMaskingReplacementString = "***"

// ShouldOmit takes a log's *string message and slices of arguments,
// and determines, based on the LoggerOpts configuration, if the
// log should be omitted (i.e. prevent it to be printed on the final writer).
func (lo LoggerOpts) ShouldOmit(msg *string, hclogArgSlices ...[]interface{}) bool {
	// Omit log if any of the configured keys is found
	// either in the logger implied arguments,
	// or in the additional arguments
	if len(lo.OmitLogWithFieldKeys) > 0 {
		for _, args := range hclogArgSlices {
			argKeys := hclogutils.ArgsToKeys(args)
			if argKeysContain(argKeys, lo.OmitLogWithFieldKeys) {
				return true
			}
		}
	}

	// Omit log if any of the configured regexp matches the log message
	if len(lo.OmitLogWithMessageRegexes) > 0 {
		for _, r := range lo.OmitLogWithMessageRegexes {
			if r.MatchString(*msg) {
				return true
			}
		}
	}

	// Omit log if any of the configured strings is contained in the log message
	if len(lo.OmitLogWithMessageStrings) > 0 {
		for _, s := range lo.OmitLogWithMessageStrings {
			if strings.Contains(*msg, s) {
				return true
			}
		}
	}

	return false
}

// ApplyMask takes a log's *string message and slices of arguments,
// and applies masking of keys' values and/or message,
// based on the LoggerOpts configuration.
//
// Note that the given input is changed-in-place by this method.
func (lo LoggerOpts) ApplyMask(msg *string, hclogArgSlices ...[]interface{}) {
	if len(lo.MaskFieldValuesWithFieldKeys) > 0 {
		for _, k := range lo.MaskFieldValuesWithFieldKeys {
			for _, args := range hclogArgSlices {
				// Here we loop `i` with steps of 2, starting from position 1 (i.e. `1, 3, 5, 7...`).
				// We then look up the key for each argument, by looking at `i-1`.
				// This ensures that in case of malformed arg slices that don't have
				// an even number of elements, we simply skip the last k/v pair.
				for i := 1; i < len(args); i += 2 {
					switch argK := args[i-1].(type) {
					case string:
						if k == argK {
							args[i] = logMaskingReplacementString
						}
					default:
						if k == fmt.Sprintf("%s", argK) {
							args[i] = logMaskingReplacementString
						}
					}
				}
			}
		}
	}

	// Replace any part of the log message matching any of the configured regexp,
	// with a masking replacement string
	if len(lo.MaskMessageRegexes) > 0 {
		for _, r := range lo.MaskMessageRegexes {
			*msg = r.ReplaceAllString(*msg, logMaskingReplacementString)
		}
	}

	// Replace any part of the log message equal to any of the configured strings,
	// with a masking replacement string
	if len(lo.MaskMessageStrings) > 0 {
		for _, s := range lo.MaskMessageStrings {
			*msg = strings.ReplaceAll(*msg, s, logMaskingReplacementString)
		}
	}
}

func argKeysContain(haystack []string, needles []string) bool {
	for _, h := range haystack {
		for _, n := range needles {
			if n == h {
				return true
			}
		}
	}

	return false
}
