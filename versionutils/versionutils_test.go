package versionutils

import "testing"

func Test_UnifiedVersion1(t *testing.T) {
	const original = "004.004.004"
	actual := UnifiedVersion(original, 4)
	const expected = "000400040004"
	if actual != expected {
		t.Errorf("Expected: '%s', got '%s'", expected, actual)
	}
}

func Test_UnifiedVersion2(t *testing.T) {
	const original = "004.001.017"
	actual := UnifiedVersion(original, 2)
	const expected = "040117"
	if actual != expected {
		t.Errorf("Expected: '%s', got '%s'", expected, actual)
	}
}

type isAnRCVersionTestData struct {
	input    string
	expected bool
}

func Test_IsAnRCVersion(t *testing.T) {
	var isAnRCVersionTests = []isAnRCVersionTestData{
		{"v3.9.7-saucy", false},
		{"v3.9.7-saucy/", false},
		{"v3.12-rc1-saucy", true},
		{"v3.12-rc1-saucy/", true},
		{"v4.1.9-unstable", false},
		{"v4.1.9-unstable/", false},
		{"v4.0-rc7-vivid", true},
		{"v4.0-rc7-vivid/", true},
		{"v4.4.3-wily", false},
		{"v4.4.3-wily/", false},
		{"v4.5-rc3-wily", true},
		{"v4.5-rc3-wily/", true},
		{"v4.7-rc5", true},
		{"v4.7-rc5/", true},
		{"v4.8-rc1", true},
		{"v4.8-rc1/", true},
		{"linux", false},
		{"asd", false},
		{"", false},
	}

	for _, tt := range isAnRCVersionTests {
		if actual := IsAnRCVersion(tt.input); actual != tt.expected {
			t.Errorf("IsAnRCVersion(%q): Expected: %t, actual %t", tt.input, tt.expected, actual)
		}
	}
}
