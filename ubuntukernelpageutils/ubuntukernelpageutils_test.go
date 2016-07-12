package ubuntukernelpageutils

import "testing"

func equalStringSlices(a, b []string) bool {
	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

type removeDuplicatesTestData struct {
	str      []string
	expected []string
}

func Test_RemoveDuplicates(t *testing.T) {
	var removeDuplicatesTests = []removeDuplicatesTestData{
		{[]string{"aa", "bb", "cc", "aa", "aa", "aa", "aa", "aa"}, []string{"aa", "bb", "cc"}},
		{[]string{"aa"}, []string{"aa"}},
		{[]string{"aa", "bb", "cc", "aa"}, []string{"aa", "bb", "cc"}},
		{[]string{"aa", "bb", "cc", "aa", "aa", "aa", "123"}, []string{"aa", "bb", "cc", "123"}},
		{[]string{}, []string{}},
	}

	for _, tt := range removeDuplicatesTests {
		actual := removeDuplicates(tt.str)
		if !equalStringSlices(actual, tt.expected) {
			t.Errorf("removeDuplicates(%q): Expected: %q, actual %q", tt.str, tt.expected, actual)
		}
	}
}
