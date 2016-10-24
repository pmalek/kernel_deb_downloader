package sortedmapstringstring

import (
	"sort"
	"testing"
	"testing/quick"
)

func stringSliceEq(a, b []string) bool {

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

type sortedKeyTestData struct {
	input    map[string]string
	expected []string
}

func Test_SortedKeysByValue_returnsSliceInSortedOrder(t *testing.T) {
	var sortedKeyTests = []sortedKeyTestData{
		{map[string]string{"key1": "value1", "key2": "value2"}, []string{"key1", "key2"}},
		{map[string]string{"1": "value1", "2": "value2"}, []string{"1", "2"}},
		{map[string]string{}, []string{}},
		{map[string]string{"1": "1", "2": "2", "3": "3"}, []string{"1", "2", "3"}},
		{map[string]string{"1": "9", "2": "4", "3": "8"}, []string{"2", "3", "1"}},
		{map[string]string{"1": "2", "2": "1", "3": "5"}, []string{"2", "1", "3"}},
		{map[string]string{"1": "1", "2": "2", "3": "3", "4": "4"}, []string{"1", "2", "3", "4"}},
		{map[string]string{"1": "xyz", "2": "asd", "3": "aaa", "4": "ppp"}, []string{"3", "2", "4", "1"}},
		{map[string]string{"1": "1", "2": "2", "3": "3", "4": "4", "5": "5", "6": "6", "7": "7", "8": "8", "9": "9"}, []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}},
	}

	for _, tt := range sortedKeyTests {
		actual := SortedMapKeysByValue(tt.input)
		if !stringSliceEq(actual, tt.expected) {
			t.Errorf("SortedKeys(%q): Expected: %q, actual %q", tt.input, tt.expected, actual)
		}
	}
}

var cfg = &quick.Config{
	MaxCount: 500,
}

func Test_SortedKeysByValue_quick(t *testing.T) {
	f := func(input map[string]string) bool {
		var outSortedKeys = SortedMapKeysByValue(input)
		return len(outSortedKeys) == len(input)
	}

	if err := quick.Check(f, cfg); err != nil {
		t.Error(err)
	}
}

func Test_SortedKeysByKey_returnsSliceInSortedOrder(t *testing.T) {
	var sortedKeyTests = []sortedKeyTestData{
		{map[string]string{}, []string{}},
		{map[string]string{"": "value1"}, []string{""}},
		{map[string]string{"": "909!!!"}, []string{""}},
		{map[string]string{"key1": "value1"}, []string{"key1"}},
		{map[string]string{"key1": "value1", "key2": "value2"}, []string{"key1", "key2"}},
		{map[string]string{"key1": "value8", "key2": "value2"}, []string{"key1", "key2"}},
		{map[string]string{"key1": "WHETEVER", "key2": "WHETEVER"}, []string{"key1", "key2"}},
		{map[string]string{"key2": "WHETEVER", "key1": "WHETEVER"}, []string{"key1", "key2"}},
		{map[string]string{"1": "1", "2": "2", "3": "3"}, []string{"1", "2", "3"}},
		{map[string]string{"1": "ALSKDJ", "2": "ALSKDJ", "3": "3sldkasj"}, []string{"1", "2", "3"}},
		{map[string]string{"1": "2222", "2": "9999", "3": "3sdlkasj"}, []string{"1", "2", "3"}},
		{map[string]string{"040801": "", "040802": "", "040803": "", "040800": ""}, []string{"040800", "040801", "040802", "040803"}},
	}

	for _, tt := range sortedKeyTests {
		actual := SortedMapKeysByKey(tt.input)
		if !stringSliceEq(actual, tt.expected) {
			t.Errorf("SortedKeys(%q): Expected: %q, actual %q", tt.input, tt.expected, actual)
		}
	}
}

func Test_SortedKeysByKey_quick_lenStaysTheSame(t *testing.T) {
	f := func(input map[string]string) bool {
		var outSortedKeys = SortedMapKeysByKey(input)
		return len(outSortedKeys) == len(input)
	}

	if err := quick.Check(f, cfg); err != nil {
		t.Error(err)
	}
}

func Test_SortedKeysByKey_quick_keysAreAlwaysSorted(t *testing.T) {
	f := func(input map[string]string) bool {
		var outSortedKeys = SortedMapKeysByKey(input)
		return sort.StringsAreSorted(outSortedKeys)
	}

	if err := quick.Check(f, cfg); err != nil {
		t.Error(err)
	}
}
