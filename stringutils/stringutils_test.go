package stringutils

import (
	"testing"
)

func Test_Stripchars_strips_properly(t *testing.T) {
	const original = "qwerty.asdfgh."
	actual := Stripchars(original, ".")
	const expected = "qwertyasdfgh"
	if actual != expected {
		t.Errorf("Expected: '%s', got '%s'", expected, actual)
	}
}

func Test_Stripchars_doesnt_strip_unspecified_characters(t *testing.T) {
	const original = "abc"
	actual := Stripchars(original, ".qwe")
	const expected = original
	if actual != expected {
		t.Errorf("Expected: '%s', got '%s'", expected, actual)
	}
}

func Test_Stripchars_stripping_all_characters_gives_empty_string(t *testing.T) {
	const original = "123456"
	actual := Stripchars(original, "123456")
	const expected = ""
	if actual != expected {
		t.Errorf("Expected: '%s', got '%s'", expected, actual)
	}
}

func Test_LeftPad1(t *testing.T) {
	const original = "aa"
	actual := LeftPad(original, "0", 3)
	const expected = "000aa"
	if actual != expected {
		t.Errorf("Expected: '%s', got '%s'", expected, actual)
	}
}

func Test_LeftPad2(t *testing.T) {
	const original = "abc"
	actual := LeftPad(original, "g", 3)
	const expected = "gggabc"
	if actual != expected {
		t.Errorf("Expected: '%s', got '%s'", expected, actual)
	}
}

func Test_LeftPad3(t *testing.T) {
	const original = "1"
	actual := LeftPad(original, "0", 0)
	const expected = "1"
	if actual != expected {
		t.Errorf("Expected: '%s', got '%s'", expected, actual)
	}
}

func Test_RightPad1(t *testing.T) {
	const original = "abc"
	actual := RightPad(original, "g", 3)
	const expected = "abcggg"
	if actual != expected {
		t.Errorf("Expected: '%s', got '%s'", expected, actual)
	}
}

func Test_RightPad2(t *testing.T) {
	const original = "1"
	actual := RightPad(original, "0", 5)
	const expected = "100000"
	if actual != expected {
		t.Errorf("Expected: '%s', got '%s'", expected, actual)
	}
}

func Test_RightPad3(t *testing.T) {
	const original = "1"
	actual := RightPad(original, "0", 0)
	const expected = "1"
	if actual != expected {
		t.Errorf("Expected: '%s', got '%s'", expected, actual)
	}
}

func Test_LeftPad2Len1(t *testing.T) {
	const original = "1"
	actual := LeftPad2Len(original, "0", 10)
	const expected = "0000000001"
	if actual != expected {
		t.Errorf("Expected: '%s', got '%s'", expected, actual)
	}
}

func Test_LeftPad2Len2(t *testing.T) {
	const original = "1"
	actual := LeftPad2Len(original, "a", 5)
	const expected = "aaaa1"
	if actual != expected {
		t.Errorf("Expected: '%s', got '%s'", expected, actual)
	}
}

func Test_RightPad2Len1(t *testing.T) {
	const original = "1"
	actual := RightPad2Len(original, "0", 10)
	const expected = "1000000000"
	if actual != expected {
		t.Errorf("Expected: '%s', got '%s'", expected, actual)
	}
}

func Test_RightPad2Len2(t *testing.T) {
	const original = "abc"
	actual := RightPad2Len(original, "1", 5)
	const expected = "abc11"
	if actual != expected {
		t.Errorf("Expected: '%s', got '%s'", expected, actual)
	}
}
