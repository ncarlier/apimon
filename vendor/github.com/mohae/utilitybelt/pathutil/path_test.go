package pathutil

import (
	"testing"
)

type mt struct{}

func TestDirDirWalk(t *testing.T) {
	tests := []struct {
		path        string
		expected    map[string]mt
		expectedErr string
	}{
		{
			path:        "invalid",
			expected:    nil,
			expectedErr: "invalid does not exist",
		},
		{
			path: "../test_files/pixies",
			expected: map[string]mt{
				"I-cant-forget.txt":            {},
				"Ive-been-waiting-for-you.txt": {},
				"born-in-chicago.txt":          {},
				"debaser-I.txt":                {},
				"debaser-II.txt":               {},
				"gigantic.txt":                 {},
				"something-against-you.txt":    {},
				"where-is-my-mind.txt":         {},
				"surfer-rosa":                  {},
				"doolittle":                    {},
			},
			expectedErr: "",
		},
		{
			path: "../test_files/pink-floyd",
			expected: map[string]mt{
				"echos.txt":                    {},
				"one-of-these-days.txt":        {},
				"san-tropez.txt":               {},
				"seamus.txt":                   {},
				"I-V.txt":                      {},
				"VI-IX.txt":                    {},
				"have-a-cigar.txt":             {},
				"wish-you-were-here-part1.txt": {},
				"wish-you-were-here-part2.txt": {},
				"wish-you-were-here-part3.txt": {},
				"wish-you-were-here":           {},
				"meddle":                       {},
				"shine-on-you-crazy-diamond":   {},
			},
			expectedErr: "",
		},
		{
			path: "../test_files",
			expected: map[string]mt{
				"I-cant-forget.txt":                      {},
				"Ive-been-waiting-for-you.txt":           {},
				"born-in-chicago.txt":                    {},
				"debaser-I.txt":                          {},
				"debaser-II.txt":                         {},
				"gigantic.txt":                           {},
				"something-against-you.txt":              {},
				"where-is-my-mind.txt":                   {},
				"echos.txt":                              {},
				"one-of-these-days.txt":                  {},
				"san-tropez.txt":                         {},
				"seamus.txt":                             {},
				"I-V.txt":                                {},
				"VI-IX.txt":                              {},
				"have-a-cigar.txt":                       {},
				"wish-you-were-here-part1.txt":           {},
				"wish-you-were-here-part2.txt":           {},
				"wish-you-were-here-part3.txt":           {},
				"tmbg-ana-ng.txt":                        {},
				"tmbg-particle-man.txt":                  {},
				"tmbg-sapphire-bullets-of-pure-love.txt": {},
				"pixies":                     {},
				"pink-floyd":                 {},
				"surfer-rosa":                {},
				"doolittle":                  {},
				"wish-you-were-here":         {},
				"meddle":                     {},
				"shine-on-you-crazy-diamond": {},
			},
			expectedErr: "",
		},
	}

	for i, test := range tests {
		d := &Dir{Files: []file{}}
		err := d.Walk(test.path)
		if err != nil {
			if err.Error() != test.expectedErr {
				t.Errorf("%d: expected %q got %q", i, test.expectedErr, err)
			}
			continue
		}
		if test.expectedErr != "" {
			t.Errorf("%d: xpected error %s", i, err)
			continue
		}
		for _, f := range d.Files {
			if _, ok := test.expected[f.Info.Name()]; !ok {
				t.Errorf("%d: %s was indexed but not found in the expected filename list", i, f.Info.Name())
			}
		}
	}
}

func TestPathExists(t *testing.T) {
	tests := []struct {
		path     string
		expected bool
		errS     string
	}{
		{"", false, ""},
		{"../test_files", true, ""},
		{"../test_files/pixies/born-in-chicago.txt", true, ""},
		{"../tst", false, ""},
		{"../test_files/pink-floyd/animals", false, ""},
	}

	for i, test := range tests {
		exists, err := PathExists(test.path)
		if err != nil {
			if err.Error() != test.errS {
				t.Errorf("%d: expected %q, got %q", i, test.errS, err.Error())
			}
			continue
		}
		if test.errS != "" {
			t.Errorf("%d: %s was expected, but no error was encountered.", i, test.errS)
			continue
		}
		if exists != test.expected {
			t.Errorf("%d: expected %v got %v for %s", i, test.expected, exists, test.path)
		}
	}
}
