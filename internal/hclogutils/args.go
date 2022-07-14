package hclogutils

import (
	"fmt"
)

// MapsToArgs will shallow merge field maps into a slice of key/value pairs
// arguments (i.e. `[k1, v1, k2, v2, ...]`) expected by hc-log.Logger methods.
func MapsToArgs(maps ...map[string]interface{}) []interface{} {
	switch len(maps) {
	case 0:
		return nil
	case 1:
		result := make([]interface{}, 0, len(maps[0])*2)

		for k, v := range maps[0] {
			result = append(result, k, v)
		}

		return result
	default:
		// Pre-allocate a map to merge all the maps into,
		// that has at least the capacity equivalent to the number
		// of maps to merge
		mergedMap := make(map[string]interface{}, len(maps))

		// Merge all the maps into one;
		// in case of clash, only the last key is preserved
		for _, m := range maps {
			for k, v := range m {
				mergedMap[k] = v
			}
		}

		// As we have merged all maps into one, we can use this
		// same function recursively for the `switch case 1`.
		return MapsToArgs(mergedMap)
	}
}

// ArgsToKeys will extract all keys from a slice of key/value pairs
// arguments (i.e. `[k1, v1, k2, v2, ...]`) expected by hc-log.Logger methods.
//
// Note that, in case of an odd number of arguments, the last key captured
// will refer to a value that does not actually exist.
func ArgsToKeys(args []interface{}) []string {
	// Pre-allocate enough capacity to fit all the keys,
	// i.e. all the elements in the input array in even position
	keys := make([]string, 0, len(args)/2)

	for i := 0; i < len(args); i += 2 {
		// All keys should be strings, but in case they are not
		// we format them to string
		switch k := args[i].(type) {
		case string:
			keys = append(keys, k)
		default:
			keys = append(keys, fmt.Sprintf("%s", k))
		}
	}

	return keys
}
