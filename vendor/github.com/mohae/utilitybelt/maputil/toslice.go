// Package maputil provides some helper functions for working with maps.
//
// Currently it is limited to taking maps and returns them as slices of keys and
// values with their indexes corresponding.  Currently only supports maps with strings
// as keys and a limited amount of value types.  Support will be added for other map
// types as needed.
package maputil

// StringInterfaceToSlices takes a map[string]interface{} and returns slices of its keys and
// values with their indexes matching
func StringInterfaceToSlices(m map[string]interface{}) (keys []string, values []interface{}) {
	if m == nil {
		return nil, nil
	}
	for k, v := range m {
		keys = append(keys, k)
		values = append(values, v)
	}
	return keys, values
}

// StringStringToSlices takes a map[string]string and returns its keys and values as
// string slices.
func StringStringToSlices(m map[string]string) (keys, values []string) {
	if m == nil {
		return nil, nil
	}
	for k, v := range m {
		keys = append(keys, k)
		values = append(values, v)
	}
	return keys, values
}

// StringIntToSlices takes a map[string]int and returns its keys and values as
// int slices.
func StringIntToSlices(m map[string]int) (keys []string, values []int) {
	if m == nil {
		return nil, nil
	}
	for k, v := range m {
		keys = append(keys, k)
		values = append(values, v)
	}
	return keys, values
}

// StringBoolToSlices takes a map[string]string and returns its keys and values as
// bool slices.
func StringBoolToSlices(m map[string]bool) (keys []string, values []bool) {
	if m == nil {
		return nil, nil
	}
	for k, v := range m {
		keys = append(keys, k)
		values = append(values, v)
	}
	return keys, values
}
