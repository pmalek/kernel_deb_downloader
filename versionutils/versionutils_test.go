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

func Test_IsAnRcVersion1(t *testing.T) {
	actual := IsAnRCVersion("v4.5-rc3-wily/")
	const expected = true
	if actual != expected {
		t.Errorf("Expected: '%s', got '%s'", expected, actual)
	}
}

func Test_IsAnRcVersion2(t *testing.T) {
	actual := IsAnRCVersion("v4.4.3-wily/")
	const expected = false
	if actual != expected {
		t.Errorf("Expected: '%s', got '%s'", expected, actual)
	}
}

func Test_IsAnRcVersion3(t *testing.T) {
	actual := IsAnRCVersion("v4.0-rc7-vivid/")
	const expected = true
	if actual != expected {
		t.Errorf("Expected: '%s', got '%s'", expected, actual)
	}
}

func Test_IsAnRcVersion4(t *testing.T) {
	actual := IsAnRCVersion("v4.1.9-unstable/")
	const expected = false
	if actual != expected {
		t.Errorf("Expected: '%s', got '%s'", expected, actual)
	}
}
