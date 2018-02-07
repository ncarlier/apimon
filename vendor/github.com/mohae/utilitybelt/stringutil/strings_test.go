package stringutil

import "testing"

func TestSetPrefix(t *testing.T) {
	tests := []struct {
		original string
		prefix   string
		expected string
	}{
		{"Ford Prefect", "", "Ford Prefect"},
		{"", "Mr.", "Mr."},
		{"Ford Prefect", "Mr. ", "Mr. Ford Prefect"},
		{"Mr. Ford Prefect", "Mr. ", "Mr. Ford Prefect"},
		{"Mr. Mr. Ford Prefect", "Mr. ", "Mr. Mr. Ford Prefect"},
	}

	for _, test := range tests {
		ret := SetPrefix(test.original, test.prefix)
		if ret != test.expected {
			t.Errorf("Expected: %s, Got: %s\n", test.expected, ret)
		}
	}
}

func TestSetSuffix(t *testing.T) {
	tests := []struct {
		original string
		suffix   string
		expected string
	}{
		{"Go", "", "Go"},
		{"", "pher", "pher"},
		{"Go", "pher", "Gopher"},
		{"Gopher", "pher", "Gopher"},
		{"Gopherpher", "pher", "Gopherpher"},
	}

	for _, test := range tests {
		ret := SetSuffix(test.original, test.suffix)
		if ret != test.expected {
			t.Errorf("Expected: %s, Got: %s\n", test.expected, ret)
		}
	}
}
